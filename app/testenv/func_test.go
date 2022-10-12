package testenv

import (
	"errors"
	"freemasonry.cc/blockchain/core"
	chattypes "freemasonry.cc/blockchain/x/chat/types"
	commtypes "freemasonry.cc/blockchain/x/comm/types"
	"freemasonry.cc/blockchain/x/pledge"
	pledgetypes "freemasonry.cc/blockchain/x/pledge/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	types3 "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

/**
******************************************
******************************************
******************************************
******************************************
***  ***
******************************************
******************************************
******************************************
******************************************
 */

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (this *IntegrationTestSuite) CreateSmartValidator(accAddr sdk.AccAddress, coin sdk.Coin) (stakingtypes.ValidatorI, error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	rate, _ := sdk.NewDecFromStr("0.1")
	commission := stakingtypes.NewCommissionRates(rate, sdk.OneDec(), sdk.OneDec())
	description := stakingtypes.Description{
		Moniker:         "test",
		Identity:        "test",
		Website:         "test",
		SecurityContact: "test",
		Details:         "test",
	}
	validatorAddr := sdk.ValAddress(accAddr)
	msg, err := commtypes.NewMsgCreateSmartValidator(validatorAddr, "oK1vO9MU47i/IkPxI/oHdpp34759zQ1vFWw1VjwCJVo=", coin, description, commission, sdk.OneInt())
	if err != nil {
		return nil, err
	}
	uctx := sdk.WrapSDKContext(ctx)
	_, err = this.commServ.CreateSmartValidator(uctx, msg)
	if err != nil {
		return nil, err
	}

	staking.EndBlocker(ctx, *this.stakingKep)

	//
	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	//BeginBlock
	staking.BeginBlocker(ctx, *this.stakingKep)
	validator, ok := this.stakingKep.GetValidator(ctx, validatorAddr)
	if !ok {
		return nil, errors.New("")
	}
	return validator, nil
}

//
func (this *IntegrationTestSuite) RegValidator(accAddr sdk.AccAddress, pk1 cryptotypes.PubKey) (stakingtypes.ValidatorI, error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	//
	validatorAddr := sdk.ValAddress(accAddr)
	validatorDposCoin := sdk.NewCoin(core.GovDenom, sdk.NewInt(1000000))
	this.T().Log(":", validatorDposCoin)
	this.T().Log(":", validatorAddr.String())

	desc := stakingtypes.NewDescription("testname", "", "", "", "")
	comm := stakingtypes.NewCommissionRates(sdk.OneDec(), sdk.OneDec(), sdk.OneDec())

	msgCreateValidator, err := types3.NewMsgCreateValidator(validatorAddr, pk1, validatorDposCoin, desc, comm, sdk.NewInt(1))
	if err := msgCreateValidator.ValidateBasic(); err != nil {
		return nil, err
	}

	_, err = this.stakingServ.CreateValidator(uctx, msgCreateValidator)
	if err != nil {
		this.T().Fatal("：", err)
		return nil, err
	}
	staking.EndBlocker(ctx, *this.stakingKep)

	//
	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})

	//BeginBlock
	staking.BeginBlocker(ctx, *this.stakingKep)
	validator, ok := this.stakingKep.GetValidator(ctx, validatorAddr)
	if !ok {
		return nil, errors.New("")
	}
	return validator, nil
}

//
func (this *IntegrationTestSuite) RegGateway(accAddr sdk.AccAddress, delegation string, indexNumber []string) (gatewayAddr string, err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	gatewayRegisterMsg := commtypes.NewMsgGatewayRegister(accAddr.String(), "1", "192.168.0.117", delegation, indexNumber)

	if err = gatewayRegisterMsg.ValidateBasic(); err != nil {
		return
	}

	_, err = this.commServ.GatewayRegister(uctx, gatewayRegisterMsg)
	if err != nil {
		return
	}
	gatewayAddr = sdk.ValAddress(accAddr).String()
	return
}

