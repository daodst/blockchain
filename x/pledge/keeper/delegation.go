package keeper

import (
	"fmt"
	"freemasonry.cc/blockchain/core"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"freemasonry.cc/blockchain/x/pledge/types"
)

// initialize starting info for a new delegation
func (k Keeper) initializeDelegation(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) {
	// period has already been incremented - we want to store the period ended by this delegation action
	previousPeriod := k.GetValidatorCurrentRewards(ctx, val).Period - 1

	// increment reference count for the period we're going to track
	k.incrementReferenceCount(ctx, val, previousPeriod)

	validator := k.Validator(ctx, val)
	delegation := k.Delegation(ctx, del, val)

	// calculate delegation stake in tokens
	// we don't store directly, so multiply delegation shares * (tokens per share)
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	stake := validator.TokensFromSharesTruncated(delegation.GetShares())
	k.SetDelegatorStartingInfo(ctx, val, del, types.NewDelegatorStartingInfo(previousPeriod, stake, uint64(ctx.BlockHeight())))
}

// calculate the rewards accrued by a delegation between two periods
func (k Keeper) calculateDelegationRewardsBetween(ctx sdk.Context, val types.ValidatorI,
	startingPeriod, endingPeriod uint64, stake sdk.Dec) (rewards sdk.DecCoins) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainPledgeKeeper)

	log.Debug("startingPeriod:", startingPeriod, "endingPeriod:", endingPeriod)
	// sanity check
	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}

	// sanity check
	if stake.IsNegative() {
		panic("stake should not be negative")
	}

	// return staking * (ending - starting)
	starting := k.GetValidatorHistoricalRewards(ctx, val.GetOperator(), startingPeriod)
	ending := k.GetValidatorHistoricalRewards(ctx, val.GetOperator(), endingPeriod)
	difference := ending.CumulativeRewardRatio.Sub(starting.CumulativeRewardRatio)
	if difference.IsAnyNegative() {
		panic("negative rewards should not be possible")
	}
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	rewards = difference.MulDecTruncate(stake)
	return
}

// calculate the total rewards accrued by a delegation
func (k Keeper) CalculateDelegationRewards(ctx sdk.Context, val types.ValidatorI, del types.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainPledgeKeeper)

	// fetch starting info for delegation
	startingInfo := k.GetDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())

	log.Debug("startingInfo", "startingPeriod:", startingInfo.PreviousPeriod, "Stake:", startingInfo.Stake, "height:", startingInfo.Height)

	if startingInfo.Height == uint64(ctx.BlockHeight()) {
		// started this height, no rewards yet
		return
	}

	startingPeriod := startingInfo.PreviousPeriod
	stake := startingInfo.Stake

	// Iterate through slashes and withdraw with calculated staking for
	// distribution periods. These period offsets are dependent on *when* slashes
	// happen - namely, in BeginBlock, after rewards are allocated...
	// Slashes which happened in the first block would have been before this
	// delegation existed, UNLESS they were slashes of a redelegation to this
	// validator which was itself slashed (from a fault committed by the
	// redelegation source validator) earlier in the same BeginBlock.
	startingHeight := startingInfo.Height
	// Slashes this block happened after reward allocation, but we have to account for them for the stake sanity check below.
	endingHeight := uint64(ctx.BlockHeight())
	if endingHeight > startingHeight {
		k.IterateValidatorSlashEventsBetween(ctx, del.GetValidatorAddr(), startingHeight, endingHeight,
			func(height uint64, event types.ValidatorSlashEvent) (stop bool) {
				endingPeriod := event.ValidatorPeriod
				if endingPeriod > startingPeriod {
					rewards = rewards.Add(k.calculateDelegationRewardsBetween(ctx, val, startingPeriod, endingPeriod, stake)...)

					// Note: It is necessary to truncate so we don't allow withdrawing
					// more rewards than owed.
					stake = stake.MulTruncate(sdk.OneDec().Sub(event.Fraction))
					startingPeriod = endingPeriod
				}
				return false
			},
		)
	}

	// A total stake sanity check; Recalculated final stake should be less than or
	// equal to current stake here. We cannot use Equals because stake is truncated
	// when multiplied by slash fractions (see above). We could only use equals if
	// we had arbitrary-precision rationals.
	currentStake := val.TokensFromShares(del.GetShares())

	if stake.GT(currentStake) {
		// AccountI for rounding inconsistencies between:
		//
		//     currentStake: calculated as in staking with a single computation
		//     stake:        calculated as an accumulation of stake
		//                   calculations across validator's distribution periods
		//
		// These inconsistencies are due to differing order of operations which
		// will inevitably have different accumulated rounding and may lead to
		// the smallest decimal place being one greater in stake than
		// currentStake. When we calculated slashing by period, even if we
		// round down for each slash fraction, it's possible due to how much is
		// being rounded that we slash less when slashing by period instead of
		// for when we slash without periods. In other words, the single slash,
		// and the slashing by period could both be rounding down but the
		// slashing by period is simply rounding down less, thus making stake >
		// currentStake
		//
		// A small amount of this error is tolerated and corrected for,
		// however any greater amount should be considered a breach in expected
		// behaviour.
		marginOfErr := sdk.SmallestDec().MulInt64(3)
		if stake.LTE(currentStake.Add(marginOfErr)) {
			stake = currentStake
		} else {
			panic(fmt.Sprintf("calculated final stake for delegator %s greater than current stake"+
				"\n\tfinal stake:\t%s"+
				"\n\tcurrent stake:\t%s",
				del.GetDelegatorAddr(), stake, currentStake))
		}
	}

	// calculate rewards for final period
	rewards = rewards.Add(k.calculateDelegationRewardsBetween(ctx, val, startingPeriod, endingPeriod, stake)...)
	return rewards
}

