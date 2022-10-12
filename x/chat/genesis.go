package chat

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	"freemasonry.cc/blockchain/x/chat/keeper"
	"freemasonry.cc/blockchain/x/chat/types"
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
	return &types.GenesisState{
		//Params:     k.GetParams(ctx),
		//TokenPairs: k.GetTokenPairs(ctx),
	}
}