//
func (this *IntegrationTestSuite) GatewayDistributeNumber(accAddr sdk.AccAddress, validatorAddr sdk.ValAddress, numbers []string) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	msgGatewayIndexNum := commtypes.NewMsgGatewayIndexNum(accAddr.String(), validatorAddr.String(), numbers)
	if err = msgGatewayIndexNum.ValidateBasic(); err != nil {
		return err
	}

	//
	_, err = this.commServ.GatewayIndexNum(uctx, msgGatewayIndexNum)

	return err
}

//dpos
func (this *IntegrationTestSuite) GatewayDelegation(accAddr sdk.AccAddress, validatorAddr sdk.ValAddress, coin sdk.Coin) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	msgDelegation := stakingtypes.NewMsgDelegate(accAddr, validatorAddr, coin)

	//dpos
	_, err = this.stakingServ.Delegate(uctx, msgDelegation)
	if err != nil {
		return err
	}
	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	//BeginBlock
	staking.BeginBlocker(ctx, *this.stakingKep)
	staking.EndBlocker(ctx, *this.stakingKep)
	return nil
}

//dpos
func (this *IntegrationTestSuite) GatewayUndelegation(accAddr sdk.AccAddress, validatorAddr sdk.ValAddress, coin sdk.Coin, numbers []string) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	msgUndelegation := commtypes.NewMsgGatewayUndelegation(accAddr.String(), validatorAddr.String(), coin, numbers)

	//
	_, err = this.commServ.GatewayUndelegate(uctx, msgUndelegation)
	if err != nil {
		return err
	}
	env.Height += 1
	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	//BeginBlock
	staking.BeginBlocker(ctx, *this.stakingKep)
	staking.EndBlocker(ctx, *this.stakingKep)
	return nil
}

//
func (this *IntegrationTestSuite) RegChat(accAddr sdk.AccAddress, gatewayAddr, amount, numberPrefix string) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	//
	regAmount := sdk.MustNewDecFromStr(amount)
	registercoin := sdk.NewCoin(core.BaseDenom, regAmount.TruncateInt())

	//this.T().Log(":", accAddr.String())
	//this.T().Log(":", gatewayAddr)
	//this.T().Log(":", registercoin)
	chatRegisterMsg := chattypes.NewMsgRegister(accAddr.String(), gatewayAddr, numberPrefix, registercoin)

	if err = chatRegisterMsg.ValidateBasic(); err != nil {
		return err
	}

	_, err = this.chatKep.Register(uctx, chatRegisterMsg)
	if err != nil {
		return err
	}

	//
	pledgeParams := this.pledgeKep.GetParams(ctx)
	//this.T().Log("-------------------------")
	//this.T().Log("[pledge]")
	//this.T().Log(":", pledgeParams.AttDestroyPercent)
	//this.T().Log(":", pledgeParams.AttGatewayPercent)
	//this.T().Log("dpos:", pledgeParams.AttDposPercent)
	//this.T().Log(":", pledgeParams.PreAttCoin)
	//this.T().Log(":", pledgeParams.PreAttAccount)
	//
	//this.T().Log("-------------------------")
	//this.T().Log("[]")
	pledgeSum, err := this.pledgeKep.GetPledgeSum(ctx, accAddr.String())
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
		return err
	}

	kouchubili := pledgeParams.AttDestroyPercent.Add(pledgeParams.AttGatewayPercent).Add(pledgeParams.AttDposPercent)
	kouchu := regAmount.Mul(kouchubili)
	//this.T().Log(":", accAddr.String())
	//this.T().Log(":", regAmount)
	//this.T().Log(":", kouchu)
	//this.T().Log(":", pledgeSum)

	if !pledgeSum.Equal(regAmount.Sub(kouchu).TruncateInt()) {
		this.T().Error("")
		return err
	}
	return nil
}

//
func (this *IntegrationTestSuite) PledgeRewardsQuery(accAddr sdk.AccAddress, validatorAddr sdk.ValAddress, pledgeAmount sdk.Dec) (coins sdk.Coins, err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	return this.pledgeKep.QueryDelegationRewards(ctx, validatorAddr, accAddr)
}

