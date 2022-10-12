package keeper

import (
	"context"
	"freemasonry.cc/blockchain/core"
	chatTypes "freemasonry.cc/blockchain/x/chat/types"
	"freemasonry.cc/blockchain/x/pledge/types"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	Keeper
	logPrefix string
}

//
func (s msgServer) WithdrawDelegatorReward(goCtx context.Context, msg *types.MsgWithdrawDelegatorReward) (*types.MsgWithdrawDelegatorRewardResponse, error) {
	logs := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)

	//todo ，
	accAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, types.ErrAddressFormat
	}

	allChatPledge, err := s.GetAllChatPledge(ctx, accAddress)
	if err != nil {
		return nil, types.ErrGetAllPledge
	}

	if len(allChatPledge) == 0 {
		return nil, nil
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	allRewardCoins := sdk.NewCoins()
	for _, resp := range allChatPledge {

		valAddr, err := sdk.ValAddressFromBech32(resp.Delegation.ValidatorAddress)
		if err != nil {
			logs.Info("WithdrawDelegatorReward", "err1", err)
			return nil, err
		}

		reward, err := s.WithdrawDelegationRewards(ctx, delegatorAddress, valAddr)
		if err != nil {
			logs.Info("WithdrawDelegatorReward", "err2", err)
			return nil, err
		}
		reward = reward.Sort()

		allRewardCoins = allRewardCoins.Add(reward...).Sort()

	}

	//
	logs.Info("------")
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.ChatWithDrawAllEventType,

			//
			sdk.NewAttribute(types.ChatWithDrawAllEventType, msg.DelegatorAddress),
			//
			sdk.NewAttribute(types.ChatWithDrawAllEventTypeFromAddress, msg.DelegatorAddress),
			//
			sdk.NewAttribute(types.ChatWithDrawAllEventTypeFromBalance, s.BankKeeper.GetAllBalances(ctx, delegatorAddress).String()),
		),
	})

	return &types.MsgWithdrawDelegatorRewardResponse{}, nil
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper, logPrefix: "pledge | msgServer | "}
}

//
func (s msgServer) Delegate(goCtx context.Context, msg *types.MsgPledge) (*types.MsgEmptyResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	ctx := sdk.UnwrapSDKContext(goCtx)
	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if valErr != nil {
		return nil, valErr
	}

	validator, found := s.GetValidator(ctx, valAddr)
	if !found {
		return nil, types.ErrNoValidatorFound
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	fromAddress, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}

	bondDenom := s.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}

	gatewayInfo, err := s.CommKeeper.GetGatewayInfo(ctx, validator.OperatorAddress)
	if err != nil {
		return nil, types.ErrGetGatewayInfo
	}

	if gatewayInfo.Status != 0 {
		if delegatorAddress.String() != validator.OperatorAddress {
			return nil, types.ErrValidatotStatus
		}
	}

	//
	//
	kParams := s.GetParams(ctx)
	minPledgeCoin := kParams.MinMortgageCoin

	//
	delegations := s.GetAllDelegatorDelegations(ctx, delegatorAddress)
	delegationResps, err := DelegationsToDelegationResponses(ctx, s.Keeper, delegations)
	if err != nil {
		return nil, err
	}
	pledgeTotal := sdk.ZeroInt()
	if len(delegationResps) > 0 {
		for _, resp := range delegationResps {
			pledgeTotal = pledgeTotal.Add(resp.Balance.Amount)
		}
	}
	if pledgeTotal.Add(msg.Amount.Amount).LT(minPledgeCoin.Amount) {
		return nil, types.ErrPledgeLimit
	}

	//
	validatorAccAddr := sdk.AccAddress(valAddr)

	feeAllDec, err := s.ChatPoundage(ctx, delegatorAddress, validatorAccAddr, msg.Amount)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}

	//
	pledgeAmount := msg.Amount.Amount.Sub(feeAllDec.TruncateInt())

	// NOTE: source funds are always unbonded
	_, err = s.Keeper.Delegate(ctx, fromAddress, delegatorAddress, pledgeAmount, validator)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, types.ModuleName, "delegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", msg.Type()},
				float32(msg.Amount.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			chatTypes.EventTypeDevide,

			sdk.NewAttribute(chatTypes.EventFeeAddress, core.ContractAddressFee.String()),
			//
			sdk.NewAttribute(chatTypes.EventTypeDevideType, chatTypes.EventTypeDevidePledge),
			//
			sdk.NewAttribute(chatTypes.DevideEventFromAddress, msg.FromAddress),
			//
			sdk.NewAttribute(chatTypes.DevideEventToAddress, s.AccountKeeper.GetModuleAddress(types.ModuleName).String()),
			//
			sdk.NewAttribute(chatTypes.DevideEventAmount, feeAllDec.TruncateInt().String()),
			//
			sdk.NewAttribute(chatTypes.DevideEventDenom, msg.Amount.Denom),
			//
			sdk.NewAttribute(chatTypes.DevideEventBalance, s.BankKeeper.GetAllBalances(ctx, delegatorAddress).String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	})

	return &types.MsgEmptyResponse{}, nil
}

//
func (s msgServer) Undelegate(goCtx context.Context, msg *types.MsgUnpledge) (*types.MsgEmptyResponse, error) {
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
	validatorAccAddr := sdk.AccAddress(msg.ValidatorAddress)

	feeAllDec, err := s.ChatPoundage(ctx, delegatorAddress, validatorAccAddr, msg.Amount)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}

	//
	unPledgeAmount := msg.Amount.Amount.Sub(feeAllDec.TruncateInt())

	shares, err := s.ValidateUnbondAmount(
		ctx, delegatorAddress, addr, unPledgeAmount,
	)
	if err != nil {
		return nil, err
	}

	bondDenom := s.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, bondDenom,
		)
	}

	undelegateCoin, err := s.Keeper.Undelegate(ctx, delegatorAddress, addr, shares)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{

		//
		sdk.NewEvent(
			chatTypes.EventTypeDevide,

			sdk.NewAttribute(chatTypes.EventFeeAddress, core.ContractAddressFee.String()),
			//
			sdk.NewAttribute(chatTypes.EventTypeDevideType, chatTypes.EventTypeDevideRedeem),
			//
			sdk.NewAttribute(chatTypes.DevideEventFromAddress, msg.DelegatorAddress),
			//
			sdk.NewAttribute(chatTypes.DevideEventToAddress, s.AccountKeeper.GetModuleAddress(types.ModuleName).String()),
			//
			sdk.NewAttribute(chatTypes.DevideEventAmount, feeAllDec.TruncateInt().String()),
			//
			sdk.NewAttribute(chatTypes.DevideEventDenom, msg.Amount.Denom),
			//
			sdk.NewAttribute(chatTypes.DevideEventBalance, s.BankKeeper.GetAllBalances(ctx, delegatorAddress).String()),
		),

		//
		sdk.NewEvent(
			types.EventTypeUnbond,
			sdk.NewAttribute(types.EventUnPledgeFromAddress, s.AccountKeeper.GetModuleAddress(types.BondedPoolName).String()),
			sdk.NewAttribute(types.EventUnPledgeToAddress, msg.DelegatorAddress),
			sdk.NewAttribute(types.EventUnPledgeAmount, undelegateCoin.Amount.String()),
			sdk.NewAttribute(types.EventUnPledgeDenom, s.BondDenom(ctx)),
			sdk.NewAttribute(types.EventUnPledgeFromBalance, s.BankKeeper.GetAllBalances(ctx, delegatorAddress).String()),
		),
	})

	return &types.MsgEmptyResponse{}, nil
}
