package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

// constants
const (
	// module name
	ModuleName     = "chat"
	ModuleBurnName = "chat_burn"
	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName
)

// ModuleAddress is the native module address for EVM
var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

// prefix bytes for the chat persistent store
const (
	// MobileSuffixLength 
	MobileSuffixLength = 4 

	// MobileSuffixMax int
	MobileSuffixMax = 9999

	
	KeyPrefixRegisterInfo = "chat_register_info_"

	
	KeyPrefixChatSendGift = "chat_send_gift_"

	
	KeyPrefixMobileOwner = "chat_mobile_owner_"

	
	KeyPrefixGatewayIssueToken = "gateway_issue_token_"

	
	KeyChatAddress = "chat_address_"
)