//
func (this *IntegrationTestSuite) PledgeRewardsGet(accAddr sdk.AccAddress, validatorAddr sdk.ValAddress, rewardAmount sdk.Coins) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	acc1BanlanceInit := this.bankKep.GetBalance(ctx, accAddr, core.BaseDenom)

	receiveMsg := pledgetypes.NewMsgMsgWithdrawDelegatorReward(accAddr.String(), validatorAddr.String())
	if err = receiveMsg.ValidateBasic(); err != nil {
		return err
	}

	_, err = this.pledgeServ.WithdrawDelegatorReward(uctx, receiveMsg)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
		return err
	}
	acc1Banlance := this.bankKep.GetBalance(ctx, accAddr, core.BaseDenom)
	this.T().Log("acc1:", acc1Banlance)

	if !acc1Banlance.Amount.Sub(rewardAmount.AmountOf(core.BaseDenom)).Equal(acc1BanlanceInit.Amount) {
		acc1lq := rewardAmount.AmountOf(core.BaseDenom)
		this.T().Log("acc1:", acc1BanlanceInit.Amount)
		this.T().Log("acc1:", acc1lq)
		this.T().Log("acc1:", acc1BanlanceInit.Amount.Add(acc1lq))
		this.T().Log("acc1:", acc1Banlance.Amount)
		this.T().Error("")
		return err
	}
	return nil
}

//
func (this *IntegrationTestSuite) ChatPledge(accAddr sdk.AccAddress, gatewayAddr string, pledgeAmount sdk.Dec) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	pledgeSumInit := sdk.NewInt(0)
	pledgeSumInit, err = this.pledgeKep.GetPledgeSum(ctx, accAddr.String())
	if err != nil {
		this.T().Log("1")
		return err
	}

	pledgeMsg := pledgetypes.NewMsgPledge(accAddr.String(), accAddr.String(), gatewayAddr, sdk.NewCoin(core.BaseDenom, pledgeAmount.TruncateInt()))
	if err = pledgeMsg.ValidateBasic(); err != nil {
		return err
	}

	_, err = this.pledgeServ.Delegate(uctx, pledgeMsg)
	if err != nil {
		this.T().Log("")
		return err
	}
	//
	pledgeParams := this.pledgeKep.GetParams(ctx)
	//
	pledgeSum, err := this.pledgeKep.GetPledgeSum(ctx, accAddr.String())
	if err != nil {
		this.T().Log("2")
		return err
	}

	kouchubili := pledgeParams.AttDestroyPercent.Add(pledgeParams.AttGatewayPercent).Add(pledgeParams.AttDposPercent)
	kouchu := pledgeAmount.Mul(kouchubili)
	//this.T().Log(":", accAddr.String())
	//this.T().Log(":", pledgeAmount)
	//this.T().Log(":", kouchu)
	//this.T().Log(":", pledgeSum)

	if !pledgeSum.Sub(pledgeSumInit).Equal(pledgeAmount.Sub(kouchu).TruncateInt()) {
		this.T().Error("")
		return err
	}
	return nil
}

//
func (this *IntegrationTestSuite) UnPledge(accAddr sdk.AccAddress, validatorAddr sdk.ValAddress, amount sdk.Int) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)
	//
	coin := sdk.NewCoin(core.BaseDenom, amount)
	msgUnpledge := pledgetypes.NewMsgUnpledge(accAddr.String(), validatorAddr.String(), coin)

	//
	balanceInit := this.bankKep.GetBalance(ctx, accAddr, core.BaseDenom).Amount
	if err != nil {
		this.T().Log("1")
		return err
	}

	if err = msgUnpledge.ValidateBasic(); err != nil {
		return err
	}
	_, err = this.pledgeServ.Undelegate(uctx, msgUnpledge)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
		return err
	}

	//，endblock
	this.SimulateSignatureBlocking(accAddr, validatorAddr)

	ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	//
	pledgeParams := this.pledgeKep.GetParams(ctx)
	//
	balanceCurr := this.bankKep.GetBalance(ctx, accAddr, core.BaseDenom).Amount
	if err != nil {
		this.T().Log("2")
		return err
	}

	amountDec := sdk.NewDecFromInt(amount)
	kouchubili := pledgeParams.AttDestroyPercent.Add(pledgeParams.AttGatewayPercent).Add(pledgeParams.AttDposPercent)
	shouxufei := amountDec.Mul(kouchubili)
	this.T().Log(":", accAddr.String())
	this.T().Log(":", balanceInit)
	this.T().Log(":", balanceCurr)
	this.T().Log(":", amount)
	this.T().Log(":", shouxufei)
	this.T().Log(":", amountDec.Sub(shouxufei))
	this.T().Log(":", balanceCurr.Sub(balanceInit))

	if !balanceCurr.Sub(balanceInit).Equal(amountDec.Sub(shouxufei).TruncateInt()) {
		this.T().Error("")
		return err
	}
	return nil
}

