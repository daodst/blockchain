package keeper

import (
	"context"
	"encoding/json"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/chat/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	Keeper
	logPrefix string
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper, logPrefix: "chat | msgServer | "}
}


func (m msgServer) SetChatInfo(goCtx context.Context, msg *types.MsgSetChatInfo) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(m), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	userChatInfo, err := m.GetRegisterInfo(ctx, msg.GetFromAddress())
	if err != nil {
		logs.WithError(err).Error("GetRegisterInfo err")
		return &types.MsgEmptyResponse{}, err
	}
	var nodeFlag string

	
	if msg.GatewayAddress != userChatInfo.NodeAddress {
		nodeFlag = "1"
	}

	userChatInfo.ChatBlacklist = msg.ChatBlacklist
	userChatInfo.ChatWhitelist = msg.ChatWhitelist
	userChatInfo.NodeAddress = msg.GatewayAddress
	userChatInfo.AddressBook = msg.AddressBook
	userChatInfo.UpdateTime = msg.UpdateTime
	userChatInfo.ChatBlackEncList = msg.ChatBlacklistEnc
	userChatInfo.ChatWhiteEncList = msg.ChatWhitelistEnc

	err = m.SetRegisterInfo(ctx, userChatInfo)
	if err != nil {
		logs.WithError(err).Error("SetRegisterInfo err")
		return &types.MsgEmptyResponse{}, core.ErrChatInfoSet
	}

	eventAddressBook, err := json.Marshal(msg.AddressBook)
	if err != nil {
		logs.WithError(err).Error("eventAddressBook Marshal err", err)
		return &types.MsgEmptyResponse{}, core.ErrChatInfoSet
	}

	eventChatBlacklist, err := json.Marshal(msg.ChatBlacklist)
	if err != nil {
		logs.WithError(err).Error("eventChatBlacklist Marshal err", err)
		return &types.MsgEmptyResponse{}, core.ErrChatInfoSet
	}

	eventChatWhitelist, err := json.Marshal(msg.ChatWhitelist)
	if err != nil {
		logs.WithError(err).Error("eventChatWhitelist Marshal err", err)
		return &types.MsgEmptyResponse{}, core.ErrChatInfoSet
	}

	
	var gatewayNumberIndex string
	if msg.GatewayAddress != "" {
		gateWayInfo, err := m.commKeeper.GetGatewayInfo(ctx, msg.GatewayAddress)
		if err != nil {
			logs.WithError(err).Error("GetGatewayInfo err")
			return nil, err
		}
		if gateWayInfo.Status == 1 {
			return nil, core.ErrNumberOfGateWay
		}
		gatewayNumberIndex = gateWayInfo.GatewayNum[0].NumberIndex
	}

	accFromAddress, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		logs.WithError(err).Error("accFromAddress err")
		return &types.MsgEmptyResponse{}, core.ErrAddressFormat
	}

	ctx.EventManager().EmitEvents(
		[]sdk.Event{
			sdk.NewEvent(
				types.TypeMsgSetChatInfo,
				sdk.NewAttribute(types.SetChatInfoEventTypeFromAddress, msg.FromAddress),
				sdk.NewAttribute(types.SetChatInfoEventTypeNodeAddress, msg.GatewayAddress),
				sdk.NewAttribute(types.SetChatInfoEventTypeNodeFlag, nodeFlag),
				sdk.NewAttribute(types.SetChatInfoEventTypeAddressBook, string(eventAddressBook)),
				sdk.NewAttribute(types.SetChatInfoEventTypeChatBlacklist, string(eventChatBlacklist)),
				sdk.NewAttribute(types.SetChatInfoEventTypeChatWhitelist, string(eventChatWhitelist)),
				sdk.NewAttribute(types.SetChatInfoEventTypeGatewayEventPrefixMobile, gatewayNumberIndex),
				sdk.NewAttribute(types.SetChatInfoEventTypeGatewayEventFromBalance, m.bankKeeper.GetAllBalances(ctx, accFromAddress).String()),
			),
		},
	)

	return &types.MsgEmptyResponse{}, nil
}


