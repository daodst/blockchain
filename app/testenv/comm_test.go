package testenv

import (
	"freemasonry.cc/blockchain/core"
	commtypes "freemasonry.cc/blockchain/x/comm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

//，
func (this *IntegrationTestSuite) initComm() (sdk.AccAddress, sdk.ValAddress, error) {
	app := this.Env.App
	env := this.Env
	env.Height = int64(1)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("account1"))

	app.AccountKeeper.SetAccount(ctx, acc1) //

	//acc1
	amountDec := sdk.MustNewDecFromStr("100000000000000000000000000000")
	acc1Coins := sdk.NewCoins(
		sdk.NewCoin(core.BaseDenom, amountDec.TruncateInt()),
		sdk.NewCoin(core.GovDenom, amountDec.TruncateInt()),
	)

	this.Require().NoError(
		this.bankKep.SendCoinsFromModuleToAccount(ctx, moduleAccountName, acc1.GetAddress(), acc1Coins),
	)

	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	this.T().Log(":", acc1.GetAddress().String())
	this.T().Log(":", app.BankKeeper.GetAllBalances(ctx, acc1.GetAddress()))

	//
	coin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.MustNewDecFromStr("10000000000000000000000000").TruncateInt())
	validator, err := this.CreateSmartValidator(acc1.GetAddress(), coin)
	if err != nil {
		this.T().Error(err)
		return acc1.GetAddress(), nil, err
	}

	return acc1.GetAddress(), validator.GetOperator(), nil
}

//dpos
func (this *IntegrationTestSuite) TestCreateSmartValidator() {

	_, validatorAddr, err := this.initComm()
	if err != nil {
		this.T().Log("")
		this.T().Error(err)
		return
	}

	app := this.Env.App
	env := this.Env
	env.Height += 1
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	validator, ok := this.stakingKep.GetValidator(ctx, validatorAddr)

	this.Require().Equal(true, ok)

	this.T().Log(":", validator)
}

//
func (this *IntegrationTestSuite) TestGatewayRegister() {
	valAccAddr, _, err := this.initComm()
	if err != nil {
		this.T().Log("")
		this.T().Error(err)
		return
	}

	app := this.Env.App
	env := this.Env
	env.Height = int64(1)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	env.Height += 1
	gatewayAddr, err := this.RegGateway(valAccAddr, "10000000000000000000000", []string{"123456"})
	this.Require().NoError(err)

	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	gateway, err := this.commKep.GetGatewayInfo(ctx, gatewayAddr)
	this.Require().NoError(err)

	this.T().Log(":", gateway)

}

//
func (this *IntegrationTestSuite) TestGatewayIndexNum() {
	valAccAddr, validatorAddr, err := this.initComm()
	if err != nil {
		this.T().Log("")
		this.T().Error(err)
		return
	}

	app := this.Env.App
	env := this.Env
	env.Height = int64(1)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	env.Height += 1
	gatewayAddr, err := this.RegGateway(valAccAddr, "5000000000000000000000", []string{})
	this.Require().NoError(err)

	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	gateway, err := this.commKep.GetGatewayInfo(ctx, gatewayAddr)
	this.Require().NoError(err)

	this.T().Log(":", gateway)
	this.T().Log(":", gateway.GatewayQuota)

	for _, v := range gateway.GatewayNum {
		this.T().Log(":", v.NumberIndex, v.Status, v.Validity)
	}

	delegator, err := app.StakingKeeper.GetDelegatorValidator(ctx, valAccAddr, validatorAddr)
	this.Require().NoError(err)
	this.T().Log(":", delegator.Tokens)

	env.Height += 1
	numbers := []string{"1888888"}
	this.T().Log(":", numbers)
	err = this.GatewayDistributeNumber(valAccAddr, validatorAddr, numbers)
	this.Require().NoError(err)

	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	gateway, err = this.commKep.GetGatewayInfo(ctx, gatewayAddr)
	this.Require().NoError(err)

	this.T().Log(":", gateway)
	this.T().Log(":", gateway.GatewayNum)
}

