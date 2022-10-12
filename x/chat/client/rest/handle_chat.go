package rest

import (
	"errors"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/chat/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

// RegisterHandlerFn 
func RegisterHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {

	//
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	var register types.MsgRegister
	err := util.Json.Unmarshal(msgBytes, &register)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}

	accFromAddress, err := sdk.AccAddressFromBech32(register.FromAddress)
	if err != nil {
		return err
	}

	//
	balStatus, errStr := judgeBalance(ctx, accFromAddress, register.MortgageAmount.Amount.ToDec(), register.MortgageAmount.Denom)
	if !balStatus {
		log.Error("judgeBalance fail | ", errStr)
		return errors.New(errStr)
	}
	return nil
}

func SetChatFeeHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {

	return nil
}

func SendGiftHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {

	return nil
}
func AddressBookSaveHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {

	return nil
}

func GetRewardsHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {

	return nil
}

func MobileTransferHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {

	return nil
}

func ChangeGatewayHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {

	return nil
}

func BurnGetMobileHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {

	return nil
}

func SetChatInfoHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {

	return nil
}
