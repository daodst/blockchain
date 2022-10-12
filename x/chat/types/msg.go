package types

import (
	"freemasonry.cc/blockchain/cmd/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	types "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgRegister{}
	_ sdk.Msg = &MsgSendGift{}
	_ sdk.Msg = &MsgSetChatFee{}
	_ sdk.Msg = &MsgAddressBookSave{}
	_ sdk.Msg = &MsgSetChatInfo{}
)

const (
	TypeMsgRegister        = "register"
	TypeMsgSetChatFee      = "set_chat_fee"
	TypeMsgSendGift        = "send_gift"
	TypeMsgAddressBookSave = "address_book_save"
	TypeMsgMobileTransfer  = "mobile_transfer"
	TypeMsgChangeGateway   = "change_gateway"
	TypeMsgBurnGetMobile   = "burn_get_mobile"
	TypeMsgSetChatInfo     = "set_chat_info"
)

//
func NewMsgSetChatInfo(fromAddress, nodeAddress, chatRestrictedMode string, addressBook, chatBlacklist, chatWhitelist []string, chatFee types.Coin) *MsgSetChatInfo {
	return &MsgSetChatInfo{
		FromAddress:        fromAddress,
		NodeAddress:        nodeAddress,
		AddressBook:        addressBook,
		ChatBlacklist:      chatBlacklist,
		ChatRestrictedMode: chatRestrictedMode,
		ChatFee:            chatFee,
		ChatWhitelist:      chatWhitelist,
	}
}

func (msg MsgSetChatInfo) Route() string { return RouterKey }
func (msg MsgSetChatInfo) Type() string  { return TypeMsgSetChatInfo }
func (msg MsgSetChatInfo) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgSetChatInfo) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgSetChatInfo) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}

	_, err = sdk.ValAddressFromBech32(msg.NodeAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}

	if msg.ChatFee.Denom != config.BaseDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "coin error")
	}

	if msg.ChatFee.Amount.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "amount error")
	}
	//fee | any | list
	if msg.ChatRestrictedMode != ChatRestrictedModeFee && msg.ChatRestrictedMode != ChatRestrictedModeAny && msg.ChatRestrictedMode != ChatRestrictedModeList {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "chat restricted mode error")
	}

	return nil
}
func (m MsgSetChatInfo) XXX_MessageName() string {
	return TypeMsgSetChatInfo
}

// NewMsgBurnGetMobile **************************/
func NewMsgBurnGetMobile(fromAddress, mobilePrefix string) *MsgBurnGetMobile {
	return &MsgBurnGetMobile{
		FromAddress:  fromAddress,
		MobilePrefix: mobilePrefix,
	}
}

func (msg MsgBurnGetMobile) Route() string { return RouterKey }
func (msg MsgBurnGetMobile) Type() string  { return TypeMsgBurnGetMobile }
func (msg MsgBurnGetMobile) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgBurnGetMobile) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgBurnGetMobile) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}

	return nil
}
func (m MsgBurnGetMobile) XXX_MessageName() string {
	return TypeMsgBurnGetMobile
}

// NewMsgChangeGateway **************************/
func NewMsgChangeGateway(fromAddress, gateway string) *MsgChangeGateway {
	return &MsgChangeGateway{
		FromAddress: fromAddress,
		Gateway:     gateway,
	}
}

func (msg MsgChangeGateway) Route() string { return RouterKey }
func (msg MsgChangeGateway) Type() string  { return TypeMsgChangeGateway }
func (msg MsgChangeGateway) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgChangeGateway) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgChangeGateway) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}

	_, err = sdk.ValAddressFromBech32(msg.Gateway)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid gateway address")
	}

	return nil
}
func (m MsgChangeGateway) XXX_MessageName() string {
	return TypeMsgChangeGateway
}

// NewMsgMobileTransfer **************************/
func NewMsgMobileTransfer(fromAddress, toAddress, mobile string) *MsgMobileTransfer {
	return &MsgMobileTransfer{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Mobile:      mobile,
	}
}

