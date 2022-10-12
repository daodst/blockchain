package keeper

import (
	"freemasonry.cc/blockchain/core"
	abci "github.com/tendermint/tendermint/abci/types"

	"freemasonry.cc/blockchain/x/pledge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AllocateTokens handles distribution of the collected fees
// bondedVotes is a list of (validator address, validator voted on last block flag) for all
// validators in the bonded set.
func (k Keeper) AllocateTokens(
	ctx sdk.Context, sumPreviousPrecommitPower, totalPreviousPower int64,
	previousProposer sdk.ConsAddress, bondedVotes []abci.VoteInfo,
) {
	logs := core.BuildLog(core.GetFuncName(), core.LmChainPledgeKeeper)
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.AccountKeeper.GetModuleAccount(ctx, k.FeeCollectorName)
	feesCollectedInt := k.BankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// transfer collected fees to the distribution module account
	logs.Debug(":", feesCollected)
	err := k.BankKeeper.SendCoinsFromModuleToModule(ctx, k.FeeCollectorName, types.ModuleName, feesCollectedInt)
	if err != nil {
		panic(err)
	}

	// temporary workaround to keep CanWithdrawInvariant happy
	// general discussions here: https://github.com/cosmos/cosmos-sdk/issues/2906#issuecomment-441867634
	feePool := k.GetFeePool(ctx)
	if totalPreviousPower == 0 {
		feePool.CommunityPool = feePool.CommunityPool.Add(feesCollected...)
		k.SetFeePool(ctx, feePool)
		return
	}

	voteMultiplier := sdk.OneDec()

	// allocate tokens proportionally to voting power
	// TODO consider parallelizing later, ref https://github.com/cosmos/cosmos-sdk/pull/3099#discussion_r246276376

	//todo 、

	allValidators := k.GetAllValidators(ctx)
	logs.Debug("allValidatorsCount：", len(allValidators))
	for _, v := range allValidators {
		//power
		logs.Debug("tokens：", v.GetTokens())

		power := v.GetConsensusPower(k.PowerReduction(ctx))
		logs.Debug("power:", power)
		powerFraction := sdk.NewDec(power).QuoTruncate(sdk.NewDec(totalPreviousPower))
		logs.Debug("totalPreviousPower:", totalPreviousPower)
		logs.Debug("powerFraction:", powerFraction)
		logs.Debug("voteMultiplier", voteMultiplier)
		reward := feesCollected.MulDecTruncate(voteMultiplier).MulDecTruncate(powerFraction)
		logs.Debug("：", reward)
		k.AllocateTokensToValidator(ctx, v, reward)
	}

	//for _, vote := range bondedVotes {
	//	validator := k.ValidatorByConsAddr(ctx, vote.Validator.Address)
	//
	//	// TODO consider microslashing for missing votes.
	//	// ref https://github.com/cosmos/cosmos-sdk/issues/2525#issuecomment-430838701
	//	powerFraction := sdk.NewDec(vote.Validator.Power).QuoTruncate(sdk.NewDec(totalPreviousPower))
	//	reward := feesCollected.MulDecTruncate(voteMultiplier).MulDecTruncate(powerFraction)
	//	k.AllocateTokensToValidator(ctx, validator, reward)
	//}

}

// AllocateTokensToValidator allocate tokens to a particular validator, splitting according to commission
func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val types.ValidatorI, tokens sdk.DecCoins) {
	// split tokens between validator and delegators according to commission

	shared := tokens

	// update current rewards
	currentRewards := k.GetValidatorCurrentRewards(ctx, val.GetOperator())
	currentRewards.Rewards = currentRewards.Rewards.Add(shared...)
	k.SetValidatorCurrentRewards(ctx, val.GetOperator(), currentRewards)

	// update outstanding rewards
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, tokens.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
		),
	)
	outstanding := k.GetValidatorOutstandingRewards(ctx, val.GetOperator())
	outstanding.Rewards = outstanding.Rewards.Add(tokens...)

	//fmt.Println("OutstandingRewards:---", outstanding.Rewards)

	k.SetValidatorOutstandingRewards(ctx, val.GetOperator(), outstanding)
}