//
func (k Keeper) QueryDelegationRewards(ctx sdk.Context, valAddr sdk.ValAddress, delAddr sdk.AccAddress) (sdk.Coins, error) {
	logs := core.BuildLog(core.GetFuncName(), core.LmChainPledgeKeeper)

	delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return sdk.NewCoins(), types.ErrNoDelegatorForAddress
	}

	// check existence of delegator starting info
	if !k.HasDelegatorStartingInfo(ctx, delegation.GetValidatorAddr(), delegation.GetDelegatorAddr()) {
		return nil, types.ErrEmptyDelegationDistInfo
	}

	// end current period and calculate rewards
	validator, found := k.GetValidator(ctx, delegation.GetValidatorAddr())
	if !found {
		return sdk.NewCoins(), types.ErrNoValidatorFound
	}

	//todo  period
	endingPeriod := k.IncrementValidatorPeriod(ctx, validator)
	rewardsRaw := k.CalculateDelegationRewards(ctx, validator, delegation, endingPeriod)

	outstanding := k.GetValidatorOutstandingRewardsCoins(ctx, delegation.GetValidatorAddr())

	// defensive edge case may happen on the very final digits of the decCoins due to operation order of the distribution mechanism.
	rewards := rewardsRaw.Intersect(outstanding)
	if !rewards.IsEqual(rewardsRaw) {
		logs.Info(
			"rounding error withdrawing rewards from validator",
			"delegator", delegation.GetDelegatorAddr().String(),
			"validator", validator.GetOperator().String(),
			"got", rewards.String(),
			"expected", rewardsRaw.String(),
		)
	}
	// truncate coins, return remainder to community pool
	coins, _ := rewards.TruncateDecimal()
	return coins, nil
}

