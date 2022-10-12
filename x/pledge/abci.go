package pledge

import (
	"freemasonry.cc/blockchain/core"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"freemasonry.cc/blockchain/x/pledge/keeper"
	"freemasonry.cc/blockchain/x/pledge/types"
)

// BeginBlocker will persist the current header and validator set as a historical entry and prune the oldest entry based on the HistoricalEntries parameter
// BeginBlocker，HistoricalEntries
func BeginBlocker(ctx sdk.Context, k keeper.Keeper, req abci.RequestBeginBlock) {
	logs := core.BuildLog(core.GetFuncName(), core.LmChainBeginBlock)
	logs.Debug("pledge BeginBlocker start ---")
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	//staking begin blocker
	k.TrackHistoricalInfo(ctx)

	logs.Debug("pledge BeginBlocker TrackHistoricalInfo end ---")

	//mint begin blocker
	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	logs.Debug("pledge BeginBlocker get params end ---")

	// recalculate inflation rate
	totalStakingSupply := k.StakingTokenSupply(ctx)
	bondedRatio := k.BondedRatio(ctx)
	minter.Inflation = minter.NextInflationRate(params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)
	k.SetMinter(ctx, minter)

	logs.Debug("pledge BeginBlocker set minter end ---")

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	logs.Debug("pledge BeginBlocker MintCoins end ---")

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)

	//
	pledgeFeeAccount := k.AccountKeeper.GetModuleAccount(ctx, types.FeeCollectorName)
	pledgeModuleAccountF := k.BankKeeper.GetBalance(ctx, pledgeFeeAccount.GetAddress(), "att")

	logs.Debug(" ---：", mintedCoins)
	logs.Debug("：", ctx.BlockHeight(), "pledgeModuleAccountF:-----", pledgeModuleAccountF)

	if err != nil {
		panic(err)
	}

	logs.Debug("pledge BeginBlocker AddCollectedFees end ---")

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
			sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)

	logs.Debug("GetVotes: --- ", req.LastCommitInfo.GetVotes())

	distBeginBlocker(ctx, req, k)

}

func distBeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// determine the total power signing the block
	var previousTotalPower, sumPreviousPrecommitPower int64
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
		if voteInfo.SignedLastBlock {
			sumPreviousPrecommitPower += voteInfo.Validator.Power
		}
	}

	// TODO this is Tendermint-dependent
	// ref https://github.com/cosmos/cosmos-sdk/issues/3095
	if ctx.BlockHeight() > 1 {
		previousProposer := k.GetPreviousProposerConsAddr(ctx)
		k.AllocateTokens(ctx, sumPreviousPrecommitPower, previousTotalPower, previousProposer, req.LastCommitInfo.GetVotes())
	}

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.SetPreviousProposerConsAddr(ctx, consAddr)
}

// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	//staking end blocker
	return k.BlockValidatorUpdates(ctx)
}
