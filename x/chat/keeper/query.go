package keeper

import (
	// this line is used by starport scaffolding # 1

	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/chat/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			res []byte
			err error
		)
		switch path[0] {
		case types.QueryUserInfo: //
			return QueryUserInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryParams:
			return queryParams(ctx, req, k, legacyQuerierCdc)
		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		return res, err
	}
}

func QueryUserInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := util.BuildLog(util.GetFuncName(), util.LmChainKeeper)
	var params types.QueryUserInfoParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	userInfo, err := k.GetRegisterInfo(ctx, params.Address)
	if err != nil {
		log.Info("GetRegisterInfo Err:", err)
		return nil, err
	}

	userInfoByte, err := util.Json.Marshal(userInfo)
	if err != nil {
		log.Info("GetRegisterInfo Marshal Err:", err)
		return nil, err
	}

	if userInfoByte == nil {
		log.Info("GetRegisterInfo Not Fount:")
	}

	return userInfoByte, nil
}

//
func queryParams(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
