package keeper

import (
	"context"
	"encoding/base64"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/comm/types"
	"github.com/armon/go-metrics"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"time"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	Keeper
	logPrefix string
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper, logPrefix: "comm | msgServer | "}
}

//
func (k msgServer) CreateSmartValidator(goCtx context.Context, msg *types.MsgCreateSmartValidator) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	valAddress := sdk.ValAddress(addr)
	_, found := k.stakingKeeper.GetValidator(ctx, valAddress)
	if found {
		return nil, stakingTypes.ErrValidatorOwnerExists
	}
	err = k.createValidator(ctx, addr, valAddress, *msg, msg.Value)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	return &types.MsgEmptyResponse{}, nil
}

//
func (k msgServer) GatewayRegister(goCtx context.Context, msg *types.MsgGatewayRegister) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	if msg.Delegation == "" {
		msg.Delegation = "0"
	}
	delInt, ok := sdk.NewIntFromString(msg.Delegation)
	if !ok {
		return &types.MsgEmptyResponse{}, sdkerrors.Wrapf(
			types.ErrDelegationCoin, "invalid delegation : got %s", msg.Delegation)
	}
	params := k.GetParams(ctx)
	delegation := sdk.NewCoin(sdk.DefaultBondDenom, delInt)

	valAddress := sdk.ValAddress(addr)
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddress)
	if !found { //,
		return nil, types.ErrValidatorNotFound
	}
	//
	if !delegation.IsZero() {
		err = k.delegate(ctx, addr, valAddress, validator, delegation)
		if err != nil {
			return nil, err
		}
	}
	//
	delegate, found := k.stakingKeeper.GetDelegation(ctx, addr, valAddress)
	if !found {
		return nil, stakingTypes.ErrNoDelegation
	}
	//  
	if delegate.Shares.LT(params.MinDelegate.ToDec()) {
		return &types.MsgEmptyResponse{}, types.ErrGatewayDelegation
	}
	delegation = sdk.NewCoin(sdk.DefaultBondDenom, delegate.Shares.TruncateInt())
	//
	err = k.SetGateway(ctx, *msg, delegation, valAddress.String())
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	//
	k.hooks.AfterCreateGateway(ctx, validator)

	return &types.MsgEmptyResponse{}, nil
}

//
func (k msgServer) GatewayIndexNum(goCtx context.Context, msg *types.MsgGatewayIndexNum) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainCommKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if valErr != nil {
		return nil, valErr
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingTypes.ErrNoValidatorFound
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	//，
	if validator.GetOperator().String() == sdk.ValAddress(delegatorAddress).String() {
		//
		gateway, err := k.GetGatewayInfo(ctx, msg.ValidatorAddress)
		if err != nil {
			if err == types.ErrGatewayNotExist {
				return &types.MsgEmptyResponse{}, nil
			}
			return nil, err
		}
		//
		delegation, found := k.stakingKeeper.GetDelegation(ctx, delegatorAddress, valAddr)
		if !found {
			return nil, stakingTypes.ErrNoDelegation
		}
		params := k.GetParams(ctx)
		num := delegation.Shares.QuoInt(params.MinDelegate)
		gateway.GatewayQuota = num.TruncateInt64()
		if msg.IndexNumber != nil && len(msg.IndexNumber) > 0 {
			//
			if num.Sub(sdk.NewInt(int64(len(gateway.GatewayNum))).ToDec()).LT(sdk.NewDec(int64(len(msg.IndexNumber)))) {
				log.WithFields(logrus.Fields{"num": num, "GatewayNum": len(gateway.GatewayNum), "IndexNumber": len(msg.IndexNumber)}).Error(types.ErrGatewayNum)
				return nil, types.ErrGatewayNum
			}
			indexNumArray, err := k.GatewayNumFilter(ctx, msg.ValidatorAddress, msg.IndexNumber)
			if err != nil {
				return nil, err
			}
			//
			err = k.SetGatewayNum(ctx, indexNumArray)
			if err != nil {
				return nil, err
			}
			//,
			err = k.GatewayRedeemNumFilter(ctx, indexNumArray)
			if err != nil {
				return nil, err
			}
			gateway.GatewayNum = append(gateway.GatewayNum, indexNumArray...)
			//
			gateway.Status = 0
		}
		//
		err = k.UpdateGatewayInfo(ctx, *gateway)
		if err != nil {
			return nil, err
		}
	}
	return &types.MsgEmptyResponse{}, nil
}

