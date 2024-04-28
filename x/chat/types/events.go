package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// chat events
const (

	
	EventFeeAddress = "fee_address"

	//register event 
	EventTypeRegister       = "chat_register"
	EventTypeChatAddress    = "chat_address"
	EventTypeRegAddress     = "reg_address"
	EventTypeGatewayAddress = "gateway_address"
	EventTypeFromBalance    = "from_balance"
	EventPrefixMobile       = "prefix_mobile"

	// set chat info
	SetChatInfoEventTypeFromAddress              = "set_chat_info_from_address"
	SetChatInfoEventTypeNodeAddress              = "set_chat_info_node_address"
	SetChatInfoEventTypeNodeFlag                 = "set_chat_info_node_flag"
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
	BurngetMobileEventTypeGetMobile     = "burn_get_mobile_get_mobile"

	//chat mobile transfer
	ChatMobileTransferEventTypeFromAddress = "chat_mobile_transfer_from_address"
	ChatMobileTransferEventTypeToAddress   = "chat_mobile_transfer_to_address"
	ChatMobileTransferEventTypeMobile      = "chat_mobile_transfer_mobile"
	ChatMobileTransferEventTypeFromBalance = "chat_mobile_transfer_from_balance"
)

// Event type for Transfer(address from, address to, uint256 value)
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
