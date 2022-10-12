package keeper

import (
	"context"
	"encoding/json"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/chat/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

var _ types.MsgServer = &Keeper{}

type msgServer struct {
	Keeper
	logPrefix string
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper, logPrefix: "chat | msgServer | "}
}

//
func (k Keeper) SetChatInfo(goCtx context.Context, msg *types.MsgSetChatInfo) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetFuncName(), core.LmChainChatKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)
	userChatInfo, err := k.GetRegisterInfo(ctx, msg.GetFromAddress())
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrUserNotFound
	}

	userChatInfo.ChatBlacklist = msg.ChatBlacklist
	userChatInfo.ChatWhitelist = msg.ChatWhitelist
	userChatInfo.ChatFee = msg.ChatFee
	userChatInfo.NodeAddress = msg.NodeAddress
	userChatInfo.ChatRestrictedMode = msg.ChatRestrictedMode
	userChatInfo.AddressBook = msg.AddressBook

	err = k.SetRegisterInfo(ctx, userChatInfo)
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrChatInfoSet
	}

	eventAddressBook, err := json.Marshal(msg.AddressBook)
	if err != nil {
		logs.Info("eventAddressBook Marshal err", err)
		return &types.MsgEmptyResponse{}, types.ErrChatInfoSet
	}

	eventChatBlacklist, err := json.Marshal(msg.ChatBlacklist)
	if err != nil {
		logs.Info("eventChatBlacklist Marshal err", err)
		return &types.MsgEmptyResponse{}, types.ErrChatInfoSet
	}

	eventChatWhitelist, err := json.Marshal(msg.ChatWhitelist)
	if err != nil {
		logs.Info("eventChatWhitelist Marshal err", err)
		return &types.MsgEmptyResponse{}, types.ErrChatInfoSet
	}

	//
	gateWayInfo, err := k.commKeeper.GetGatewayInfo(ctx, msg.NodeAddress)
	if err != nil {
		return nil, err
	}

	accFromAddress, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrAddressFormat
	}

	ctx.EventManager().EmitEvents(
		[]sdk.Event{
			sdk.NewEvent(
				types.TypeMsgSetChatInfo,
				sdk.NewAttribute(types.SetChatInfoEventTypeFromAddress, msg.FromAddress),
				sdk.NewAttribute(types.SetChatInfoEventTypeNodeAddress, msg.NodeAddress),
				sdk.NewAttribute(types.SetChatInfoEventTypeChatRestrictedMode, msg.ChatRestrictedMode),
				sdk.NewAttribute(types.SetChatInfoEventTypeChatFeeAmount, msg.ChatFee.Amount.String()),
				sdk.NewAttribute(types.SetChatInfoEventTypeChatFeeDenom, msg.ChatFee.Denom),
				sdk.NewAttribute(types.SetChatInfoEventTypeAddressBook, string(eventAddressBook)),
				sdk.NewAttribute(types.SetChatInfoEventTypeChatBlacklist, string(eventChatBlacklist)),
				sdk.NewAttribute(types.SetChatInfoEventTypeChatWhitelist, string(eventChatWhitelist)),
				sdk.NewAttribute(types.SetChatInfoEventTypeGatewayEventPrefixMobile, gateWayInfo.GatewayNum[0].NumberIndex),
				sdk.NewAttribute(types.SetChatInfoEventTypeGatewayEventFromBalance, k.bankKeeper.GetAllBalances(ctx, accFromAddress).String()),
			),
		},
	)

	return &types.MsgEmptyResponse{}, nil
}

