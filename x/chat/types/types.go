package types

import (
	"github.com/cosmos/cosmos-sdk/types"
)

const (
	MsgTypeRegister        = "chat/MsgRegister"
	MsgTypeMortgage        = "chat/MsgMortgage"
	MsgTypeSetChatFee      = "chat/MsgTypeSetChatFee"
	MsgTypeSendGift        = "chat/MsgTypSendGift"
	MsgTypeAddressBookSave = "chat/MsgTypeAddressBookSave"
	MsgTypeGetRewards      = "chat/MsgTypeGetRewards"
	MsgTypeMobileTransfer  = "chat/MsgTypeMobileTransfer"
	MsgTypeChangeGateway   = "chat/MsgTypeChangeGateway"
	MsgTypeBurnGetMobile   = "chat/MsgTypeBurnGetMobile"
	MsgTypeSetChatInfo     = "chat/MsgTypeSetChatInfo"

	ChatRestrictedModeFee  = "fee"
	ChatRestrictedModeAny  = "any"
	ChatRestrictedModeList = "list"
)

//
type UserInfo struct {
	//
	FromAddress string `json:"from_address" yaml:"from_address"`
	//
	NodeAddress string `json:"node_address" yaml:"node_address"`
	//
	AddressBook []string `json:"address_book" yaml:"address_book"`
	//
	ChatBlacklist []string `json:"chat_blacklist" yaml:"chat_blacklist"`
	// （fee | any | list）
	ChatRestrictedMode string `json:"chat_restricted_mode" yaml:"chat_restricted_mode"`
	//
	ChatWhitelist []string `json:"chat_whitelist" yaml:"chat_whitelist"`
	//
	ChatFee types.Coin `json:"chat_fee" yaml:"chat_fee"`
	//
	Mobile []string `json:"mobile" yaml:"mobile"`
}

const (
	TransferTypeToModule  = "to_module"
	TransferTypeToAccount = "to_account"
)

//type RegisterData struct {
//	core.TxBase
//	core.TradeBase
//	FromAddress    string     `json:"from_address"`    //
//	NodeAddress    string     `json:"node_address"`    //
//	MortgageAmount types.Coin `json:"mortgage_amount"` //
//	MobilePrefix   string     `json:"mobile_prefix"`   //
//}

type MortgageInfo struct {
	MortgageRemain     types.Coin           `json:"mortgage_remain"`      //
	MortgageDevideInfo []MortgageDevideInfo `json:"mortgage_devide_info"` //
}

type MortgageDevideInfo struct {
	MortgageAddress string `json:"mortgage_address"` //
	MortgageAmount  string `json:"mortgage_amount"`  //
	ShowBalance     bool   `json:"show_balance"`     //
}

type LastReceiveLog struct {
	Height int64      `json:"height"`
	Value  types.Coin `json:"value"`
}

type MortgageAddLog struct {
	Height        int64      `json:"height"`
	MortgageValue types.Coin `json:"mortgage_value"`
}
