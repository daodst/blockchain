package testenv

import (
	"freemasonry.cc/blockchain/core"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"strconv"
)

//，，
func (this *IntegrationTestSuite) initChat(validatorAcc authtypes.AccountI, numberPrefix []string) (string, error) {
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

	return gatewayAddr, nil
}

//
func (this *IntegrationTestSuite) TestRegister() {
	app := this.Env.App
	env := this.Env
	env.Height += 1
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	validatorAcc := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("valacc"))
	gatewayAddr, err := this.initChat(validatorAcc, []string{"1234567"})
	this.Require().NoError(err)

	this.T().Log("-----")
	this.T().Log(":", gatewayAddr)

	accList := [10]sdk.AccAddress{}

	//acc1
	amountDec := sdk.MustNewDecFromStr("100000000000000000000000000000")
	coins := sdk.NewCoins(
		sdk.NewCoin(core.BaseDenom, amountDec.TruncateInt()),
		sdk.NewCoin(core.GovDenom, amountDec.TruncateInt()),
	)

	//
	for i := 0; i < len(accList); i++ {
		accInterface := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("account"+strconv.Itoa(i)))
		app.AccountKeeper.SetAccount(ctx, accInterface) //
		accList[i] = accInterface.GetAddress()
		this.T().Log(":", accList[i].String())

		this.SendCoinsFromModuleToAccount(accList[i], coins)
	}

	for i := 0; i < len(accList); i++ {
		env.Height += 1
		err = this.RegChat(accList[i], gatewayAddr, "200000000000000000000", "1234567")
		this.Require().NoError(err)

		this.T().Log(accList[i].String(), "-----")

		userInfo, err := this.chatKep.GetRegisterInfo(ctx, accList[i].String())
		this.Require().NoError(err)

		this.T().Log(accList[i].String(), "")
		this.T().Log(":", userInfo.FromAddress, ":", userInfo.ChatFee, ":", userInfo.Mobile, ":", userInfo.NodeAddress)
	}

}

//
func (this *IntegrationTestSuite) TestSetChatFee() {
	app := this.Env.App
	env := this.Env
	env.Height += 1
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	validatorAcc := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("valacc"))
	gatewayAddr, err := this.initChat(validatorAcc, []string{"1234567"})
	this.Require().NoError(err)

	this.T().Log("-----")
	this.T().Log(":", gatewayAddr)

	accList := [10]sdk.AccAddress{}

	//acc1
	amountDec := sdk.MustNewDecFromStr("100000000000000000000000000000")
	coins := sdk.NewCoins(
		sdk.NewCoin(core.BaseDenom, amountDec.TruncateInt()),
		sdk.NewCoin(core.GovDenom, amountDec.TruncateInt()),
	)

	//
	for i := 0; i < len(accList); i++ {
		accInterface := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("account"+strconv.Itoa(i)))
		app.AccountKeeper.SetAccount(ctx, accInterface) //
		accList[i] = accInterface.GetAddress()
		this.T().Log(":", accList[i].String())

		this.SendCoinsFromModuleToAccount(accList[i], coins)
	}

	for i := 0; i < len(accList); i++ {
		env.Height += 1
		err = this.RegChat(accList[i], gatewayAddr, "200000000000000000000", "1234567")
		this.Require().NoError(err)

		this.T().Log(accList[i].String(), "-----")

		env.Height += 1
		err = this.SetChatFee(accList[i], sdk.MustNewDecFromStr("100000000000000000").TruncateInt())
		this.Require().NoError(err)
		this.T().Log(accList[i].String(), "-----")

		env.Height += 1
		userInfo, err := this.chatKep.GetRegisterInfo(ctx, accList[i].String())
		this.Require().NoError(err)

		this.T().Log(accList[i].String(), "")
		this.T().Log(":", userInfo.FromAddress, ":", userInfo.ChatFee, ":", userInfo.Mobile, ":", userInfo.NodeAddress)
	}
}

//
func (this *IntegrationTestSuite) TestAddressBookSave() {
	app := this.Env.App
	env := this.Env
	env.Height += 1
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	validatorAcc := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("valacc"))
	gatewayAddr, err := this.initChat(validatorAcc, []string{"1234567"})
	this.Require().NoError(err)

	this.T().Log("-----")
	this.T().Log(":", gatewayAddr)

	accList := [10]sdk.AccAddress{}

	amountDec := sdk.MustNewDecFromStr("100000000000000000000000000000")
	coins := sdk.NewCoins(
		sdk.NewCoin(core.BaseDenom, amountDec.TruncateInt()),
		sdk.NewCoin(core.GovDenom, amountDec.TruncateInt()),
	)

	//
	for i := 0; i < len(accList); i++ {
		accInterface := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("account"+strconv.Itoa(i)))
		app.AccountKeeper.SetAccount(ctx, accInterface) //
		accList[i] = accInterface.GetAddress()
		this.T().Log(":", accList[i].String())

		this.SendCoinsFromModuleToAccount(accList[i], coins)
	}

	for i := 0; i < len(accList); i++ {
		env.Height += 1
		err = this.RegChat(accList[i], gatewayAddr, "200000000000000000000", "1234567")
		this.Require().NoError(err)

		this.T().Log(accList[i].String(), "-----")

		env.Height += 1
		err = this.SetChatFee(accList[i], sdk.MustNewDecFromStr("100000000000000000").TruncateInt())
		this.Require().NoError(err)
		this.T().Log(accList[i].String(), "-----")

		env.Height += 1
		err = this.AddressBookSave(accList[i], []string{"15863927002", "15863927001"})
		this.Require().NoError(err)
		this.T().Log(accList[i].String(), "-----")

		env.Height += 1
		userInfo, err := this.chatKep.GetRegisterInfo(ctx, accList[i].String())
		this.Require().NoError(err)

		this.T().Log(accList[i].String(), "")
		this.T().Log(":", userInfo.FromAddress, ":", userInfo.ChatFee, ":", userInfo.Mobile, ":", userInfo.NodeAddress)
	}
}

