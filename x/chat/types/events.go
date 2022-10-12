package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// chat events
const (

	//
	EventFeeAddress = "fee_address"

	//register event 

	EventTypeFromAddress              = "from_address"
	EventTypeNodeAddress              = "node_address"
	EventTypeMortGageAmount           = "mortgage_amount"
	EventTypeMortGageDenom            = "mortgage_denom"
	EventTypeMortgageRemain           = "mortgage_remain"
	EventTypeGetMobile                = "get_mobile"
	EventTypePledgeFee                = "pledge_fee"
	EventTypePledgeAmount             = "pledge_amount"
	EventTypeFromBalance              = "from_balance"
	EventPrefixMobile                 = "prefix_mobile"
	EventDefaultChatRestrictedMode    = "default_chat_restricted_mode"
	EventDefaultChatFee               = "default_chat_fee"
	SetChatInfoEventTypeChatFeeAmount = "chat_fee_amount"
	SetChatInfoEventTypeChatFeeDenom  = "chat_fee_denom"

	//Event_Devide 
	EventTypeDevide = "event_devide"
	//
	DevideEventFromAddress = "from_address"
	//
	DevideEventToAddress = "to_address"
	//
	DevideEventAmount = "devide_amount"
	//
	DevideEventDenom = "devide_denom"
	//
	DevideEventBalance = "balance"

	//
	EventTypeDevideType     = "event_devide_type"
	EventTypeDevideRegister = "event_devide_register"  //
	EventTypeDevidePledge   = "event_devide_pledge"    //
	EventTypeDevideSendGift = "event_devide_send_gift" //
	EventTypeDevideRedeem   = "event_devide_redeem"    //

	//mortgage event 

	MortgageEventTypeType           = "mortgage_type"
	MortgageEventTypeFromAddress    = "mortgage_from_address"
	MortgageEventTypeDenom          = "mortgage_denom"
	MortgageEventTypeMortgageAmount = "mortgage_mortgage_amount"
	MortgageEventTypeFromBalance    = "mortgage_from_balance"

	MortgageEventTypePledgeFeeAmount = "mortgage_pledge_fee_amount"
	//
	MortgageEventTypeModuleAddress = "mortgage_module_address"
	// send gift

	SendGiftEventTypeFromAddress  = "send_gift_from_address"
	SendGiftEventTypeToAddress    = "send_gift_to_address"
	SendGiftEventTypeGateAddress  = "send_gift_gate_address"
	SendGiftEventTypeGiftId       = "send_gift_gift_id"
	SendGiftEventTypeGiftValue    = "send_gift_gift_value"
	SendGiftEventTypeGiftDenom    = "send_gift_gift_denom"
	SendGiftEventTypeGiftAmount   = "send_gift_gift_amount"
	SendGiftEventTypeGiftValueAll = "send_gift_gift_value_all"
	SendGiftEventTypeGiftReceive  = "send_gift_gift_receive"

	// get rewards

	GetRewardEventTypeFromAddress       = "get_rewards_from_address"
	GetRewardEventTypeMortgageAmountAdd = "get_rewards_mortgage_amount_add"
	GetRewardEventTypeMortgageAmountNew = "get_rewards_mortgage_amount_new"
	GetRewardEventTypeDenom             = "get_rewards_denom"

	// change gateway

	ChangeGatewayEventTypeFromAddress = "change_gateway_from_address"
	ChangeGatewayEventTypeOldGateWay  = "change_gateway_old_gateway"
	ChangeGatewayEventTypeNewGateWay  = "change_gateway_new_gateway"

	// set chat fee

	SetChatFeeEventTypeFromAddress = "set_chat_fee_from_address"
	SetChatFeeEventTypeFee         = "set_chat_fee_fee"
	SetChatFeeEventTypeDenom       = "set_chat_fee_denom"

	// set chat info
	SetChatInfoEventTypeFromAddress              = "set_chat_info_from_address"
	SetChatInfoEventTypeNodeAddress              = "set_chat_info_node_address"
	SetChatInfoEventTypeChatRestrictedMode       = "set_chat_info_chat_restricted_mode"
	SetChatInfoEventTypeChatFee                  = "set_chat_info_chat_fee"
	SetChatInfoEventTypeChatBlacklist            = "set_chat_info_chat_blacklist"
	SetChatInfoEventTypeChatWhitelist            = "set_chat_info_chat_whitelist"
	SetChatInfoEventTypeAddressBook              = "set_chat_info_address_book"
	SetChatInfoEventTypeGatewayEventPrefixMobile = "set_chat_info_gateway_event_prefix_mobile"
	SetChatInfoEventTypeGatewayEventFromBalance  = "set_chat_info_gateway_event_from_balance"

	// burn get mobile
	BurngetMobileEventTypeFromAddress   = "burn_get_mobile_from_address"
	BurngetMobileEventTypeDenom         = "burn_get_mobile_denom"
	BurngetMobileEventTypeAmount        = "burn_get_mobile_amount"
	BurngetMobileEventTypeFromBalance   = "burn_get_mobile_from_balance"
	BurngetMobileEventTypeModuleAddress = "burn_get_mobile_module_address"
)

// Event type for Transfer(address from, address to, uint256 value)
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