func (k Keeper) withdrawDelegationRewards(ctx sdk.Context, val types.ValidatorI, del types.DelegationI) (sdk.Coins, error) {
	logs := core.BuildLog(core.GetFuncName(), core.LmChainPledgeKeeper)
	// check existence of delegator starting info
	if !k.HasDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr()) {
		return nil, types.ErrEmptyDelegationDistInfo
	}

	// end current period and calculate rewards

	//todo  period
	endingPeriod := k.IncrementValidatorPeriod(ctx, val)
	rewardsRaw := k.CalculateDelegationRewards(ctx, val, del, endingPeriod)
	logs.Debug("CalculateDelegationRewards:", rewardsRaw)

	outstanding := k.GetValidatorOutstandingRewardsCoins(ctx, del.GetValidatorAddr())

	logs.Debug("，：", outstanding.String())

	// defensive edge case may happen on the very final digits of the decCoins due to operation order of the distribution mechanism.
	rewards := rewardsRaw.Intersect(outstanding)
	if !rewards.IsEqual(rewardsRaw) {
		logs.Info(
			"rounding error withdrawing rewards from validator",
			"delegator", del.GetDelegatorAddr().String(),
			"validator", val.GetOperator().String(),
			"got", rewards.String(),
			"expected", rewardsRaw.String(),
		)
	}
	// truncate coins, return remainder to community pool
	coins, remainder := rewards.TruncateDecimal()

	// add coins to user account
	if !coins.IsZero() {
		withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, del.GetDelegatorAddr())
		err := k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, coins)
		if err != nil {
			return nil, err
		}
	}

	// update the outstanding rewards and the community pool only if the transaction was successful
	k.SetValidatorOutstandingRewards(ctx, del.GetValidatorAddr(), types.ValidatorOutstandingRewards{Rewards: outstanding.Sub(rewards)})
	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(remainder...)
	k.SetFeePool(ctx, feePool)

	// decrement reference count of starting period
	startingInfo := k.GetDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())
	startingPeriod := startingInfo.PreviousPeriod
	k.decrementReferenceCount(ctx, del.GetValidatorAddr(), startingPeriod)

	// remove delegator starting info
	k.DeleteDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())

	return coins, nil
}

// ValidateUnbondAmount validates that a given unbond or redelegation amount is valied based on upon the converted shares. If the amount is valid, the total amount of respective shares is returned, otherwise an error is returned.
func (k Keeper) ValidateUnbondAmount(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, amt sdk.Int,
) (shares sdk.Dec, err error) {
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return shares, types.ErrNoValidatorFound
	}

	del, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return shares, types.ErrNoDelegation
	}

	shares, err = validator.SharesFromTokens(amt)
	if err != nil {
		return shares, err
	}

	sharesTruncated, err := validator.SharesFromTokensTruncated(amt)
	if err != nil {
		return shares, err
	}

	delShares := del.GetShares()
	if sharesTruncated.GT(delShares) {
		return shares, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid shares amount")
	}

	// Cap the shares at the delegation's shares. Shares being greater could occur
	// due to rounding, however we don't want to truncate the shares or take the
	// minimum because we want to allow for the full withdraw of shares from a
	// delegation.
	if shares.GT(delShares) {
		shares = delShares
	}

	return shares, nil
}

func (k Keeper) Undelegate(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sharesAmount sdk.Dec,
) (sdk.Coin, error) {

	logs := core.BuildLog(core.GetFuncName(), core.LmChainPledgeKeeper)
	//validator, found := k.GetValidator(ctx, valAddr)
	//if !found {
	//	return time.Time{}, types.ErrNoDelegatorForAddress
	//}

	returnAmount, err := k.Unbond(ctx, delAddr, valAddr, sharesAmount)
	if err != nil {
		return sdk.Coin{}, err
	}
	logs.Info("Undelegate returnAmount:", returnAmount.String())

	// transfer the validator tokens to the not bonded pool
	//if validator.IsBonded() {
	//	k.bondedTokensToNotBonded(ctx, returnAmount)
	//}

	//completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
	//ubd := k.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
	//k.InsertUBDQueue(ctx, ubd, completionTime)

	//，
	undelegateCoin := sdk.NewCoin(k.BondDenom(ctx), returnAmount)
	logs.Info("")
	err = k.BankKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.BondedPoolName, delAddr, sdk.NewCoins(undelegateCoin))
	if err != nil {
		logs.Info("Undelegate err:", err)
		return sdk.Coin{}, err
	}

	return undelegateCoin, nil
}

