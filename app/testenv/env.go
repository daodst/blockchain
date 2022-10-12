package testenv

import (
	"fmt"
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/cmd/config"
	"freemasonry.cc/blockchain/core"
	chatkeeper "freemasonry.cc/blockchain/x/chat/keeper"
	chattypes "freemasonry.cc/blockchain/x/chat/types"
	commkeeper "freemasonry.cc/blockchain/x/comm/keeper"
	commtypes "freemasonry.cc/blockchain/x/comm/types"
	pledgekeeper "freemasonry.cc/blockchain/x/pledge/keeper"
	pledgeTypes "freemasonry.cc/blockchain/x/pledge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"
)

const (
	genesisAccountKey = "genesis"
	moduleAccountName = "chat"
)

var (
	// = 1
	genesisAmountAtt = sdk.MustNewDecFromStr("1000000000000000000000000000000000000000").TruncateInt() //att
	genesisAmountFm  = sdk.MustNewDecFromStr("1000000000000000000000000000000000000000").TruncateInt() //fm
)

var (
	genesisAccount    = authtypes.NewEmptyModuleAccount(genesisAccountKey) //
	chatModuleAccount = authtypes.NewEmptyModuleAccount(moduleAccountName, authtypes.Minter)

	// = 1
	genesisAccountTokens = sdk.NewCoins(
		sdk.NewCoin(core.BaseDenom, sdk.MustNewDecFromStr("10000000000000000000000").TruncateInt()),
		sdk.NewCoin(core.GovDenom, sdk.MustNewDecFromStr("10000000000000000000000").TruncateInt()),
	)
)

type IntegrationTestSuite struct {
	suite.Suite
	Env        *TestAppEnv
	bankKep    bankkeeper.Keeper
	accountKep *authkeeper.AccountKeeper

	commKep  *commkeeper.Keeper
	commServ commtypes.MsgServer

	stakingKep  *stakingkeeper.Keeper
	stakingServ stakingtypes.MsgServer

	chatKep  *chatkeeper.Keeper
	chatServ chattypes.MsgServer

	pledgeKep  *pledgekeeper.Keeper
	pledgeServ pledgeTypes.MsgServer
}

//
func (this *IntegrationTestSuite) SetupTest() {
	this.T().Log("SetupTest()")

	env, err := InitTestEnv()
	if err != nil {
		this.T().Fatal(err)
	}
	this.Env = env

	this.commKep = &env.App.CommKeeper

	this.commServ = commkeeper.NewMsgServerImpl(*this.commKep)

	this.chatKep = &env.App.ChatKeeper

	this.chatServ = chatkeeper.NewMsgServerImpl(*this.chatKep)

	this.pledgeKep = &env.App.PledgeKeeper

	this.pledgeServ = pledgekeeper.NewMsgServerImpl(*this.pledgeKep)

	this.bankKep = env.App.BankKeeper
	this.accountKep = &env.App.AccountKeeper

	this.stakingKep = &env.App.StakingKeeper
	this.stakingServ = stakingkeeper.NewMsgServerImpl(*this.stakingKep)
}

type TestAppEnv struct {
	App    *app.Evmos
	Ctx    sdk.Context
	Height int64
}

//
func InitTestEnv() (*TestAppEnv, error) {

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

	env := &TestAppEnv{
		Ctx: ctx,
		App: chatApp,
	}

	chatApp.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	chatApp.BankKeeper.SetParams(ctx, banktypes.DefaultParams())

	//queryHelper := baseapp.NewQueryServerTestHelper(ctx, chatApp.InterfaceRegistry())
	//types.RegisterQueryServer(queryHelper, chatApp.BankKeeper)
	//queryClient := types.NewQueryClient(queryHelper)
	authtypes.NewModuleAddress("chat_burn")
	chatApp.AccountKeeper.SetModuleAccount(ctx, genesisAccount) //
	chatApp.AccountKeeper.SetModuleAccount(ctx, chatModuleAccount)

	//
	err := chatApp.BankKeeper.MintCoins(ctx, moduleAccountName, sdk.NewCoins(
		sdk.NewCoin(core.BaseDenom, genesisAmountAtt),
		sdk.NewCoin(core.GovDenom, genesisAmountFm),
	))
	if err != nil {
		return env, err
	}

	//
	err = chatApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, moduleAccountName, genesisAccount.GetAddress(), genesisAccountTokens)
	if err != nil {
		return env, err
	}

	balance0 := chatApp.BankKeeper.GetAllBalances(ctx, genesisAccount.GetAddress())

	fmt.Println("", genesisAccount.GetAddress().String(), ":", balance0)
	//chatApp.StakingKeeper.Delegation(ctx,genesisAccount.GetAddress(),addr1)

	pledgeKey := chatApp.GetKey(pledgeTypes.StoreKey)
	pledgeSub := chatApp.GetSubspace(pledgeTypes.ModuleName)

	pledgeKep := pledgekeeper.NewKeeper(pledgeKey, chatApp.AppCodec(), pledgeSub, chatApp.AccountKeeper, chatApp.BankKeeper, chatApp.CommKeeper, pledgeTypes.FeeCollectorName)
	pledgeParams := pledgeKep.GetParams(ctx)
	pledgeParams.UnbondingTime = 0 //0
	pledgeKep.SetParams(ctx, pledgeParams)

	chatApp.PledgeKeeper = pledgeKep
	commK := chatApp.GetKey(commtypes.StoreKey)
	commSub := chatApp.GetSubspace(commtypes.ModuleName)

	commKep := commkeeper.NewKeeper(commK, chatApp.AppCodec(), commSub, chatApp.AccountKeeper, chatApp.BankKeeper, &chatApp.StakingKeeper)
	commKep.SetParams(ctx, commtypes.DefaultParams())
	chatApp.CommKeeper = commKep
	chatApp.CommKeeper.SetHooks(
		commtypes.NewMultiStakingHooks(
			pledgeKep.CommonHooks(),
		),
	)

	chatApp.PledgeKeeper = pledgeKep

	cfg := sdk.GetConfig()
	config.SetBech32Prefixes(cfg)
	config.SetBip44CoinType(cfg)
	return env, nil
}
