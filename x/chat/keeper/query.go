package keeper

import (
	// this line is used by starport scaffolding # 1

	"freemasonry.cc/blockchain/core"
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
		case types.QueryUserInfo: 
			return QueryUserInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryParams:
			return queryParams(ctx, req, k, legacyQuerierCdc)
		case types.QueryUserInfos:
			return QueryUserInfos(ctx, req, k, legacyQuerierCdc)
		case types.QueryUserByMobile:
			return QueryUserByMobile(ctx, req, k, legacyQuerierCdc)
		case types.QueryUsersChatInfo:
			return QueryUsersChatInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryAddrByChatAddr:
			return QueryAddrByChatAddr(ctx, req, k, legacyQuerierCdc)
		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		return res, err
	}
}

func QueryAddrByChatAddr(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainChatQuery)
	var params string

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	res, err := k.GetAddrFromChatAddr(ctx, params)
	if err != nil {
		return nil, err
	}

	resbyte, err := codec.MarshalJSONIndent(legacyQuerierCdc, res)
	if err != nil {
		log.WithError(err).Info("QueryUsersChatInfo Marshal Err")
		return nil, err
	}

	return resbyte, nil
}

func QueryUsersChatInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainChatQuery)
	var params types.QueryUserInfosParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	res := k.GetGatewayProfixMobiles(ctx, params.Addresses)

	byte, err := codec.MarshalJSONIndent(legacyQuerierCdc, res)
	if err != nil {
		log.WithError(err).Info("QueryUsersChatInfo Marshal Err")
		return nil, err
	}

	return byte, nil
}

func QueryUserByMobile(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainChatQuery)
	var params types.QueryUserByMobileParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	address, err := k.GetUserByMobile(ctx, params.Mobile)
	if err != nil {
		log.WithError(err).Info("GetUserByMobile Err")
		return nil, err
	}

	//accAddress, err := sdk.AccAddressFromBech32(address)
	//if err != nil {
	//	log.WithError(err).Info("Address Format Error")
	//	return nil, err
	//}

	allUserInfo := types.AllUserInfo{}
	
	userInfo, err := k.GetRegisterInfo(ctx, address)
	if err != nil {
		log.WithError(err).Info("GetRegisterInfo Err")
		return nil, err
	}
	allUserInfo.UserInfo = userInfo

	//todo 
	//pledgeLevelInfo, err := k.pledgeKeeper.QueryPledgeLevelByAccAddress(ctx, accAddress)
	//allUserInfo.PledgeLevel = pledgeLevelInfo.Level

	
	gateWayInfo, err := k.commKeeper.GetGatewayInfo(ctx, userInfo.NodeAddress)
	if err != nil {
		log.WithError(err).Error("GetGatewayInfo err")
		return nil, err
	}
	if gateWayInfo.Status == 1 {
		allUserInfo.GatewayProfixMobile = ""
	}
	allUserInfo.GatewayProfixMobile = gateWayInfo.GatewayNum[0].NumberIndex

	
	allUserInfo.IsExist = 0

	userInfoByte, err := codec.MarshalJSONIndent(legacyQuerierCdc, allUserInfo)
	if err != nil {
		log.WithError(err).Info("GetUserByMobile Marshal Err")
		return nil, err
	}

	if userInfoByte == nil {
		log.Warning("GetRegisterInfo Not Fount:")
	}

	return userInfoByte, nil
}

func QueryUserInfos(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainChatQuery)
	var params types.QueryUserInfosParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	userInfos, err := k.GetUserInfos(ctx, params.Addresses)
	if err != nil {
		log.WithError(err).Info("GetUserInfos Err")
		return nil, err
	}

	userInfosByte, err := codec.MarshalJSONIndent(legacyQuerierCdc, userInfos)
	if err != nil {
		log.WithError(err).Info("GetUserInfos Marshal Err")
		return nil, err
	}

	if userInfosByte == nil {
		log.Warning("GetUserInfos Not Fount:")
	}

	return userInfosByte, nil
}

func QueryUserInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainChatQuery)
	var params types.QueryUserInfoParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	userInfo, err := k.GetRegisterInfo(ctx, params.Address)
	if err != nil {
		log.Error(params.Address + ": not found")
		log.WithError(err).Info("GetRegisterInfo Err")
		return nil, err
	}

	userInfoByte, err := codec.MarshalJSONIndent(legacyQuerierCdc, userInfo)
	if err != nil {
		log.WithError(err).Info("GetRegisterInfo Marshal Err")
		return nil, err
	}

	if userInfoByte == nil {
		log.Warning("GetRegisterInfo Not Fount:")
	}

	return userInfoByte, nil
}


func queryParams(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainChatQuery)

	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		log.WithError(err).Error("Chat Params MarshalJSONIndent")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