//
func (this *IntegrationTestSuite) TestGatewayDelegate() {

	valAccAddr, validatorAddr, err := this.initComm()
	if err != nil {
		this.T().Log("")
		this.T().Error(err)
		return
	}

	app := this.Env.App
	env := this.Env

	env.Height += 1
	gatewayAddr, err := this.RegGateway(valAccAddr, "5000000000000000000000", []string{})
	this.Require().NoError(err)

	env.Height += 1
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	gateway, err := this.commKep.GetGatewayInfo(ctx, gatewayAddr)
	this.Require().NoError(err)

	this.T().Log(":", gateway)
	this.T().Log(":", gateway.GatewayQuota)
	this.T().Log(":", len(gateway.GatewayNum))
	for _, v := range gateway.GatewayNum {
		this.T().Log(":", v.NumberIndex, ":", v.Status)
	}
	this.T().Log("--------------------------")

	delegator, err := app.StakingKeeper.GetDelegatorValidator(ctx, valAccAddr, validatorAddr)
	this.Require().NoError(err)
	this.T().Log(":", delegator.Tokens)

	env.Height += 1
	numbers := []string{"1888888"}
	this.T().Log(":", numbers)
	err = this.GatewayDistributeNumber(valAccAddr, validatorAddr, numbers)
	this.Require().NoError(err)

	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	delegate, _ := this.Env.App.StakingKeeper.GetDelegation(ctx, valAccAddr, validatorAddr)

	this.T().Log(":", delegate.Shares.TruncateInt())

	//
	coin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.MustNewDecFromStr("20000000000000000000000").TruncateInt())
	err = this.GatewayDelegation(valAccAddr, validatorAddr, coin)
	this.Require().NoError(err)
	this.T().Log(":", coin)

	gateway, err = this.commKep.GetGatewayInfo(ctx, sdk.ValAddress(valAccAddr).String())
	this.Require().NoError(err)

	this.T().Log(":", gateway)
	this.T().Log(":", gateway.GatewayQuota)
	this.T().Log(":", len(gateway.GatewayNum))
	for _, v := range gateway.GatewayNum {
		this.T().Log(":", v.NumberIndex, ":", v.Status)
	}
	this.T().Log("--------------------------")

	delegate, _ = this.Env.App.StakingKeeper.GetDelegation(ctx, valAccAddr, validatorAddr)

	this.T().Log(":", delegate.Shares.TruncateInt())
}

//
func (this *IntegrationTestSuite) TestGatewayUndelegate() {
	valAccAddr, validatorAddr, err := this.initComm()
	if err != nil {
		this.T().Log("")
		this.T().Error(err)
		return
	}

	app := this.Env.App
	env := this.Env

	env.Height += 1
	gatewayAddr, err := this.RegGateway(valAccAddr, "5000000000000000000000", []string{})
	this.Require().NoError(err)

	env.Height += 1
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	gateway, err := this.commKep.GetGatewayInfo(ctx, gatewayAddr)
	this.Require().NoError(err)

	this.T().Log(":", gateway)
	this.T().Log(":", gateway.GatewayQuota)

	for _, v := range gateway.GatewayNum {
		this.T().Log(":", v.NumberIndex, v.Status, v.Validity)
	}

	delegator, err := app.StakingKeeper.GetDelegatorValidator(ctx, valAccAddr, validatorAddr)
	this.Require().NoError(err)
	this.T().Log(":", delegator.Tokens)

	env.Height += 1
	numbers := []string{"1888888"}
	this.T().Log(":", numbers)
	err = this.GatewayDistributeNumber(valAccAddr, validatorAddr, numbers)
	this.Require().NoError(err)

	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	gateway, err = this.commKep.GetGatewayInfo(ctx, gatewayAddr)
	this.Require().NoError(err)

	this.T().Log(":", gateway)
	this.T().Log(":", gateway.GatewayNum)

	//
	coin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.MustNewDecFromStr("2000000000000000000000").TruncateInt())
	err = this.GatewayUndelegation(valAccAddr, validatorAddr, coin, []string{"1888888"})
	this.Require().NoError(err)

	gateway, err = this.commKep.GetGatewayInfo(ctx, sdk.ValAddress(valAccAddr).String())
	this.Require().NoError(err)

	this.T().Log(":", gateway)
	this.T().Log(":", gateway.GatewayNum)

	delegate, _ := this.Env.App.StakingKeeper.GetDelegation(ctx, valAccAddr, validatorAddr)

	this.T().Log(":", delegate)
}

