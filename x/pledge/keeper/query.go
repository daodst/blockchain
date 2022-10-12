package keeper

import (
    // this line is used by starport scaffolding # 1

    "freemasonry.cc/blockchain/x/pledge/types"
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
        case types.QueryDelegation: //
            return queryDelegation(ctx, req, k, legacyQuerierCdc)
        case types.QueryDelegatorDelegations:
            return queryDelegatorDelegations(ctx, req, k, legacyQuerierCdc)
        case types.QueryParams:
            return queryParams(ctx, req, k, legacyQuerierCdc)
        case types.QueryPrePledge:
            return queryPrePledge(ctx, req, k, legacyQuerierCdc)
        default:
            err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
        }

        return res, err
    }
}

func QueryUserInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

    return nil, nil
}

func queryDelegatorDelegations(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
    var params types.QueryDelegatorParams

    err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
    }

    delegations := k.GetAllDelegatorDelegations(ctx, params.DelegatorAddr)
    delegationResps, err := DelegationsToDelegationResponses(ctx, k, delegations)

    if err != nil {
        return nil, err
    }

    if delegationResps == nil {
        delegationResps = types.DelegationResponses{}
    }

    res, err := codec.MarshalJSONIndent(legacyQuerierCdc, delegationResps)

    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
    }

    return res, nil
}

func DelegationsToDelegationResponses(
    ctx sdk.Context, k Keeper, delegations types.Delegations,
) (types.DelegationResponses, error) {
    resp := make(types.DelegationResponses, len(delegations))

    for i, del := range delegations {
        delResp, err := DelegationToDelegationResponse(ctx, k, del)
        if err != nil {
            return nil, err
        }

        resp[i] = delResp
    }

    return resp, nil
}

// return all delegations for a delegator
func (k Keeper) GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []types.Delegation {
    delegations := make([]types.Delegation, 0)

    store := ctx.KVStore(k.storeKey)
    delegatorPrefixKey := types.GetDelegationsKey(delegator)

    iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
    defer iterator.Close()

    i := 0

    for ; iterator.Valid(); iterator.Next() {
        delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
        delegations = append(delegations, delegation)
        i++
    }

    return delegations
}

func queryDelegation(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
    var params types.QueryDelegatorValidatorRequest

    err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
    }

    delAddr, err := sdk.AccAddressFromBech32(params.DelegatorAddr)
    if err != nil {
        return nil, err
    }

    valAddr, err := sdk.ValAddressFromBech32(params.ValidatorAddr)
    if err != nil {
        return nil, err
    }

    delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
    if !found {
        return nil, types.ErrNoDelegation
    }

    delegationResp, err := DelegationToDelegationResponse(ctx, k, delegation)
    if err != nil {
        return nil, err
    }

    res, err := codec.MarshalJSONIndent(legacyQuerierCdc, delegationResp)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
    }

    return res, nil
}

func DelegationToDelegationResponse(ctx sdk.Context, k Keeper, del types.Delegation) (types.DelegationResponse, error) {
    val, found := k.GetValidator(ctx, del.GetValidatorAddr())
    if !found {
        return types.DelegationResponse{}, types.ErrNoValidatorFound
    }

    delegatorAddress, err := sdk.AccAddressFromBech32(del.DelegatorAddress)
    if err != nil {
        return types.DelegationResponse{}, err
    }

    return types.NewDelegationResp(
        delegatorAddress,
        del.GetValidatorAddr(),
        del.Shares,
        sdk.NewCoin(k.BondDenom(ctx), val.TokensFromShares(del.Shares).TruncateInt()),
    ), nil
}

func queryParams(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
    params := k.GetParams(ctx)

    res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
    }

    return res, nil
}

func queryPrePledge(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

    var address string

    err := legacyQuerierCdc.UnmarshalJSON(req.Data, &address)

    prePledgeAmount, err := k.GetPledgeDelegate(ctx, address)
    if err != nil {
        return nil, err
    }

    res, err := codec.MarshalJSONIndent(legacyQuerierCdc, prePledgeAmount)
    if err != nil {
        return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
    }

    return res, nil
}