//
func (this *IntegrationTestSuite) SetChatFee(accAddr sdk.AccAddress, amount sdk.Int) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	coin := sdk.NewCoin(core.BaseDenom, amount)
	msg := chattypes.NewMsgSetChatFee(accAddr.String(), coin)

	if err = msg.ValidateBasic(); err != nil {
		return err
	}
	_, err = this.chatServ.SetChatFee(uctx, msg)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
		return err
	}
	return nil
}

//
func (this *IntegrationTestSuite) AddressBookSave(accAddr sdk.AccAddress, addrbook []string) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	msg := chattypes.NewMsgAddressBookSave(accAddr.String(), addrbook)

	if err = msg.ValidateBasic(); err != nil {
		return err
	}
	_, err = this.chatServ.AddressBookSave(uctx, msg)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
		return err
	}
	return nil
}

//
func (this *IntegrationTestSuite) BurnGetMobile(accAddr sdk.AccAddress, mobilePrefix string) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	msg := chattypes.NewMsgBurnGetMobile(accAddr.String(), mobilePrefix)

	if err = msg.ValidateBasic(); err != nil {
		return err
	}
	_, err = this.chatServ.BurnGetMobile(uctx, msg)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
		return err
	}
	return nil
}

//
func (this *IntegrationTestSuite) ChangeGateway(accAddr sdk.AccAddress, gatewayAddr string) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	msg := chattypes.NewMsgChangeGateway(accAddr.String(), gatewayAddr)

	if err = msg.ValidateBasic(); err != nil {
		return err
	}
	_, err = this.chatServ.ChangeGateway(uctx, msg)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
		return err
	}
	return nil
}

//
func (this *IntegrationTestSuite) MobileTransfer(fromAddr sdk.AccAddress, toAddr sdk.AccAddress, mobileNumber string) (err error) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	uctx := sdk.WrapSDKContext(ctx)

	msg := chattypes.NewMsgMobileTransfer(fromAddr.String(), toAddr.String(), mobileNumber)

	if err = msg.ValidateBasic(); err != nil {
		return err
	}
	_, err = this.chatServ.MobileTransfer(uctx, msg)
	if err != nil {
		this.T().Log("")
		this.T().Fatal(err)
		return err
	}
	return nil
}

//
func (this *IntegrationTestSuite) SimulateSignatureBlocking(accAddr sdk.AccAddress, validatorAddr sdk.ValAddress) {
	env := this.Env
	app := env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
	this.pledgeKep.SetPreviousProposerConsAddr(ctx, sdk.ConsAddress(accAddr))

	for i := 0; i < 50; i++ {
		env.Height += 1
		ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: env.Height})
		req := abci.RequestBeginBlock{
			Hash: []byte("test"),
			LastCommitInfo: abci.LastCommitInfo{
				Round: int32(env.Height),
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

		req.Header.ProposerAddress = accAddr
		pledge.BeginBlocker(ctx, *this.pledgeKep, req)
	}
}

//
func (this *IntegrationTestSuite) SendCoinsFromModuleToAccount(accAddr sdk.AccAddress, coins sdk.Coins) error {
	app := this.Env.App
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: this.Env.Height})
	return this.bankKep.SendCoinsFromModuleToAccount(ctx, moduleAccountName, accAddr, coins)
}