//
func (this *IntegrationTestSuite) TestGatewayBeginRedelegate() {
	app := this.Env.App
	height := int64(1)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: height})

	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("account1"))
	acc2 := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("account2"))

	app.AccountKeeper.SetAccount(ctx, acc1) //

	//acc1
	acc1CoinDec := sdk.MustNewDecFromStr("10000000000000000000")
	acc1Coin := sdk.NewCoin("att", acc1CoinDec.TruncateInt())
	acc1Coins := sdk.NewCoins(acc1Coin)

	acc1FmCoinDec := sdk.MustNewDecFromStr("100000000000000000000000")
	acc1FmCoin := sdk.NewCoin("fm", acc1FmCoinDec.TruncateInt())
	acc1FmCoins := sdk.NewCoins(acc1FmCoin)
	//acc2
	acc2CoinDec := sdk.MustNewDecFromStr("10000000000000000000")
	acc2Coin := sdk.NewCoin("att", acc2CoinDec.TruncateInt())
	acc2Coins := sdk.NewCoins(acc2Coin)

	acc2FmCoinDec := sdk.MustNewDecFromStr("100000000000000000000000")
	acc2FmCoin := sdk.NewCoin("fm", acc2FmCoinDec.TruncateInt())
	acc2FmCoins := sdk.NewCoins(acc2FmCoin)

	this.Require().NoError(
		this.bankKep.SendCoinsFromModuleToAccount(ctx, moduleAccountName, acc1.GetAddress(), acc1Coins),
		this.bankKep.SendCoinsFromModuleToAccount(ctx, moduleAccountName, acc1.GetAddress(), acc1FmCoins),
		this.bankKep.SendCoinsFromModuleToAccount(ctx, moduleAccountName, acc2.GetAddress(), acc2Coins),
		this.bankKep.SendCoinsFromModuleToAccount(ctx, moduleAccountName, acc2.GetAddress(), acc2FmCoins),
	)

	this.bankKep.SendCoinsFromModuleToAccount(ctx, moduleAccountName, genesisAccount.GetAddress(), genesisAccountTokens)

	height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: height})

	this.T().Log(":", acc1.GetAddress().String())
	this.T().Log(":", app.BankKeeper.GetAllBalances(ctx, acc1.GetAddress()))
	this.T().Log("2:", acc2.GetAddress().String())
	this.T().Log("2:", app.BankKeeper.GetAllBalances(ctx, acc2.GetAddress()))

	height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: height})
	rate, _ := sdk.NewDecFromStr("0.1")
	commission := stakingTypes.NewCommissionRates(rate, sdk.OneDec(), sdk.OneDec())
	description := stakingTypes.Description{
		Moniker:         "test",
		Identity:        "test",
		Website:         "test",
		SecurityContact: "test",
		Details:         "test",
	}

	coinInt, _ := sdk.NewIntFromString("12000000000000000000000")
	coin := sdk.NewCoin(sdk.DefaultBondDenom, coinInt)
	msg, err := commtypes.NewMsgCreateSmartValidator(sdk.ValAddress(acc1.GetAddress()), "oK1vO9MU47i/IkPxI/oHdpp34759zQ1vFWw1VjwCJVo=", coin, description, commission, sdk.OneInt())
	msg2, err := commtypes.NewMsgCreateSmartValidator(sdk.ValAddress(acc2.GetAddress()), "2Yshx2YWNofX/u8UWDKASrS8hOsD3LA1N26vt93Q2N8=", coin, description, commission, sdk.OneInt())
	this.Require().NoError(err)
	uctx := sdk.WrapSDKContext(ctx)
	//1
	_, err = this.commServ.CreateSmartValidator(uctx, msg)
	this.Require().NoError(err)
	//2
	_, err = this.commServ.CreateSmartValidator(uctx, msg2)
	this.Require().NoError(err)

	height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: height})
	//BeginBlock
	staking.BeginBlocker(ctx, this.Env.App.StakingKeeper)
	staking.EndBlocker(ctx, this.Env.App.StakingKeeper)
	validator, found := this.Env.App.StakingKeeper.GetValidator(ctx, sdk.ValAddress(acc1.GetAddress()))
	if !found {
		return
	}
	validator2, found := this.Env.App.StakingKeeper.GetValidator(ctx, sdk.ValAddress(acc2.GetAddress()))
	if !found {
		return
	}
	this.T().Log(":", validator.GetOperator().String(), ":", validator.GetDelegatorShares())
	this.T().Log("2:", validator2.GetOperator().String(), ":", validator2.GetDelegatorShares())
	coinInt, _ = sdk.NewIntFromString("2000000000000000000000")
	coin = sdk.NewCoin(sdk.DefaultBondDenom, coinInt)
	msgRedelegate := commtypes.NewMsgGatewayBeginRedelegate(acc1.GetAddress(), validator.GetOperator(), validator2.GetOperator(), coin, []string{})
	//
	_, err = this.commServ.GatewayBeginRedelegate(uctx, msgRedelegate)
	this.Require().NoError(err)
	//BeginBlock
	staking.BeginBlocker(ctx, this.Env.App.StakingKeeper)
	staking.EndBlocker(ctx, this.Env.App.StakingKeeper)

	validator, found = this.Env.App.StakingKeeper.GetValidator(ctx, sdk.ValAddress(acc1.GetAddress()))
	if !found {
		return
	}
	validator2, found = this.Env.App.StakingKeeper.GetValidator(ctx, sdk.ValAddress(acc2.GetAddress()))
	if !found {
		return
	}
	this.T().Log(":", validator.GetOperator().String(), ":", validator.GetDelegatorShares())
	this.T().Log("2:", validator2.GetOperator().String(), ":", validator2.GetDelegatorShares())

	//
	msgGatewayRegister := commtypes.NewMsgGatewayRegister(acc1.GetAddress().String(), "", "127.0.0.1", "", []string{"1888888"})
	_, err = this.commServ.GatewayRegister(uctx, msgGatewayRegister)
	this.Require().NoError(err)
	height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: height})

	gateway, err := this.commKep.GetGatewayInfo(ctx, sdk.ValAddress(acc1.GetAddress()).String())
	this.Require().NoError(err)

	this.T().Log(":", gateway)
	this.T().Log(":", gateway.GatewayNum)

	coinInt, _ = sdk.NewIntFromString("2000000000000000000000")
	coin = sdk.NewCoin(sdk.DefaultBondDenom, coinInt)
	msgRedelegate = commtypes.NewMsgGatewayBeginRedelegate(acc1.GetAddress(), validator.GetOperator(), validator2.GetOperator(), coin, []string{"1888888"})
	//
	_, err = this.commServ.GatewayBeginRedelegate(uctx, msgRedelegate)
	this.Require().NoError(err)
	//BeginBlock
	staking.BeginBlocker(ctx, this.Env.App.StakingKeeper)
	staking.EndBlocker(ctx, this.Env.App.StakingKeeper)
	this.T().Log(":", validator.GetOperator().String(), ":", validator.GetDelegatorShares())
	this.T().Log("2:", validator2.GetOperator().String(), ":", validator2.GetDelegatorShares())
	gateway, err = this.commKep.GetGatewayInfo(ctx, sdk.ValAddress(acc1.GetAddress()).String())
	this.Require().NoError(err)

	this.T().Log(":", gateway)
	this.T().Log(":", gateway.GatewayNum)

}

func TestA(t *testing.T) {
	t.Log(sdk.NewInt(123).String())
}
