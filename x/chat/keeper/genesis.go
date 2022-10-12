package keeper

import (
	types2 "freemasonry.cc/blockchain/x/chat/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

// InitGenesis initializes the bank module's state from a given genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context) {
	k.SetParams(ctx, types2.DefaultParams())
}

// ExportGenesis returns the bank module's genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {

	return &types.GenesisState{}
}
