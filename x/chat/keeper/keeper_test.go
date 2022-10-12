package keeper_test

import (
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/cmd/config"
	"freemasonry.cc/blockchain/x/chat/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"
	"testing"
)

const (
	genesisAccountKey = "genesis"
)

var (
	genesisAmountAtt = sdk.MustNewDecFromStr("9999999999999999999999999999").TruncateInt()
	genesisAmountFm  = sdk.MustNewDecFromStr("9999999999999999999999999999").TruncateInt()
)

var (
	genesisAccount    = authtypes.NewEmptyModuleAccount(genesisAccountKey) //
	chatModuleAccount = authtypes.NewEmptyModuleAccount("chat", authtypes.Minter)

	genesisTokens = sdk.NewCoins(
		sdk.NewCoin(config.BaseDenom, genesisAmountAtt),
		sdk.NewCoin("fm", genesisAmountFm),
	)
)

type IntegrationTestSuite struct {
	suite.Suite

	app         *app.Evmos
	ctx         sdk.Context
	queryClient types.QueryClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

//
func (this *IntegrationTestSuite) SetupTest() {
	this.T().Log("SetupTest()")

	feemarketTypes := feemarkettypes.GenesisState{
		Params: feemarkettypes.Params{
			NoBaseFee:                false,
			BaseFeeChangeDenominator: 8,
			ElasticityMultiplier:     2,
			EnableHeight:             0,
			BaseFee:                  sdk.NewInt(1000000000),
		},
		BlockGas: 0,
	}
	chatApp := app.Setup(false, &feemarketTypes)
	ctx := chatApp.BaseApp.NewContext(false, tmproto.Header{})
	chatApp.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	chatApp.BankKeeper.SetParams(ctx, banktypes.DefaultParams())

	//queryHelper := baseapp.NewQueryServerTestHelper(ctx, chatApp.InterfaceRegistry())
	//types.RegisterQueryServer(queryHelper, chatApp.BankKeeper)
	//queryClient := types.NewQueryClient(queryHelper)
	authtypes.NewModuleAddress("chat_burn")
	chatApp.AccountKeeper.SetModuleAccount(ctx, genesisAccount) //
	chatApp.AccountKeeper.SetModuleAccount(ctx, chatModuleAccount)
	//
	this.Require().NoError(
		chatApp.BankKeeper.MintCoins(ctx, "chat", genesisTokens.Add(sdk.NewCoin(config.BaseDenom, genesisAmountAtt))),
	)

	this.Require().NoError(
		chatApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, "chat", genesisAccount.GetAddress(), genesisTokens),
	)

	balance0 := chatApp.BankKeeper.GetBalance(ctx, genesisAccount.GetAddress(), "att")

	this.T().Log("", genesisAccount.GetAddress().String(), ":", balance0)
	//chatApp.StakingKeeper.Delegation(ctx,genesisAccount.GetAddress(),addr1)
	this.app = chatApp
	this.ctx = ctx

	//commK := chatApp.GetKey(types2.StoreKey)
	//commSub := chatApp.GetSubspace(types2.ModuleName)

	//commkeeper := commKeeper.NewKeeper(commK, chatApp.AppCodec(), commSub, chatApp.AccountKeeper, chatApp.BankKeeper, &chatApp.StakingKeeper)

	//pledgekeeper := pledgekeeper.NewKeeper(commK, chatApp.AppCodec(), commSub, chatApp.AccountKeeper, chatApp.BankKeeper, types4.FeeCollectorName)

	//chatkeeper := keeper.NewKeeper(chatApp.GetKey(types.StoreKey), chatApp.AppCodec(), chatApp.GetSubspace(types.ModuleName), chatApp.AccountKeeper, chatApp.BankKeeper, commkeeper, pledgekeeper)

	cfg := sdk.GetConfig()
	config.SetBech32Prefixes(cfg)
	config.SetBip44CoinType(cfg)
}
