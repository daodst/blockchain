package contract

import (
	"freemasonry.cc/blockchain/contracts"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
	//
	initCoinMetadata(ctx, k)
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
	}
}

func RegisterCoin(ctx sdk.Context, k keeper.Keeper) {
	//dst
	//dstMetaData, found := k.BankKeeper.GetDenomMetaData(ctx, core.BaseDenom)
	//if found {
	//	_, err := k.RegisterCoin(ctx, dstMetaData)
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	usdtMetaData, found := k.BankKeeper.GetDenomMetaData(ctx, core.UsdtDenom)
	if found {
		_, err := k.RegisterCoin(ctx, usdtMetaData)
		if err != nil {
			panic(err)
		}
	}
}

func initCoinMetadata(ctx sdk.Context, k keeper.Keeper) {
	coinMetadata := []banktypes.Metadata{
		//{
		//	Description: "dst coin",
		//	Base:        core.BaseDenom,
		//	DenomUnits: []*banktypes.DenomUnit{
		//		{
		//			Denom:    core.BaseDenom,
		//			Exponent: 18,
		//		},
		//	},
		//
		//	Display: core.BaseDenom,
		//	Name:    "cosmos dst",
		//	Symbol:  "edst",
		//},
		{
			Description: "usdt coin",
			Base:        core.UsdtDenom,
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    core.UsdtDenom,
					Exponent: 18,
				},
			},
			Display: core.UsdtDenom,
			Name:    "bsc bridge usdt",
			Symbol:  "usdt",
		},
	}
	for _, metadatum := range coinMetadata {
		k.BankKeeper.SetDenomMetaData(ctx, metadatum)
	}
}

// 
func DeployTokenFactoryContract(ctx sdk.Context, k keeper.Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	contractAddress := authtypes.NewModuleAddress(types.TokenFactoryContractDeploy)
	from := common.BytesToAddress(contractAddress)
	contractAddr := crypto.CreateAddress(from, 0)
	log.Info("********************token factory contract address:", contractAddr)
	StakeFactoryMedal := contracts.AppTokenIssueJSONContract
	resp, err := k.CallEVMWithData(ctx, from, nil, StakeFactoryMedal.Bin, true)
	if err != nil {
		log.WithError(err).Error("token factory contract deploy error")
		panic(err)
	}
	log.Info("********************************token factory contract deploy res:", resp.String())
	err = k.SetTokenFactoryContractAddress(ctx, contractAddr.String())

	conAddr := k.GetTokenFactoryContractAddress(ctx)
	log.Info("*************************get token factory contract addr:", conAddr)

	if err != nil {
		panic(err)
	}
}
