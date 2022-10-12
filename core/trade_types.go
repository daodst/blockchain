package core

import (
	"freemasonry.cc/trerr"
)

//

var (
	TradeTypeTransfer            = RegisterTranserType("transfer", "", "transfer accounts")
	TradeTypeDelegation          = RegisterTranserType("bonded", "POS", "POS mortgage")
	TradeTypeDelegationFee       = RegisterTranserType("bonded-fee", "POS", "POS mortgage service charge")
	TradeTypeUnbondDelegation    = RegisterTranserType("unbonded", "POS", "POS redemption")
	TradeTypeUnbondDelegationFee = RegisterTranserType("unbonded-fee", "POS", "POS redemption fee")
	TradeTypeFee                 = RegisterTranserType("fee", "", "Service charge expenditure")
	TradeTypeDelegationReward    = RegisterTranserType("delegate-reward", "POS", "POS mortgage reward")
	TradeTypeCommissionReward    = RegisterTranserType("commission-reward", "POS", "POS Commission reward")
	TradeTypeCommunityReward     = RegisterTranserType("community-reward", "", "Community reward")
	TradeTypeValidatorUnjail     = RegisterTranserType("unjail", "POS", "POS release")
	TradeTypeValidatorMinerBonus = RegisterTranserType("validator-miner-bonus", "POS", "POS incentive")
	TradeTypeValidatorCreate     = RegisterTranserType("validator-create", "", "Create validator")
	TradeTypeValidatorEditor     = RegisterTranserType("validator-editor", "", "Editor validator")
	TradeTypeCrossChainOut       = RegisterTranserType("cross-chain-out", "", "Cross Chain Out")
	TradeTypeCrossChainFee       = RegisterTranserType("cross-chain-fee", "", "Cross Chain Fee")
	TradeTypeCrossChainIn        = RegisterTranserType("cross-chain-in", "", "Cross Chain In")
	TradeTypeGatewayRegister     = RegisterTranserType("gateway-register", "", "Gateway register")
	TradeTypeChatRegister        = RegisterTranserType("event_devide_register", "", "Chat register share")
	TradeTypeChatMortgage        = RegisterTranserType("event_devide_mortgate", "", "Chat mortgage share")
	TradeTypeChatSendGift        = RegisterTranserType("event_devide_send_gift", "", "Chat send gift share")
	TradeTypeChatBonus           = RegisterTranserType("bonus", "fm ", "fm bonus")
	TradeTypeChatPledge          = RegisterTranserType("chat_delegate", "chat ", "chat pledge")
	TradeTypeChatUnpledge        = RegisterTranserType("chat_unpledge", "chat ", "chat unpledge")
	TradeTypeDevidePledge        = RegisterTranserType("devide_pledge", "", "devide pledge")
)

var tradeTypeText = make(map[string]string)
var tradeTypeTextEn = make(map[string]string)

//
func RegisterTranserType(key, value, enValue string) TranserType {
	tradeTypeTextEn[key] = enValue
	tradeTypeText[key] = value
	return TranserType(key)
}

//
func GetTranserTypeConfig() map[string]string {
	if trerr.Language == "EN" {
		return tradeTypeTextEn
	} else {
		return tradeTypeText
	}
}

type TranserType string

func (this TranserType) GetValue() string {
	if text, ok := tradeTypeText[string(this[:])]; ok {
		return text
	} else {
		return ""
	}
}

func (this TranserType) GetKey() string {
	return string(this[:])
}
