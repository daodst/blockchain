package rest

import (
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/pledge/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
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

func delegationHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32delegator := vars["delegatorAddr"]
		bech32validator := vars["validatorAddr"]

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		validatorAddr, err := sdk.ValAddressFromBech32(bech32validator)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		params := types.QueryDelegatorValidatorRequest{DelegatorAddr: delegatorAddr.String(), ValidatorAddr: validatorAddr.String()}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.RouterKey, types.QueryDelegation), bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}
