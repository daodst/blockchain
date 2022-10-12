package types

const (

	//
	QueryTotalDelegate = "total_delegate"
)

type QueryUserInfoParams struct {
	Address string
}

type QueryDelegatorValidatorRequest struct {
	DelegatorAddr string `json:"delegator_addr"`
	ValidatorAddr string `json:"validator_addr"`
}
