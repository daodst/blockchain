package rest

import (
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/comm/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

type resp struct {
	Status int         `json:"status"`
	Info   string      `json:"info"`
	Data   interface{} `json:"data"`
}


func gatewayNumIsValidHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		resp := resp{Data: true, Status: 1}
		gatewayNumber := vars["gatewayNumber"]
		if len(gatewayNumber) != 7 { 
			if rest.CheckInternalServerError(w, types.ErrGatewayNumLength) {
				return
			}
		}
		resBytes, _, err := clientCtx.QueryWithData("custom/comm/"+types.QueryGatewayNum, nil)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		gatewayNumList := make(map[string]types.GatewayNumIndex)
		if resBytes != nil {
			err := util.Json.Unmarshal(resBytes, &gatewayNumList)
			if rest.CheckInternalServerError(w, err) {
				return
			}
		}
		if _, ok := gatewayNumList[gatewayNumber]; ok {
			num := gatewayNumList[gatewayNumber]
			if num.Status != 2 { 
				resp.Data = false
				resp.Info = "number already registered"
				SendReponse(w, clientCtx, resp)
				return
			}
		}
		SendReponse(w, clientCtx, resp)
	}
}


func gatewayInfoHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		resp := resp{}
		gatewayAddress := vars["gatewayAddress"]
		params := types.QueryGatewayInfoParams{GatewayAddress: gatewayAddress}
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		valAddr, err := sdk.ValAddressFromBech32(gatewayAddress)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		resBytes, _, err := clientCtx.QueryWithData("custom/comm/"+types.QueryGatewayInfo, bz)
		if err != nil {
			
			if strings.Contains(err.Error(), types.ErrGatewayNotExist.Error()) {
				resp.Status = 1
				SendReponse(w, clientCtx, resp)
				return
			} else {
				if rest.CheckInternalServerError(w, err) {
					return
				}
			}
		}
		if resBytes != nil {
			respData := struct {
				types.Gateway
				Account string `json:"account"`
			}{}
			err = util.Json.Unmarshal(resBytes, &respData)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			respData.Account = sdk.AccAddress(valAddr).String()
			resp.Status = 1
			resp.Data = respData
		}
		SendReponse(w, clientCtx, resp)
	}
}


func gatewayNumberCountHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := resp{}
		var paramsByte []byte
		if r.Body != nil {
			paramsByte, _ = ioutil.ReadAll(r.Body)
		}
		req := types.GatewayNumberCountReq{}
		err := util.Json.Unmarshal(paramsByte, &req)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		coin := sdk.NewCoin(sdk.DefaultBondDenom, core.MustRealString2LedgerIntNoMin(req.Amount))
		params := types.GatewayNumberCountParams{GatewayAddress: req.GatewayAddress, Amount: coin}
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		resBytes, _, err := clientCtx.QueryWithData("custom/comm/"+types.QueryGatewayNumberCount, bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		var count int64
		if resBytes != nil {
			err = util.Json.Unmarshal(resBytes, &count)
			if rest.CheckInternalServerError(w, err) {
				return
			}
		}
		resp.Status = 1
		resp.Data = count
		SendReponse(w, clientCtx, resp)
	}
}


func gatewayNumberUnbondCountHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := resp{}
		var paramsByte []byte
		if r.Body != nil {
			paramsByte, _ = ioutil.ReadAll(r.Body)
		}
		req := types.GatewayNumberCountReq{}
		err := util.Json.Unmarshal(paramsByte, &req)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		coin := sdk.NewCoin(sdk.DefaultBondDenom, core.MustRealString2LedgerIntNoMin(req.Amount))
		params := types.GatewayNumberCountParams{GatewayAddress: req.GatewayAddress, Amount: coin}
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		resBytes, _, err := clientCtx.QueryWithData("custom/comm/"+types.QueryGatewayNumberUnbondCount, bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		var count int64
		if resBytes != nil {
			err = util.Json.Unmarshal(resBytes, &count)
			if rest.CheckInternalServerError(w, err) {
				return
			}
		}
		resp.Status = 1
		resp.Data = count
		SendReponse(w, clientCtx, resp)
	}
}