//
func (k Keeper) BurnGetMobile(goCtx context.Context, msg *types.MsgBurnGetMobile) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetFuncName(), core.LmChainChatKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)
	userInfo, err := k.GetRegisterInfo(ctx, msg.GetFromAddress())
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrUserNotFound
	}

	chatParams := k.GetParams(ctx)

	if len(userInfo.Mobile) >= int(chatParams.MaxPhoneNumber) {
		return &types.MsgEmptyResponse{}, types.ErrUserMobileCount
	}

	// 

	//
	accFromAddress, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		logs.Info("accFromAddress err", err)
		return &types.MsgEmptyResponse{}, types.ErrAddressFormat
	}

	burnCoins := sdk.NewCoins(chatParams.DestroyPhoneNumberCoin)
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, accFromAddress, types.ModuleBurnName, burnCoins)
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrBurn
	}

	err = k.bankKeeper.BurnCoins(ctx, types.ModuleBurnName, burnCoins)
	if err != nil {
		logs.Info("BurnCoins err", err)
		return &types.MsgEmptyResponse{}, types.ErrBurn
	}

	//
	mobile, err := k.RegisterMobile(ctx, userInfo.NodeAddress, msg.FromAddress, msg.MobilePrefix)
	if err != nil {
		logs.Info("RegisterMobile err", err)
		return &types.MsgEmptyResponse{}, err
	}

	//
	userInfo.Mobile = append(userInfo.Mobile, mobile)
	err = k.SetRegisterInfo(ctx, userInfo)
	if err != nil {
		logs.Info("SetRegisterInfo err", err)
		return &types.MsgEmptyResponse{}, types.ErrUserUpdate
	}

	fromBlance := k.bankKeeper.GetAllBalances(ctx, accFromAddress).String()

	//
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.TypeMsgBurnGetMobile,
			sdk.NewAttribute(types.BurngetMobileEventTypeFromAddress, msg.FromAddress),
			sdk.NewAttribute(types.BurngetMobileEventTypeModuleAddress, k.accountKeeper.GetModuleAddress(types.ModuleName).String()),
			sdk.NewAttribute(types.BurngetMobileEventTypeAmount, chatParams.DestroyPhoneNumberCoin.Amount.String()),
			sdk.NewAttribute(types.BurngetMobileEventTypeDenom, chatParams.DestroyPhoneNumberCoin.Denom),
			sdk.NewAttribute(types.BurngetMobileEventTypeFromBalance, fromBlance),
		),
	)

	return &types.MsgEmptyResponse{}, nil
}

//
func (k Keeper) ChangeGateway(goCtx context.Context, msg *types.MsgChangeGateway) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	userInfo, err := k.GetRegisterInfo(ctx, msg.GetFromAddress())
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrUserNotFound
	}

	oldGateway := userInfo.NodeAddress

	userInfo.NodeAddress = msg.Gateway

	err = k.SetRegisterInfo(ctx, userInfo)
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrUserUpdate
	}

	//
	ctx.EventManager().EmitEvents(
		[]sdk.Event{
			sdk.NewEvent(
				types.TypeMsgChangeGateway,
				sdk.NewAttribute(types.ChangeGatewayEventTypeFromAddress, msg.FromAddress),
				sdk.NewAttribute(types.ChangeGatewayEventTypeOldGateWay, oldGateway),
				sdk.NewAttribute(types.ChangeGatewayEventTypeNewGateWay, msg.Gateway),
			),
		},
	)

	return &types.MsgEmptyResponse{}, nil
}

// MobileTransfer 
func (k Keeper) MobileTransfer(goCtx context.Context, msg *types.MsgMobileTransfer) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	//
	userInfo, err := k.GetRegisterInfo(ctx, msg.GetFromAddress())
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrUserNotFound
	}

	isHave := false

	for _, mobile := range userInfo.Mobile {
		if mobile == msg.Mobile {
			isHave = true
			break
		}
	}

	if !isHave {
		return &types.MsgEmptyResponse{}, types.ErrUserNotHaveMobile
	}

	//toAddress 
	toUserInfo, err := k.GetRegisterInfo(ctx, msg.GetToAddress())
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrUserNotFound
	}

	chatParams := k.GetParams(ctx)

	//10（）
	if len(toUserInfo.Mobile) >= int(chatParams.MaxPhoneNumber) {
		return &types.MsgEmptyResponse{}, types.ErrUserMobileCount
	}

	//
	userInfoNew := make([]string, 0)
	for _, mobile := range userInfo.Mobile {
		if mobile != msg.Mobile {
			userInfoNew = append(userInfoNew, mobile)
		}
	}
	userInfo.Mobile = userInfoNew
	err = k.SetRegisterInfo(ctx, userInfo)
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrUserUpdate
	}

	toUserInfo.Mobile = append(toUserInfo.Mobile, msg.Mobile)

	err = k.SetRegisterInfo(ctx, toUserInfo)
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrUserUpdate
	}

	return &types.MsgEmptyResponse{}, nil
}

