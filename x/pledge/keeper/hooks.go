package keeper

import (
	"freemasonry.cc/blockchain/cmd/config"
	"freemasonry.cc/blockchain/x/pledge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Implements StakingHooks interface
var _ stakingTypes.StakingHooks = Keeper{}

// AfterValidatorCreated - call hook if registered
func (k Keeper) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterValidatorCreated(ctx, valAddr)
	}
}

// BeforeValidatorModified - call hook if registered
func (k Keeper) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeValidatorModified(ctx, valAddr)
	}
}

// AfterValidatorRemoved - call hook if registered
func (k Keeper) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterValidatorRemoved(ctx, consAddr, valAddr)
	}
}

// AfterValidatorBonded - call hook if registered
func (k Keeper) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterValidatorBonded(ctx, consAddr, valAddr)
	}
}

// AfterValidatorBeginUnbonding - call hook if registered
func (k Keeper) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterValidatorBeginUnbonding(ctx, consAddr, valAddr)
	}
}

// BeforeDelegationCreated - call hook if registered
func (k Keeper) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationCreated(ctx, delAddr, valAddr)
	}
}

// BeforeDelegationSharesModified - call hook if registered
func (k Keeper) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	}
}

// BeforeDelegationRemoved - call hook if registered
func (k Keeper) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationRemoved(ctx, delAddr, valAddr)
	}
}

// AfterDelegationModified - call hook if registered
func (k Keeper) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterDelegationModified(ctx, delAddr, valAddr)
	}
}

// BeforeValidatorSlashed - call hook if registered
func (k Keeper) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	if k.hooks != nil {
		k.hooks.BeforeValidatorSlashed(ctx, valAddr, fraction)
	}
}

func (k Keeper) AfterValidatorRemovedPledge(ctx sdk.Context, valAddr sdk.ValAddress) {
	// fetch outstanding
	outstanding := k.GetValidatorOutstandingRewardsCoins(ctx, valAddr)

	// force-withdraw commission
	//commission := k.GetValidatorAccumulatedCommission(ctx, valAddr).Commission
	//if !commission.IsZero() {
	//	// subtract from outstanding
	//	outstanding = outstanding.Sub(commission)
	//
	//	// split into integral & remainder
	//	coins, remainder := commission.TruncateDecimal()
	//
	//	// remainder to community pool
	//	feePool := h.k.GetFeePool(ctx)
	//	feePool.CommunityPool = feePool.CommunityPool.Add(remainder...)
	//	h.k.SetFeePool(ctx, feePool)
	//
	//	// add to validator account
	//	if !coins.IsZero() {
	//		accAddr := sdk.AccAddress(valAddr)
	//		withdrawAddr := h.k.GetDelegatorWithdrawAddr(ctx, accAddr)
	//
	//		if err := h.k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, coins); err != nil {
	//			panic(err)
	//		}
	//	}
	//}

	// Add outstanding to community pool
	// The validator is removed only after it has no more delegations.
	// This operation sends only the remaining dust to the community pool.
	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(outstanding...)
	k.SetFeePool(ctx, feePool)

	// delete outstanding
	k.DeleteValidatorOutstandingRewards(ctx, valAddr)

	// remove commission record
	k.DeleteValidatorAccumulatedCommission(ctx, valAddr)

	// clear slashes
	//k.DeleteValidatorSlashEvents(ctx, valAddr)

	// clear historical rewards
	k.DeleteValidatorHistoricalRewards(ctx, valAddr)

	// clear current rewards
	k.DeleteValidatorCurrentRewards(ctx, valAddr)
}

type CommonHooks struct {
	k Keeper
}

var _ types.CommonHooks = CommonHooks{}

func (k Keeper) CommonHooks() CommonHooks {
	return CommonHooks{k}
}

func (c CommonHooks) AfterCreateGateway(ctx sdk.Context, validator stakingTypes.Validator) {
	pledgeValidator := stakingTypes.MsgCreateValidator{}
	commission := stakingTypes.NewCommissionRates(validator.Commission.Rate, validator.Commission.MaxRate, validator.Commission.MaxChangeRate)
	description := stakingTypes.Description{
		Moniker:         validator.Description.Moniker,
		Identity:        validator.Description.Identity,
		Website:         validator.Description.Website,
		SecurityContact: validator.Description.SecurityContact,
		Details:         validator.Description.Details,
	}
	pledgeValidator.Commission = commission
	pledgeValidator.Description = description
	pledgeValidator.MinSelfDelegation = validator.MinSelfDelegation
	pledgeValidator.DelegatorAddress = sdk.AccAddress(validator.GetOperator()).String()
	pledgeValidator.ValidatorAddress = validator.OperatorAddress
	pledgeValidator.Pubkey = validator.ConsensusPubkey
	pledgeValidator.Value = sdk.NewCoin(config.BaseDenom, validator.GetBondedTokens())
	err := c.k.CreateValidator(ctx, pledgeValidator)
	if err != nil {
		panic(err)
	}
}
