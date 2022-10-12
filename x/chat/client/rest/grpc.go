package rest

import (
	"errors"
	"freemasonry.cc/blockchain/core"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/sirupsen/logrus"
)

var grpcLogPrefix = "restGrpc"

//grpc
func grpcQueryBalance(cliCtx *client.Context, address sdk.AccAddress, denom string) (coin sdk.Coin, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithFields(logrus.Fields{"addr": address, "denom": denom})
	params := bankTypes.QueryBalanceRequest{Address: address.String(), Denom: denom}
	bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return coin, errors.New(MarshalError)
	}
	resBytes, _, err := cliCtx.QueryWithData("custom/bank/balance", bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return coin, errors.New(QueryChainInforError)
	}
	err = cliCtx.LegacyAmino.UnmarshalJSON(resBytes, &coin)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return coin, errors.New(UnmarshalError + "2")
	}

	return coin, nil
}
