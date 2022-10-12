package keeper

import (
	types "freemasonry.cc/blockchain/x/pledge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the bank module's state from a given genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context) {
	k.SetParams(ctx, types.DefaultParams())

	k.SetMinter(ctx, types.Minter{
		Inflation:        sdk.NewDecWithPrec(55, 1),
		AnnualProvisions: sdk.NewDec(0),
	})

	k.SetFeePool(ctx, InitialFeePool())
}

func InitialFeePool() types.FeePool {
	return types.FeePool{
		CommunityPool: sdk.DecCoins{},
	}
}

// ExportGenesis returns the bank module's genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {

	return &types.GenesisState{}
}
