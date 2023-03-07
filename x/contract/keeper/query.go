package keeper

import (
    "freemasonry.cc/blockchain/core"
    "freemasonry.cc/blockchain/util"
    "freemasonry.cc/blockchain/x/contract/types"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    "github.com/ethereum/go-ethereum/common"
    evmtypes "github.com/evmos/ethermint/x/evm/types"
    abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
    return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
        var (
            res []byte
            err error
        )
        switch path[0] {
        case types.QueryParams:
            return queryParams(ctx, req, k, legacyQuerierCdc)
        case types.QueryNft:
            return queryNftInfo(ctx, req, k, legacyQuerierCdc)
        case types.QueryNftContractAddress:
            return queryContractAddress(ctx, k)
        case types.QueryContractCode:
            return queryContractCode(ctx, req, k, legacyQuerierCdc)
        default:
            err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
        }

        return res, err
    }
}

func queryContractCode(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
    log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
    var params evmtypes.QueryCodeRequest
    err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
    if err != nil {
        log.WithError(err).Error("UnmarshalJSON")
        return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
    }
    addr := common.HexToAddress(params.Address)
    acct := k.evmKeeper.GetAccountWithoutBalance(ctx, addr)
    var code []byte
    if acct != nil && acct.IsContract() {
        code = k.evmKeeper.GetCode(ctx, common.BytesToHash(acct.CodeHash))
    }
    return code, nil
}

func queryContractAddress(ctx sdk.Context, k Keeper) ([]byte, error) {
    contract := k.GetNftContractAddress(ctx)
    return []byte(contract), nil
}

func queryNftInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
    log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
    var params types.QueryNftInfoParams
    err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
    if err != nil {
        log.WithError(err).Error("UnmarshalJSON")
        return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
    }

    addr, _ := sdk.AccAddressFromBech32(params.Address)
    from := common.BytesToAddress(addr)
    contract := common.HexToAddress(params.ContractAddress)
    info, err := k.GetNftInfo(ctx, from, contract)
    if err != nil {
        return nil, err
    }
    infoByte, err := util.Json.Marshal(info)
    if err != nil {
        log.WithError(err).Error("Marshal")
        return nil, err
    }
    return infoByte, nil
}


func queryParams(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
    params := k.GetParams(ctx)
    res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
    }

    return res, nil
}
