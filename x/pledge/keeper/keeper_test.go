package keeper_test

import (
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/cmd/config"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/chat/types"
	commKeeper "freemasonry.cc/blockchain/x/comm/keeper"
	types2 "freemasonry.cc/blockchain/x/comm/types"
	"freemasonry.cc/blockchain/x/pledge"
	keeper2 "freemasonry.cc/blockchain/x/pledge/keeper"
	types4 "freemasonry.cc/blockchain/x/pledge/types"
	"freemasonry.cc/log"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	types3 "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
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

func init() {
	log.InitLogger(logrus.DebugLevel)
}

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
		chatApp.BankKeeper.MintCoins(ctx, "chat", genesisTokens.Add(sdk.NewCoin(core.GovDenom, genesisAmountFm))),
	)

	//this.Require().NoError(
	//	chatApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, "chat", genesisAccount.GetAddress(), genesisTokens),
	//)

	balance0 := chatApp.BankKeeper.GetBalance(ctx, genesisAccount.GetAddress(), "att")

	this.T().Log("", genesisAccount.GetAddress().String(), "att:", balance0)
	//chatApp.StakingKeeper.Delegation(ctx,genesisAccount.GetAddress(),addr1)
	this.app = chatApp
	this.ctx = ctx

	//commK := chatApp.GetKey(types2.StoreKey)
	//commSub := chatApp.GetSubspace(types2.ModuleName)
	//
	//commkeeper := commKeeper.NewKeeper(commK, chatApp.AppCodec(), commSub, chatApp.AccountKeeper, chatApp.BankKeeper, &chatApp.StakingKeeper)
	//pledgeKeeper := keeper2.NewKeeper(commK, chatApp.AppCodec(), commSub, chatApp.AccountKeeper, chatApp.BankKeeper, types4.FeeCollectorName)
	//chatkeeper := keeper.NewKeeper(chatApp.GetKey(types.StoreKey), chatApp.AppCodec(), chatApp.GetSubspace(types.ModuleName),
	//	chatApp.AccountKeeper, chatApp.BankKeeper, commkeeper, pledgeKeeper)

	//
	//chatParams := types.DefaultParams()
	//chatParams.ChatRewardLog = append(chatParams.ChatRewardLog, types.ChatReward{
	//	Height: 14403,
	//	Value:  "0.02",
	//})

	//chatkeeper.SetParams(ctx, chatParams)
	//commkeeper.SetParams(ctx, types2.DefaultParams())

	cfg := sdk.GetConfig()
	config.SetBech32Prefixes(cfg)
	config.SetBip44CoinType(cfg)
}