func (msg MsgMobileTransfer) Route() string { return RouterKey }
func (msg MsgMobileTransfer) Type() string  { return TypeMsgMobileTransfer }
func (msg MsgMobileTransfer) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgMobileTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgMobileTransfer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}

	_, err = sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}

	return nil
}
func (m MsgMobileTransfer) XXX_MessageName() string {
	return TypeMsgMobileTransfer
}

// NewMsgAddressBookSave **************************/
func NewMsgAddressBookSave(fromAddress string, AddressBook []string) *MsgAddressBookSave {
	return &MsgAddressBookSave{
		FromAddress: fromAddress,
		AddressBook: AddressBook,
	}
}

func (msg MsgAddressBookSave) Route() string { return RouterKey }
func (msg MsgAddressBookSave) Type() string  { return TypeMsgAddressBookSave }
func (msg MsgAddressBookSave) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgAddressBookSave) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgAddressBookSave) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}

	return nil
}
func (m MsgAddressBookSave) XXX_MessageName() string {
	return TypeMsgAddressBookSave
}

// NewMsgSendGift **************************/
func NewMsgSendGift(fromAddress, nodeAddress, toAddress string, giftId, giftAmount int64, giftValue types.Coin) *MsgSendGift {
	return &MsgSendGift{
		FromAddress: fromAddress,
		NodeAddress: nodeAddress,
		ToAddress:   toAddress,
		GiftId:      giftId,
		GiftAmount:  giftAmount,
		GiftValue:   giftValue,
	}
}

func (msg MsgSendGift) Route() string { return RouterKey }
func (msg MsgSendGift) Type() string  { return TypeMsgSendGift }
func (msg MsgSendGift) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgSendGift) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgSendGift) ValidateBasic() error {

	if msg.GiftValue.Denom != config.BaseDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "coin error")
	}

	if !msg.GiftValue.Amount.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "amount error")
	}
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}

	_, err = sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid to address")
	}

	_, err = sdk.ValAddressFromBech32(msg.NodeAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid node address")
	}

	return nil
}
func (m MsgSendGift) XXX_MessageName() string {
	return TypeMsgSendGift
}

// NewMsgSetChatFee **************************/
func NewMsgSetChatFee(fromAddress string, fee types.Coin) *MsgSetChatFee {
	return &MsgSetChatFee{
		FromAddress: fromAddress,
		Fee:         fee,
	}
}

func (msg MsgSetChatFee) Route() string { return RouterKey }
func (msg MsgSetChatFee) Type() string  { return TypeMsgSetChatFee }
func (msg MsgSetChatFee) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgSetChatFee) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgSetChatFee) ValidateBasic() error {

	if msg.Fee.Denom != config.BaseDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "coin error")
	}

	if !msg.Fee.Amount.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "amount error")
	}
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}
func (m MsgSetChatFee) XXX_MessageName() string {
	return TypeMsgSetChatFee
}

// NewMsgRegister **************************/
func NewMsgRegister(fromAddress, nodeAddress, mobilePrefix string, mortgageAmount types.Coin) *MsgRegister {
	return &MsgRegister{
		FromAddress:    fromAddress,
		NodeAddress:    nodeAddress,
		MortgageAmount: mortgageAmount,
		MobilePrefix:   mobilePrefix,
	}
}
func (msg MsgRegister) Route() string { return RouterKey }
func (msg MsgRegister) Type() string  { return TypeMsgRegister }
func (msg MsgRegister) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgRegister) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgRegister) ValidateBasic() error {

	if msg.MortgageAmount.Denom != config.BaseDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "coin error")
	}

	if !msg.MortgageAmount.Amount.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "amount error")
	}
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	_, err = sdk.ValAddressFromBech32(msg.NodeAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}
func (m MsgRegister) XXX_MessageName() string {
	return TypeMsgRegister
}

//**************************/