func (k Keeper) Unbond(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec,
) (amount sdk.Int, err error) {
	// check if a delegation object exists in the store
	delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return amount, types.ErrNoDelegatorForAddress
	}

	// call the before-delegation-modified hook
	//k.BeforeDelegationSharesModified(ctx, delAddr, valAddr)

	val := k.Validator(ctx, valAddr)
	del := k.Delegation(ctx, delAddr, valAddr)
	_, err = k.withdrawDelegationRewards(ctx, val, del)
	if err != nil {
		panic(err)
	}

	//
	amDec, err := k.GetPledgeDelegate(ctx, delAddr.String())
	if err != nil {
		return amount, err
	}
	if !amDec.IsZero() {
		//
		delegation.Shares = delegation.Shares.Sub(amDec)
	}
	// ensure that we have enough shares to remove
	if delegation.Shares.LT(shares) {
		return amount, sdkerrors.Wrap(types.ErrNotEnoughDelegationShares, delegation.Shares.String())
	}

	// get validator
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return amount, types.ErrNoValidatorFound
	}

	// subtract shares from delegation
	delegation.Shares = delegation.Shares.Sub(shares)

	//delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	//if err != nil {
	//	return amount, err
	//}

	//isValidatorOperator := delegatorAddress.Equals(validator.GetOperator())

	// If the delegation is the operator of the validator and undelegating will decrease the validator's
	// self-delegation below their minimum, we jail the validator.
	//if isValidatorOperator && !validator.Jailed &&
	//	validator.TokensFromShares(delegation.Shares).TruncateInt().LT(validator.MinSelfDelegation) {
	//	k.jailValidator(ctx, validator)
	//	validator = k.mustGetValidator(ctx, validator.GetOperator())
	//}

	// remove the delegation
	if delegation.Shares.IsZero() {
		k.RemoveDelegation(ctx, delegation)
	} else {
		k.SetDelegation(ctx, delegation)
		// call the after delegation modification hook
		//k.AfterDelegationModified(ctx, delegatorAddress, delegation.GetValidatorAddr())
		k.initializeDelegation(ctx, delegation.GetValidatorAddr(), delegation.GetDelegatorAddr())
	}

	// remove the shares and coins from the validator
	// NOTE that the amount is later (in keeper.Delegation) moved between staking module pools
	validator, amount = k.RemoveValidatorTokensAndShares(ctx, validator, shares)

	if validator.DelegatorShares.IsZero() { //&& validator.IsUnbonded() {
		// if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
		k.RemoveValidator(ctx, validator.GetOperator())
	}

	return amount, nil
}

// RemoveDelegation removes a delegation.
func (k Keeper) RemoveDelegation(ctx sdk.Context, delegation types.Delegation) {
	delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	// TODO: Consider calling hooks outside of the store wrapper functions, it's unobvious.
	k.BeforeDelegationRemoved(ctx, delegatorAddress, delegation.GetValidatorAddr())
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDelegationKey(delegatorAddress, delegation.GetValidatorAddr()))
}

// InsertUBDQueue inserts an unbonding delegation to the appropriate timeslice in the unbonding queue.
func (k Keeper) InsertUBDQueue(ctx sdk.Context, ubd types.UnbondingDelegation,
	completionTime time.Time) {
	dvPair := types.DVPair{DelegatorAddress: ubd.DelegatorAddress, ValidatorAddress: ubd.ValidatorAddress}

	timeSlice := k.GetUBDQueueTimeSlice(ctx, completionTime)
	if len(timeSlice) == 0 {
		k.SetUBDQueueTimeSlice(ctx, completionTime, []types.DVPair{dvPair})
	} else {
		timeSlice = append(timeSlice, dvPair)
		k.SetUBDQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// GetUBDQueueTimeSlice gets a specific unbonding queue timeslice. A timeslice
// is a slice of DVPairs corresponding to unbonding delegations that expire at a
// certain time.
func (k Keeper) GetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvPairs []types.DVPair) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetUnbondingDelegationTimeKey(timestamp))
	if bz == nil {
		return []types.DVPair{}
	}

	pairs := types.DVPairs{}
	k.cdc.MustUnmarshal(bz, &pairs)

	return pairs.Pairs
}

// SetUBDQueueTimeSlice sets a specific unbonding queue timeslice.
func (k Keeper) SetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []types.DVPair) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&types.DVPairs{Pairs: keys})
	store.Set(types.GetUnbondingDelegationTimeKey(timestamp), bz)
}