//
func (k msgServer) GatewayUndelegate(goCtx context.Context, msg *types.MsgGatewayUndelegate) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	//
	shares, err := k.stakingKeeper.ValidateUnbondAmount(
		ctx, delegatorAddress, addr, msg.Amount.Amount,
	)
	if err != nil {
		return nil, err
	}
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, addr)
	if !found {
		return nil, stakingTypes.ErrNoDelegatorForAddress
	}
	completionTime, returnAmount, err := k.Keeper.Undelegate(ctx, delegatorAddress, addr, validator, shares)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingTypes.EventTypeUnbond,
			sdk.NewAttribute(stakingTypes.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),    //
			sdk.NewAttribute(types.AttributeKeyReturnAmount, returnAmount.String()), //
			sdk.NewAttribute(stakingTypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),

			sdk.NewAttribute(stakingTypes.AttributeKeyDelegatorAddr, delegatorAddress.String()), //
			sdk.NewAttribute(stakingTypes.AttributeKeyNewShares, shares.String()),               //
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, stakingTypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	})

	//，
	err = k.GatewayNumUnbond(ctx, delegatorAddress, validator.GetOperator(), msg.IndexNumber, shares)
	if err != nil {
		return nil, err
	}
	return &types.MsgEmptyResponse{}, nil
}

//
func (k msgServer) GatewayBeginRedelegate(goCtx context.Context, msg *types.MsgGatewayBeginRedelegate) (*types.MsgBeginRedelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valSrcAddr, err := sdk.ValAddressFromBech32(msg.ValidatorSrcAddress)
	if err != nil {
		return nil, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	validator, found := k.stakingKeeper.GetValidator(ctx, valSrcAddr)
	if !found {
		return nil, stakingTypes.ErrNoValidatorFound
	}

	shares, err := k.stakingKeeper.ValidateUnbondAmount(
		ctx, delegatorAddress, valSrcAddr, msg.Amount.Amount,
	)
	if err != nil {
		return nil, err
	}

	bondDenom := k.stakingKeeper.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}

	valDstAddr, err := sdk.ValAddressFromBech32(msg.ValidatorDstAddress)
	if err != nil {
		return nil, err
	}

	completionTime, err := k.stakingKeeper.BeginRedelegation(
		ctx, delegatorAddress, valSrcAddr, valDstAddr, shares,
	)
	if err != nil {
		return nil, err
	}

	//，
	err = k.GatewayNumUnbond(ctx, delegatorAddress, validator.GetOperator(), msg.IndexNumber, shares)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, types.ModuleName, "redelegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", msg.Type()},
				float32(msg.Amount.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingTypes.EventTypeRedelegate,
			sdk.NewAttribute(stakingTypes.AttributeKeySrcValidator, msg.ValidatorSrcAddress),
			sdk.NewAttribute(stakingTypes.AttributeKeyDstValidator, msg.ValidatorDstAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(stakingTypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, stakingTypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	})
	return &types.MsgBeginRedelegateResponse{
		CompletionTime: completionTime,
	}, nil
}

//base64bech32
func ParseBech32ValConsPubkey(validatorInfoPubKeyBase64 string) (cryptotypes.PubKey, error) {
	validatorInfoPubKeyBytes, err := base64.StdEncoding.DecodeString(validatorInfoPubKeyBase64)
	if err != nil {
		return nil, err
	}
	pbk := ed25519.PubKey(validatorInfoPubKeyBytes) //ed25519
	pubkey, err := cryptocodec.FromTmPubKeyInterface(pbk)
	if err != nil {
		return nil, err
	}
	return pubkey, nil
}
