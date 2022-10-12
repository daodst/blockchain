package types

const (
	MSG_SMART_CREATE_VALIDATOR   = "comm/MsgCreateSmartValidator"
	MSG_GATEWAY_REGISTER         = "comm/MsgGatewayRegister"
	MSG_GATEWAY_INDEX_NUM        = "comm/MsgGatewayIndexNum"
	MSG_GATEWAY_UNDELEGATION     = "comm/MsgGatewayUndelegation"
	MSG_GATEWAY_BEGIN_REDELEGATE = "comm/MsgGatewayBeginRedelegate"
)

//
type Gateway struct {
	//
	GatewayAddress string `json:"gateway_address"`
	//
	GatewayName string `json:"gateway_name"`
	//
	GatewayUrl string `json:"gateway_url"`
	//
	GatewayQuota int64 `json:"gateway_quota"`
	// 0   1 
	Status int64 `json:"status"`
	//
	GatewayNum []GatewayNumIndex `json:"gateway_num"`
}

//
type GatewayNumIndex struct {
	//
	GatewayAddress string `json:"gateway_address"`
	//
	NumberIndex string `json:"number_index"`
	//
	NumberEnd []string `json:"number_end"`
	// 0:  1:  2:
	Status int64 `json:"status"`
	//()
	Validity int64 `json:"validity"`
}

//
type ValidatorInfor struct {
	ValidatorConsAddr string `json:"validator_consaddr"` //
	ValidatorStatus   string `json:"validator_status"`   //0  Unbonded 1 Unbonding 2 Bonded 3  4 
	ValidatorPubAddr  string `json:"validator_pubaddr"`  //
	ValidatorOperAddr string `json:"validator_operaddr"` // 
	AccAddr           string `json:"acc_addr"`           // 
}

type GatewayNumberCountReq struct {
	GatewayAddress string `json:"gateway_address"` //
	Amount         string `json:"amount"`
}

type IsValidReq struct {
	Number string `json:"number"` //
}
