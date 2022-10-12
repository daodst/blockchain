package testenv

import (
	"freemasonry.cc/blockchain/core"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

//，，
func (this *IntegrationTestSuite) initPledge(validatorAcc authtypes.AccountI, numberPrefix []string) (stakingtypes.ValidatorI, string, error) {
	app := this.Env.App
	env := this.Env
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	//validatorAcc := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("valacc"))
	app.AccountKeeper.SetAccount(ctx, validatorAcc) //

	amountDec := sdk.MustNewDecFromStr("100000000000000000000000000000")
	coins := sdk.NewCoins(
		sdk.NewCoin(core.BaseDenom, amountDec.TruncateInt()),
		sdk.NewCoin(core.GovDenom, amountDec.TruncateInt()),
	)

	this.Require().NoError(
		this.SendCoinsFromModuleToAccount(validatorAcc.GetAddress(), coins),
	)

	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	this.T().Log(":", validatorAcc.GetAddress().String())
	this.T().Log(":", app.BankKeeper.GetAllBalances(ctx, validatorAcc.GetAddress()))

	priv1 := secp256k1.GenPrivKey()
	pk1 := priv1.PubKey()

	//
	validator, err := this.RegValidator(validatorAcc.GetAddress(), pk1)
	this.Require().NoError(err)

	this.T().Log("-----")
	this.T().Log(":", validator.GetOperator().String())

	env.Height += 1
	//
	gatewayAddr, err := this.RegGateway(validatorAcc.GetAddress(), "10000000000000000000000", numberPrefix)
	this.Require().NoError(err)

	return validator, gatewayAddr, nil
}

// 
func (this *IntegrationTestSuite) TestPledge() {
	app := this.Env.App
	env := this.Env
	env.Height += 1

	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	validatorAcc := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("valacc"))

	_, gatewayAddr, err := this.initPledge(validatorAcc, []string{"1234567"})
	this.Require().NoError(err, "initPledge")

	this.T().Log("-----")
	this.T().Log(":", gatewayAddr)

	env.Height += 1
	err = this.RegChat(validatorAcc.GetAddress(), gatewayAddr, "200000000000000000000", "1234567")
	this.Require().NoError(err, "")

	this.T().Log("-----")

	env.Height += 1
	pledgeAmount2 := sdk.MustNewDecFromStr("10000000000000000000")
	err = this.ChatPledge(validatorAcc.GetAddress(), gatewayAddr, pledgeAmount2)
	this.Require().NoError(err, "")

	this.T().Log("-----")
}

// 
func (this *IntegrationTestSuite) TestUnpledge() {
	app := this.Env.App
	env := this.Env
	env.Height += 1

	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	validatorAcc := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("valacc"))

	validator, gatewayAddr, err := this.initPledge(validatorAcc, []string{"1234567"})
	this.Require().NoError(err, "initPledge")

	this.T().Log("-----")
	this.T().Log(":", gatewayAddr)

	env.Height += 1
	err = this.RegChat(validatorAcc.GetAddress(), gatewayAddr, "200000000000000000000", "1234567")
	this.Require().NoError(err, "")

	this.T().Log("-----")

	env.Height += 1
	pledgeAmount2 := sdk.MustNewDecFromStr("10000000000000000000")
	err = this.ChatPledge(validatorAcc.GetAddress(), gatewayAddr, pledgeAmount2)
	this.Require().NoError(err, "")

	this.T().Log("-----")

	params := this.pledgeKep.GetParams(ctx)
	this.T().Log(":", params.UnbondingTime.Hours(), "h")

	env.Height += 1
	err = this.UnPledge(validatorAcc.GetAddress(), validator.GetOperator(), sdk.NewInt(10000000000))
	if err != nil {
		this.T().Log("：")
		this.T().Fatal(err)
		return
	}
	this.T().Log("-----")
}

// 
func (this *IntegrationTestSuite) TestWithdrawDelegatorReward() {
	app := this.Env.App
	env := this.Env
	env.Height += 1

	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	validatorAcc := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("valacc"))

	validator, gatewayAddr, err := this.initPledge(validatorAcc, []string{"1234567"})
	this.Require().NoError(err, "initPledge")

	this.T().Log("-----")
	this.T().Log(":", gatewayAddr)

	env.Height += 1
	err = this.RegChat(validatorAcc.GetAddress(), gatewayAddr, "200000000000000000000", "1234567")
	this.Require().NoError(err, "")

	this.T().Log("-----")

	env.Height += 1
	pledgeAmount2 := sdk.MustNewDecFromStr("10000000000000000000")
	err = this.ChatPledge(validatorAcc.GetAddress(), gatewayAddr, pledgeAmount2)
	this.Require().NoError(err, "")

	this.T().Log("-----")

	// 
	this.SimulateSignatureBlocking(validatorAcc.GetAddress(), validator.GetOperator())

	reward, err := this.PledgeRewardsQuery(validatorAcc.GetAddress(), validator.GetOperator(), pledgeAmount2)
	this.Require().NoError(err, "")

	this.T().Log(":", reward)

	env.Height += 1
	err = this.PledgeRewardsGet(validatorAcc.GetAddress(), validator.GetOperator(), reward)
	this.Require().NoError(err, "")

	this.T().Log("-----")

	env.Height += 1
	err = this.UnPledge(validatorAcc.GetAddress(), validator.GetOperator(), sdk.NewInt(10000000000))
	if err != nil {
		this.T().Log("：")
		this.T().Fatal(err)
		return
	}
	this.T().Log("-----")

}
