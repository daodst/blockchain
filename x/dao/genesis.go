package dao

import (
	"freemasonry.cc/blockchain/x/dao/keeper"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
)

// InitGenesis import module genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	data types.GenesisState,
) {
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{}
}
