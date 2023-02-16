package contract

import (
	"fmt"
	"freemasonry.cc/blockchain/contracts"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"freemasonry.cc/blockchain/x/contract/keeper"
	"freemasonry.cc/blockchain/x/contract/types"
)

// InitGenesis import module genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	data types.GenesisState,
) {
	k.SetParams(ctx, data.Params)
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		//Params:     k.GetParams(ctx),
		//TokenPairs: k.GetTokenPairs(ctx),
	}
}

func ApplayNftContract(ctx sdk.Context, k keeper.Keeper, accountKeeper authkeeper.AccountKeeper) {
	contractAddress := authtypes.NewModuleAddress(types.ModuleName)
	from := common.BytesToAddress(contractAddress)
	contractAddr := crypto.CreateAddress(from, 0)
	fmt.Println("********************:", contractAddr)
	freeMasonryMedal := contracts.FreeMasonryMedalJSONContract
	resp, err := k.CallEVMWithData(ctx, from, nil, freeMasonryMedal.Bin, true)
	if err != nil {
		fmt.Println("************************:", err.Error())
		panic(err)
	}
	fmt.Println("********************************:", resp.String())
	err = k.SetNftContractAddress(ctx, contractAddr.String())
	if err != nil {
		panic(err)
	}
}