//
func (k Keeper) AddressBookSave(goCtx context.Context, msg *types.MsgAddressBookSave) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	userInfo, err := k.GetRegisterInfo(ctx, msg.GetFromAddress())
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrUserNotFound
	}

	userInfo.Mobile = msg.AddressBook

	err = k.SetRegisterInfo(ctx, userInfo)
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrUserUpdate
	}

	return &types.MsgEmptyResponse{}, nil
}

//
func (k Keeper) SendGift(goCtx context.Context, msg *types.MsgSendGift) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//
	giftValueAll := sdk.NewCoin(msg.GiftValue.Denom, msg.GiftValue.Amount.Mul(sdk.NewInt(msg.GiftAmount)))

	//
	//mortgateInfo, err := k.MortgageSendCoin(ctx, types.TransferTypeToAccount, msg.ToAddress, msg.FromAddress, msg.NodeAddress, giftValueAll)
	//if err != nil {
	//    return &types.MsgEmptyResponse{}, err
	//}

	accFromAddress, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrDevideError
	}

	fromBalance := k.bankKeeper.GetAllBalances(ctx, accFromAddress)

	ctx.EventManager().EmitEvents(
		[]sdk.Event{
			sdk.NewEvent(
				types.TypeMsgSendGift,
				sdk.NewAttribute(types.SendGiftEventTypeFromAddress, msg.FromAddress),
				sdk.NewAttribute(types.SendGiftEventTypeToAddress, msg.ToAddress),
				sdk.NewAttribute(types.SendGiftEventTypeGateAddress, msg.NodeAddress),
				sdk.NewAttribute(types.SendGiftEventTypeGiftId, strconv.FormatInt(msg.GiftId, 10)),
				sdk.NewAttribute(types.SendGiftEventTypeGiftValue, msg.GiftValue.Amount.String()),
				sdk.NewAttribute(types.SendGiftEventTypeGiftDenom, msg.GiftValue.Denom),
				sdk.NewAttribute(types.SendGiftEventTypeGiftAmount, strconv.FormatInt(msg.GiftAmount, 10)),
				sdk.NewAttribute(types.SendGiftEventTypeGiftValueAll, giftValueAll.Amount.String()),
				//sdk.NewAttribute(types.SendGiftEventTypeGiftReceive, giftValueAll.Denom),
			),
			sdk.NewEvent(
				types.EventTypeDevide,
				sdk.NewAttribute(types.MortgageEventTypeType, types.EventTypeDevideSendGift),
				sdk.NewAttribute(types.MortgageEventTypeFromAddress, msg.FromAddress),
				sdk.NewAttribute(types.MortgageEventTypeDenom, msg.GiftValue.Denom),
				sdk.NewAttribute(types.MortgageEventTypeMortgageAmount, giftValueAll.Amount.String()),
				sdk.NewAttribute(types.MortgageEventTypeFromBalance, fromBalance.String()),
			),
		},
	)

	return &types.MsgEmptyResponse{}, nil
}

//
func (k Keeper) SetChatFee(goCtx context.Context, msg *types.MsgSetChatFee) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//todo 
	userInfo, err := k.GetRegisterInfo(ctx, msg.FromAddress)
	//if err != nil {
	//	return &types.MsgEmptyResponse{}, types.ErrUserNotFound
	//}

	//
	userInfo.ChatFee = msg.Fee

	err = k.SetRegisterInfo(ctx, userInfo)
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrUserUpdate
	}

	//
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.TypeMsgSetChatFee,
			sdk.NewAttribute(types.SetChatFeeEventTypeFromAddress, msg.FromAddress),
			sdk.NewAttribute(types.SetChatFeeEventTypeFee, msg.Fee.Amount.String()),
			sdk.NewAttribute(types.SetChatFeeEventTypeDenom, msg.Fee.Denom),
		),
	)

	return &types.MsgEmptyResponse{}, nil
}

