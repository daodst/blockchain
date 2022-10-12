package keeper

import (
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/comm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			res []byte
			err error
		)
		switch path[0] {
		case types.QueryGatewayInfo: //
			return queryGatewayInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryGatewayList: //
			return queryGatewayList(ctx, k)
		case types.QueryGatewayNum: //
			return queryGatewayNum(ctx, k)
		case types.QueryGatewayRedeemNum: //
			return queryGatewayRedeemNum(ctx, k)
		case types.QueryValidatorByConsAddress: //
			return queryValidatorByConsAddress(ctx, req, k, legacyQuerierCdc)
		case types.QueryGatewayNumberCount:
			return queryGatewayNumberCount(ctx, req, k, legacyQuerierCdc)
		case types.QueryGatewayNumberUnbondCount:
			return queryGatewayNumberUnbondCount(ctx, req, k, legacyQuerierCdc)
		case types.QueryGasPrice:
			return queryGasPrice(ctx)
		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		return res, err
	}
}

func queryGasPrice(ctx sdk.Context) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	gasPrice := ctx.MinGasPrices()
	gasByte, err := util.Json.Marshal(gasPrice)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return gasByte, nil
}

//
func queryGatewayNumberUnbondCount(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.GatewayNumberCountParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	valAddr, err := sdk.ValAddressFromBech32(params.GatewayAddress)
	if err != nil {
		return nil, err
	}
	delAddr := sdk.AccAddress(valAddr)
	fmt.Println("deladdr:", delAddr.String())
	fmt.Println("valAddr:", valAddr.String())
	//
	delegation, found := k.stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, stakingTypes.ErrNoDelegation
	}
	//
	shares, err := k.stakingKeeper.ValidateUnbondAmount(
		ctx, delAddr, valAddr, params.Amount.Amount,
	)
	if err != nil {
		return nil, err
	}
	param := k.GetParams(ctx)
	//
	gateway, err := k.GetGatewayInfo(ctx, params.GatewayAddress)
	if err != nil {
		return nil, err
	}
	//
	balanceShares := delegation.Shares.Sub(shares)
	//
	num := balanceShares.QuoInt(param.MinDelegate)

	//  (-)
	hode := gateway.GatewayQuota - int64(len(gateway.GatewayNum))

	count := gateway.GatewayQuota - num.TruncateInt64() - hode
	if count < 0 {
		count = 0
	}
	countByte, err := util.Json.Marshal(count)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return countByte, nil
}

//
func queryGatewayNumberCount(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.GatewayNumberCountParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	param := k.GetParams(ctx)
	//
	gateway, err := k.GetGatewayInfo(ctx, params.GatewayAddress)
	if err != nil && err != types.ErrGatewayNotExist {
		return nil, err
	}
	num := params.Amount.Amount.Quo(param.MinDelegate)
	if err == types.ErrGatewayNotExist { //,
		return num.Marshal()
	}
	//  (-)
	hode := gateway.GatewayQuota - int64(len(gateway.GatewayNum))

	count := num.Int64() + hode

	countByte, err := util.Json.Marshal(count)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return countByte, nil
}

//
func queryValidatorByConsAddress(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var params types.QueryValidatorByConsAddrParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	//
	validator, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, params.ValidatorConsAddress)
	if !found {
		return nil, stakingTypes.ErrNoValidatorFound
	}
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, validator)
	if err != nil {
		log.WithError(err).Error("MarshalJSONIndent")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

//
func queryGatewayList(ctx sdk.Context, k Keeper) ([]byte, error) {
	log := util.BuildLog(util.GetFuncName(), util.LmChainKeeper)
	gatewayList, err := k.GetGatewayList(ctx)
	if err != nil {
		log.WithError(err).Error("GetGatewayList")
		return nil, err
	}
	gatewayListByte, err := util.Json.Marshal(gatewayList)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return gatewayListByte, nil
}

//
func queryGatewayInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := util.BuildLog(util.GetFuncName(), util.LmChainKeeper)
	var params types.QueryGatewayInfoParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	var gateway *types.Gateway
	//
	if params.GatewayAddress != "" {
		gateway, err = k.GetGatewayInfo(ctx, params.GatewayAddress)
		if err != nil {
			log.WithError(err).Error("GetGatewayInfo")
			return nil, err
		}
	}
	//
	if params.GatewayNumIndex != "" {
		gateway, err = k.GetGatewayInfoByNum(ctx, params.GatewayNumIndex)
		if err != nil {
			log.WithError(err).Error("GetGatewayInfoByNum")
			return nil, err
		}
	}
	if gateway != nil {
		gatewayByte, err := util.Json.Marshal(gateway)
		if err != nil {
			log.WithError(err).Error("Marshal")
			return nil, err
		}
		return gatewayByte, nil
	}
	return nil, nil
}

//
func queryGatewayNum(ctx sdk.Context, k Keeper) ([]byte, error) {
	log := util.BuildLog(util.GetFuncName(), util.LmChainKeeper)
	gatewayNumMap, err := k.GetGatewayNumMap(ctx)
	if err != nil {
		log.WithError(err).Error("GetGatewayNumMap")
		return nil, err
	}
	gatewayMapByte, err := util.Json.Marshal(gatewayNumMap)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return gatewayMapByte, nil
}

//
func queryGatewayRedeemNum(ctx sdk.Context, k Keeper) ([]byte, error) {
	log := util.BuildLog(util.GetFuncName(), util.LmChainKeeper)
	gatewayRedeemNumMap, err := k.GetGatewayRedeemNum(ctx)
	if err != nil {
		log.WithError(err).Error("GetGatewayRedeemNum")
		return nil, err
	}
	gatewayMapByte, err := util.Json.Marshal(gatewayRedeemNumMap)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return gatewayMapByte, nil
}