//
func (this *IntegrationTestSuite) TestMobileTransfer() {

}

//
func (this *IntegrationTestSuite) TestChangeGateway() {
	app := this.Env.App
	env := this.Env
	env.Height += 1
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	validatorAcc1 := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("valacc1"))
	validatorAcc2 := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("valacc2"))

	gatewayAddr1, err := this.initChat(validatorAcc1, []string{"1234567"})
	this.Require().NoError(err)

	this.T().Log("1-----")
	this.T().Log("1:", gatewayAddr1)

	gatewayAddr2, err := this.initChat(validatorAcc2, []string{"1234666"})
	this.Require().NoError(err)

	this.T().Log("2-----")
	this.T().Log("2:", gatewayAddr2)

	accList := [3]sdk.AccAddress{}

	//acc1
	amountDec := sdk.MustNewDecFromStr("100000000000000000000000000000")
	coins := sdk.NewCoins(
		sdk.NewCoin(core.BaseDenom, amountDec.TruncateInt()),
		sdk.NewCoin(core.GovDenom, amountDec.TruncateInt()),
	)

	//
	for i := 0; i < len(accList); i++ {
		accInterface := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("account"+strconv.Itoa(i)))
		app.AccountKeeper.SetAccount(ctx, accInterface) //
		accList[i] = accInterface.GetAddress()
		this.T().Log(":", accList[i].String())

		this.SendCoinsFromModuleToAccount(accList[i], coins)
	}

	//
	for i := 0; i < len(accList); i++ {
		env.Height += 1
		err = this.RegChat(accList[i], gatewayAddr1, "200000000000000000000", "1234567")
		this.Require().NoError(err)

		this.T().Log(accList[i].String(), "-----")

		userInfo, err := this.chatKep.GetRegisterInfo(ctx, accList[i].String())
		this.Require().NoError(err)

		this.T().Log(":", userInfo.FromAddress, ":", userInfo.ChatFee, ":", userInfo.Mobile, ":", userInfo.NodeAddress)
	}

	this.T().Log("----------------------------")

	//
	for i := 0; i < len(accList); i++ {
		env.Height += 1
		err = this.ChangeGateway(accList[i], gatewayAddr2)
		this.Require().NoError(err)

		this.T().Log(accList[i].String(), "", gatewayAddr2, "-----")

		userInfo, err := this.chatKep.GetRegisterInfo(ctx, accList[i].String())
		this.Require().NoError(err)

		this.T().Log(":", userInfo.FromAddress, ":", userInfo.ChatFee, ":", userInfo.Mobile, ":", userInfo.NodeAddress)
	}
}

//
func (this *IntegrationTestSuite) TestBurnGetMobile() {
	app := this.Env.App
	env := this.Env
	env.Height += 1
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	validatorAcc1 := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("valacc1"))
	gatewayAddr, err := this.initChat(validatorAcc1, []string{"1234567"})
	this.Require().NoError(err)

	this.T().Log("-----")
	this.T().Log(":", gatewayAddr)

	accList := [10]sdk.AccAddress{}

	//acc1
	amountDec := sdk.MustNewDecFromStr("100000000000000000000000000000")
	coins := sdk.NewCoins(
		sdk.NewCoin(core.BaseDenom, amountDec.TruncateInt()),
		sdk.NewCoin(core.GovDenom, amountDec.TruncateInt()),
	)

	//
	for i := 0; i < len(accList); i++ {
		accInterface := app.AccountKeeper.NewAccountWithAddress(ctx, authtypes.NewModuleAddress("account"+strconv.Itoa(i)))
		app.AccountKeeper.SetAccount(ctx, accInterface) //
		accList[i] = accInterface.GetAddress()
		this.T().Log(":", accList[i].String())

		this.SendCoinsFromModuleToAccount(accList[i], coins)
	}

	for i := 0; i < len(accList); i++ {
		env.Height += 1
		err = this.RegChat(accList[i], gatewayAddr, "200000000000000000000", "1234567")
		this.Require().NoError(err)

		this.T().Log(accList[i].String(), "-----")

		env.Height += 1
		err = this.BurnGetMobile(accList[i], "1234567")
		this.Require().NoError(err)

		this.T().Log(accList[i].String(), "-----")

		env.Height += 1
		userInfo, err := this.chatKep.GetRegisterInfo(ctx, accList[i].String())
		this.Require().NoError(err)

		this.T().Log(accList[i].String(), "")
		this.T().Log("FromAddress:", userInfo.FromAddress)
		this.T().Log("ChatFee:", userInfo.ChatFee)
		this.T().Log("Mobile:", userInfo.Mobile)
		this.T().Log("NodeAddress:", userInfo.NodeAddress)
	}
}