//
func (k Keeper) Register(goCtx context.Context, msg *types.MsgRegister) (*types.MsgEmptyResponse, error) {

	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainKeeper)

	ctx := sdk.UnwrapSDKContext(goCtx)

	//
	_, err := k.GetRegisterInfo(ctx, msg.FromAddress)
	if err == nil {
		return &types.MsgEmptyResponse{}, types.ErrUserHasExisted
	}

	accFromAddress, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrDevideError
	}

	valAddr, valErr := sdk.ValAddressFromBech32(msg.NodeAddress)
	if valErr != nil {
		return nil, valErr
	}

	validator, found := k.pledgeKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, types.ErrValidatorNotFound
	}

	validatorAccAddr := sdk.AccAddress(valAddr)

	feeAllDec := sdk.ZeroDec()
	pledgeAmount := sdk.ZeroInt()
	if !msg.MortgageAmount.Amount.IsZero() {
		//
		feeAllDec, err := k.pledgeKeeper.ChatPoundage(ctx, accFromAddress, validatorAccAddr, msg.MortgageAmount)
		if err != nil {
			return &types.MsgEmptyResponse{}, err
		}

		//
		pledgeAmount = msg.MortgageAmount.Amount.Sub(feeAllDec.TruncateInt())
		//

		if !pledgeAmount.IsPositive() {
			return nil, types.ErrPledgeFeeSet
		}

		_, err = k.pledgeKeeper.Delegate(ctx, accFromAddress, accFromAddress, pledgeAmount, validator)
		if err != nil {
			log.Info("register  --> Delegate error:", err)
			return &types.MsgEmptyResponse{}, types.ErrPledgeDelegate
		}
	}

	// 
	mobile, err := k.RegisterMobile(ctx, msg.NodeAddress, msg.FromAddress, msg.MobilePrefix)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}

	userMobile := []string{
		mobile,
	}

	//
	var userInfo types.UserInfo
	userInfo.NodeAddress = msg.NodeAddress
	userInfo.Mobile = userMobile
	userInfo.FromAddress = msg.FromAddress

	//
	userInfo.ChatRestrictedMode = types.ChatRestrictedModeFee
	userInfo.ChatFee = types.DefaultParams().ChatFee

	err = k.SetRegisterInfo(ctx, userInfo)
	if err != nil {
		return &types.MsgEmptyResponse{}, types.ErrRegister
	}

	fromBalance := k.bankKeeper.GetAllBalances(ctx, accFromAddress)

	//
	moduleAccAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	//
	gateWayInfo, err := k.commKeeper.GetGatewayInfo(ctx, msg.NodeAddress)
	if err != nil {
		return nil, err
	}

	//
	ctx.EventManager().EmitEvents(
		[]sdk.Event{
			//
			sdk.NewEvent(
				types.TypeMsgRegister,
				//
				sdk.NewAttribute(types.EventTypeFromAddress, msg.FromAddress),
				//
				sdk.NewAttribute(types.EventTypeNodeAddress, msg.NodeAddress),
				//（）
				sdk.NewAttribute(types.EventTypeMortGageAmount, msg.MortgageAmount.Amount.String()),
				//
				sdk.NewAttribute(types.EventTypeMortGageDenom, msg.MortgageAmount.Denom),
				//
				sdk.NewAttribute(types.EventTypeGetMobile, mobile),
				//
				sdk.NewAttribute(types.EventTypePledgeAmount, pledgeAmount.String()),
				//
				sdk.NewAttribute(types.EventTypePledgeFee, feeAllDec.String()),
				//
				sdk.NewAttribute(types.EventTypeFromBalance, fromBalance.String()),
				//
				sdk.NewAttribute(types.EventPrefixMobile, gateWayInfo.GatewayNum[0].NumberIndex),
				// fee
				sdk.NewAttribute(types.EventDefaultChatRestrictedMode, userInfo.ChatRestrictedMode),
				//
				sdk.NewAttribute(types.SetChatInfoEventTypeChatFeeAmount, userInfo.ChatFee.Amount.String()),
				sdk.NewAttribute(types.SetChatInfoEventTypeChatFeeDenom, userInfo.ChatFee.Denom),
			),

			//
			sdk.NewEvent(
				types.EventTypeDevide,

				//
				sdk.NewAttribute(types.EventFeeAddress, core.ContractAddressFee.String()),
				//
				sdk.NewAttribute(types.EventTypeDevideType, types.EventTypeDevideRegister),
				//
				sdk.NewAttribute(types.DevideEventFromAddress, msg.FromAddress),
				//
				sdk.NewAttribute(types.DevideEventToAddress, moduleAccAddr.String()),
				//
				sdk.NewAttribute(types.DevideEventAmount, feeAllDec.String()),
				//
				sdk.NewAttribute(types.DevideEventDenom, msg.MortgageAmount.Denom),
				//
				sdk.NewAttribute(types.DevideEventBalance, fromBalance.String()),
			),
		},
	)

	return &types.MsgEmptyResponse{}, nil
}