func (m msgServer) BurnGetMobile(goCtx context.Context, msg *types.MsgBurnGetMobile) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(m), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)

	
	if msg.GatewayAddress != "" && msg.ChatAddress != "" {
		err := m.Register(ctx, msg.FromAddress, msg.ChatAddress, msg.GatewayAddress, nil)
		if err != nil {
			return nil, err
		}
	}

	userInfo, err := m.GetRegisterInfo(ctx, msg.GetFromAddress())
	if err != nil {
		logs.WithError(err).Error("GetRegisterInfo err")
		return &types.MsgEmptyResponse{}, core.ErrUserNotFound
	}

	chatParams := m.GetParams(ctx)

	if len(userInfo.Mobile) >= int(chatParams.MaxPhoneNumber) {
		return &types.MsgEmptyResponse{}, core.ErrUserMobileCount
	}

	

	
	accFromAddress, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		logs.WithError(err).Error("accFromAddress err")
		return &types.MsgEmptyResponse{}, core.ErrAddressFormat
	}

	burnCoins := sdk.NewCoins(chatParams.DestroyPhoneNumberCoin)
	err = m.bankKeeper.SendCoinsFromAccountToModule(ctx, accFromAddress, types.ModuleBurnName, burnCoins)
	if err != nil {
		logs.WithError(err).Error("SendCoinsFromAccountToModule err")
		return &types.MsgEmptyResponse{}, core.ErrBurn
	}

	err = m.bankKeeper.BurnCoins(ctx, types.ModuleBurnName, burnCoins)
	if err != nil {
		logs.WithError(err).Error("BurnCoins err")
		return &types.MsgEmptyResponse{}, core.ErrBurn
	}

	
	mobile, err := m.RegisterMobile(ctx, userInfo.NodeAddress, msg.FromAddress, msg.MobilePrefix)
	if err != nil {
		logs.WithError(err).Error("RegisterMobile err")
		return &types.MsgEmptyResponse{}, err
	}

	
	userInfo.Mobile = append(userInfo.Mobile, mobile)
	err = m.SetRegisterInfo(ctx, userInfo)
	if err != nil {
		logs.WithError(err).Error("SetRegisterInfo err")
		return &types.MsgEmptyResponse{}, core.ErrUserUpdate
	}

	fromBlance := m.bankKeeper.GetAllBalances(ctx, accFromAddress).String()

	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.TypeMsgBurnGetMobile,
			sdk.NewAttribute(types.BurngetMobileEventTypeFromAddress, msg.FromAddress),
			sdk.NewAttribute(types.BurngetMobileEventTypeModuleAddress, m.accountKeeper.GetModuleAddress(types.ModuleName).String()),
			sdk.NewAttribute(types.BurngetMobileEventTypeAmount, chatParams.DestroyPhoneNumberCoin.Amount.String()),
			sdk.NewAttribute(types.BurngetMobileEventTypeDenom, chatParams.DestroyPhoneNumberCoin.Denom),
			sdk.NewAttribute(types.BurngetMobileEventTypeFromBalance, fromBlance),
			sdk.NewAttribute(types.BurngetMobileEventTypeGetMobile, mobile),
		),
	)

	return &types.MsgEmptyResponse{}, nil
}

// MobileTransfer 
func (m msgServer) MobileTransfer(goCtx context.Context, msg *types.MsgMobileTransfer) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(m), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	userInfo, err := m.GetRegisterInfo(ctx, msg.GetFromAddress())
	if err != nil {
		logs.WithError(err).Error("GetRegisterInfo err")
		return &types.MsgEmptyResponse{}, core.ErrUserNotFound
	}

	
	if len(userInfo.Mobile) == 0 || len(userInfo.Mobile) == 1 {
		return &types.MsgEmptyResponse{}, core.ErrUserMobileCount
	}

	isHave := false

	for _, mobile := range userInfo.Mobile {
		if mobile == msg.Mobile {
			isHave = true
			break
		}
	}

	
	if !isHave {
		return &types.MsgEmptyResponse{}, core.ErrUserNotHaveMobile
	}

	
	if msg.FromAddress == msg.ToAddress {
		return nil, core.ErrMobileTransferTo
	}

	//toAddress 
	toUserInfo, err := m.GetRegisterInfo(ctx, msg.GetToAddress())
	if err != nil {
		logs.WithError(err).Error("GetToAddressRegisterInfo err")
		return &types.MsgEmptyResponse{}, core.ErrUserNotFound
	}

	chatParams := m.GetParams(ctx)

	//10（）
	if len(toUserInfo.Mobile) >= int(chatParams.MaxPhoneNumber) {
		return &types.MsgEmptyResponse{}, core.ErrUserMobileCount
	}

	
	userInfoNew := make([]string, 0)
	for _, mobile := range userInfo.Mobile {
		if mobile != msg.Mobile {
			userInfoNew = append(userInfoNew, mobile)
		}
	}
	userInfo.Mobile = userInfoNew
	err = m.SetRegisterInfo(ctx, userInfo)
	if err != nil {
		logs.WithError(err).Error("SetRegisterInfo err")
		return &types.MsgEmptyResponse{}, core.ErrUserUpdate
	}

	toUserInfo.Mobile = append(toUserInfo.Mobile, msg.Mobile)

	err = m.SetRegisterInfo(ctx, toUserInfo)
	if err != nil {
		logs.WithError(err).Error("SetToAddressRegisterInfo err")
		return &types.MsgEmptyResponse{}, core.ErrUserUpdate
	}

	
	err = m.SetMobileOwner(ctx, msg.Mobile, msg.ToAddress)
	if err != nil {
		logs.WithError(err).Info("Err SetMobileOwner")
		return &types.MsgEmptyResponse{}, err
	}

	accFromAddress, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		logs.WithError(err).Error("AccAddressFromHex err")
		return &types.MsgEmptyResponse{}, core.ErrAddressFormat
	}

	fromBalance := m.bankKeeper.GetAllBalances(ctx, accFromAddress).String()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.TypeMsgMobileTransfer,
			sdk.NewAttribute(types.ChatMobileTransferEventTypeFromAddress, msg.FromAddress),
			sdk.NewAttribute(types.ChatMobileTransferEventTypeToAddress, msg.ToAddress),
			sdk.NewAttribute(types.ChatMobileTransferEventTypeMobile, msg.Mobile),
			sdk.NewAttribute(types.ChatMobileTransferEventTypeFromBalance, fromBalance),
		),
	)

	return &types.MsgEmptyResponse{}, nil
}