func (this *IntegrationTestSuite) TestPledge() {
	cfg := sdk.GetConfig()
	config.SetBech32Prefixes(cfg)
	config.SetBip44CoinType(cfg)
	app := this.app
	height := int64(1)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: height})

	priv1 := secp256k1.GenPrivKey()
	pk1 := priv1.PubKey()
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, sdk.AccAddress(pk1.Address()))

	app.AccountKeeper.SetAccount(ctx, acc1) //

	//acc1
	acc1CoinDec := sdk.MustNewDecFromStr("10000000000000000000")
	acc1Coin := sdk.NewCoin(config.BaseDenom, acc1CoinDec.TruncateInt())
	acc1Coins := sdk.NewCoins(acc1Coin)

	//acc1FmDec := sdk.MustNewDecFromStr("100000000000000000000")
	//acc1FmCoin := sdk.NewCoin(core.GovDenom, acc1FmDec.TruncateInt())
	//acc1FmCoins := sdk.NewCoins(acc1FmCoin)
	this.T().Log("-------------------------")
	this.T().Log("[]")
	this.Require().NoError(
		this.app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "chat", acc1.GetAddress(), acc1Coins),
	)

	this.Require().NoError(
		this.app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "chat", acc1.GetAddress(), genesisTokens),
	)

	balance0 := this.app.BankKeeper.GetBalance(ctx, acc1.GetAddress(), "fm")

	this.T().Log("acc1", genesisAccount.GetAddress().String(), "acc1fm:", balance0)

	height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: height})
	uctx := sdk.WrapSDKContext(ctx)

	this.T().Log("-------------------------")
	this.T().Log("[]")
	//
	validatorAddr := sdk.ValAddress(acc1.GetAddress())
	validatorDposCoin := sdk.NewCoin(core.GovDenom, sdk.NewInt(1000000))
	this.T().Log(":", validatorDposCoin)
	this.T().Log(":", validatorAddr.String())

	desc := stakingtypes.NewDescription("testname", "", "", "", "")
	comm := stakingtypes.NewCommissionRates(sdk.OneDec(), sdk.OneDec(), sdk.OneDec())

	msgCreateValidator, err := types3.NewMsgCreateValidator(validatorAddr, pk1, validatorDposCoin, desc, comm, sdk.NewInt(1))
	stakingServer := stakingKeeper.NewMsgServerImpl(this.app.StakingKeeper)
	_, err = stakingServer.CreateValidator(uctx, msgCreateValidator)
	if err != nil {
		this.T().Fatal("：", err)
	}
	staking.EndBlocker(ctx, this.app.StakingKeeper)

	// pledge
	//pledgeKeeper := app.PledgeKeeper
	//this.Require().NoError(
	//	pledgeKeeper.CreateValidator(ctx, *msgCreateValidator),
	//)

	//
	height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: height})
	uctx = sdk.WrapSDKContext(ctx)

	//BeginBlock
	staking.BeginBlocker(ctx, this.app.StakingKeeper)

	//
	commServer := commKeeper.NewMsgServerImpl(this.app.CommKeeper)

	chuangshi := genesisAccount.GetAddress().String()

	this.T().Log("-------------------------")
	this.T().Log(":", chuangshi)
	gatewayRegisterMsg := types2.NewMsgGatewayRegister(acc1.GetAddress().String(), "1", "192.168.0.117", "10000000000000000000000", []string{"1234567"})

	_, err = commServer.GatewayRegister(uctx, gatewayRegisterMsg)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
	}

	//
	gatewayList, err := this.app.CommKeeper.GetGatewayList(ctx)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
	}
	this.T().Log(":", gatewayList)

	this.T().Log("-------------------------")
	this.T().Log("[]")

	//
	regAmount := sdk.MustNewDecFromStr("200000000000000000000")
	registercoin := sdk.NewCoin("att", regAmount.TruncateInt())
	regGateway := gatewayList[0].GatewayAddress

	this.T().Log(":", acc1.GetAddress().String())
	this.T().Log(":", regGateway)
	this.T().Log(":", registercoin)
	chatRegisterMsg := types.NewMsgRegister(acc1.GetAddress().String(), regGateway, "1234567", registercoin)
	_, err = this.app.ChatKeeper.Register(uctx, chatRegisterMsg)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
	}

	//
	pledgeParams := this.app.PledgeKeeper.GetParams(ctx)
	this.T().Log("-------------------------")
	this.T().Log("[pledge]")
	this.T().Log(":", pledgeParams.AttDestroyPercent)
	this.T().Log(":", pledgeParams.AttGatewayPercent)
	this.T().Log("dpos:", pledgeParams.AttDposPercent)
	this.T().Log(":", pledgeParams.PreAttCoin)
	this.T().Log(":", pledgeParams.PreAttAccount)

	this.T().Log("-------------------------")
	this.T().Log("[]")
	pledgeSum, err := this.app.PledgeKeeper.GetPledgeSum(ctx, acc1.GetAddress().String())
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
	}

	kouchubili := pledgeParams.AttDestroyPercent.Add(pledgeParams.AttGatewayPercent).Add(pledgeParams.AttDposPercent)
	kouchu := regAmount.Mul(kouchubili)
	this.T().Log(":", acc1.GetAddress().String())
	this.T().Log(":", regAmount)
	this.T().Log(":", kouchu)
	this.T().Log(":", pledgeSum)

	if !pledgeSum.Equal(regAmount.Sub(kouchu).TruncateInt()) {
		this.T().Error("")
		return
	}

	//
	pledgeServer := keeper2.NewMsgServerImpl(this.app.PledgeKeeper)
	pledgeAmount2 := sdk.MustNewDecFromStr("10000000000000000000")

	this.T().Log("-------------------------")
	this.T().Log("[]")
	this.T().Log(":", acc1.GetAddress().String())
	this.T().Log(":", pledgeAmount2)

	pledgeMsg := types4.NewMsgPledge(acc1.GetAddress().String(), acc1.GetAddress().String(), gatewayList[0].GatewayAddress, sdk.NewCoin("att", pledgeAmount2.TruncateInt()))
	_, err = pledgeServer.Delegate(uctx, pledgeMsg)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
	}

	this.T().Log("-------------------------")
	this.T().Log("[]")
	this.T().Log(":", acc1.GetAddress().String())
	pledgeSum, err = this.app.PledgeKeeper.GetPledgeSum(ctx, acc1.GetAddress().String())
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
	}

	kouchu = pledgeAmount2.Mul(kouchubili)
	this.T().Log(":", pledgeAmount2)
	this.T().Log(":", kouchu)
	this.T().Log(":", pledgeSum)

	this.T().Log("-------------------------")
	this.T().Log("[]")
	// 
	acc1Banlance := this.app.BankKeeper.GetBalance(ctx, acc1.GetAddress(), "att")
	this.T().Log("acc1:", acc1Banlance)

	pledgeModuleAccount := this.app.AccountKeeper.GetModuleAccount(ctx, types4.ModuleName)
	pledgeModuleAccountB := this.app.BankKeeper.GetBalance(ctx, pledgeModuleAccount.GetAddress(), "att")
	this.T().Log(":", pledgeModuleAccountB) //k.feeCollectorName

	pledgeFeeAccount := this.app.AccountKeeper.GetModuleAccount(ctx, types4.FeeCollectorName)
	pledgeModuleAccountF := this.app.BankKeeper.GetBalance(ctx, pledgeFeeAccount.GetAddress(), "att")
	this.T().Log(":", pledgeModuleAccountF) //k.feeCollectorName

	this.app.PledgeKeeper.SetPreviousProposerConsAddr(ctx, sdk.ConsAddress(acc1.GetAddress()))

	for i := 0; i < 50; i++ {
		height += 1
		ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: height})
		req := abci.RequestBeginBlock{
			Hash: []byte("test"),
			LastCommitInfo: abci.LastCommitInfo{
				Round: int32(height),
				Votes: []abci.VoteInfo{
					{
						Validator: abci.Validator{
							Address: validatorAddr,
							Power:   178,
						},
						SignedLastBlock: true,
					},
				},
			},
		}

		req.Header.ProposerAddress = acc1.GetAddress()
		pledge.BeginBlocker(ctx, this.app.PledgeKeeper, req)
	}
	delegationRewards, err := this.app.PledgeKeeper.QueryDelegationRewards(ctx, validatorAddr, acc1.GetAddress())
	if err != nil {
		this.T().Log(":")
		this.T().Fatal(err)
	}
	this.T().Log(":", delegationRewards)

	this.T().Log("-------------------------")
	this.T().Log("[]")
	//
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: height})
	uctx = sdk.WrapSDKContext(ctx)
	receiveMsg := types4.NewMsgMsgWithdrawDelegatorReward(acc1.GetAddress().String(), gatewayList[0].GatewayAddress)
	_, err = pledgeServer.WithdrawDelegatorReward(uctx, receiveMsg)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
	}
	acc1BanlanceOld := acc1Banlance
	acc1Banlance = this.app.BankKeeper.GetBalance(ctx, acc1.GetAddress(), "att")
	this.T().Log("acc1:", acc1Banlance)

	if !acc1Banlance.Amount.Sub(delegationRewards.AmountOf("att")).Equal(acc1BanlanceOld.Amount) {
		acc1lq := delegationRewards.AmountOf("att")
		this.T().Log("acc1:", acc1BanlanceOld.Amount)
		this.T().Log("acc1:", acc1lq)
		this.T().Log("acc1:", acc1BanlanceOld.Amount.Add(acc1lq))
		this.T().Log("acc1:", acc1Banlance.Amount)
		this.T().Error("")
		return
	}

	pledgeModuleAccountB = this.app.BankKeeper.GetBalance(ctx, pledgeModuleAccount.GetAddress(), "att")
	this.T().Log(":", pledgeModuleAccountB) //k.feeCollectorName
	pledgeFeeAccount = this.app.AccountKeeper.GetModuleAccount(ctx, types4.FeeCollectorName)
	this.T().Log(":", pledgeModuleAccountF) //k.feeCollectorName

	//
	height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: height})
	uctx = sdk.WrapSDKContext(ctx)

	req := abci.RequestBeginBlock{
		Hash: []byte("test"),
		LastCommitInfo: abci.LastCommitInfo{
			Round: int32(height),
			Votes: []abci.VoteInfo{
				{
					Validator: abci.Validator{
						Address: validatorAddr,
						Power:   178,
					},
					SignedLastBlock: true,
				},
			},
		},
	}

	req.Header.ProposerAddress = acc1.GetAddress()
	pledge.BeginBlocker(ctx, this.app.PledgeKeeper, req)
	unPledgeCoin := sdk.NewCoin("att", sdk.MustNewDecFromStr("10000000000000000000").TruncateInt())

	this.T().Log("-------------------------")
	this.T().Log("[]")
	this.T().Log("：", acc1.GetAddress().String())
	this.T().Log(":", unPledgeCoin)

	//
	acc1Banlance = this.app.BankKeeper.GetBalance(ctx, acc1.GetAddress(), "att")
	this.T().Log("acc1:", acc1Banlance)

	msgUnpledge := types4.NewMsgUnpledge(acc1.GetAddress().String(), validatorAddr.String(), unPledgeCoin)

	_, err = pledgeServer.Undelegate(uctx, msgUnpledge)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
	}

	acc1Banlance = this.app.BankKeeper.GetBalance(ctx, acc1.GetAddress(), "att")
	this.T().Log("acc1:", acc1Banlance)
	this.T().Log(":", this.app.PledgeKeeper.GetValidatorOutstandingRewards(ctx, validatorAddr))
}
