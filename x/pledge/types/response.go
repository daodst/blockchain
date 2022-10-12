package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BaseResponse struct {
	Info   string `json:"info"`
	Status int    `json:"status"`
}

type TxResponse struct {
	Height string `json:"height"`
	TxHash string `json:"txhash"`
	Info   string `json:"info"`
}

//Tx
type BroadcastTxResponse struct {
	BaseResponse
	Height      int64  `json:"height"`
	TxHash      string `json:"txhash"`
	Codespace   string `json:"codespace"`     //
	Code        uint32 `json:"code"`          //
	SignedTxStr string `json:"signed_tx_str"` //
}

func (this *BaseResponse) IsSuccess() bool {
	return this.Status == 0
}

type AccountNumberSeqResponse struct {
	BaseResponse
	AccountNumber uint64 `json:"account_number"` //，-1
	Sequence      uint64 `json:"sequence"`
	NotFound      bool   `json:"not_found"` //
}

//
type BalanceResponse struct {
	BaseResponse
	Height string  `json:"height"`
	Token  []Token `json:"result"`
}

type Token struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

//
type DelegationPreviewResponse struct {
	Shares        string `json:"shares"`
	SourceAmount  string `json:"source_amount"`
	SourceShares  string `json:"source_shares"`
	BalanceAmount string `json:"balance_amount"`
	BalanceShares string `json:"balance_shares"`
}

//
type UnbondingDelegationPreviewResponse struct {
	Amount        string `json:"amount"`         //
	Shares        string `json:"shares"`         //
	BalanceAmount string `json:"balance_amount"` //
	BalanceShares string `json:"balance_shares"` //
	SourceAmount  string `json:"source_amount"`  //
	SourceShares  string `json:"source_shares"`  //
}

//
type DelegationDetailResponse struct {
	DelegationAddr string `json:"delegation_addr"`
	ValidatorAddr  string `json:"validator_addr"`
	Amount         string `json:"amount"` //
	Shares         string `json:"shares"` //
}

//pos
type PosReportFormResponse struct {
	UnbondAmount        string `json:"unbond_amount"`         //
	PosRewardReceived   string `json:"pos_reward_received"`   //pos
	PosRewardUnreceived string `json:"pos_reward_unreceived"` //pos
	MortgAmount         string `json:"mortg_amount"`          //
	Shares              string `json:"shares"`                //
	ValidatorShares     string `json:"validator_shares"`      //
	TotalShares         string `json:"total_shares"`          //
	Account             string `json:"account"`               //
	AccountTotalShares  string `json:"account_total_shares"`  //
}

//
type ValidatorRegisterLimit struct {
	MortgAmount string `json:"mortg_amount"` //
	Status      string `json:"status"`
}

//
type ValidatorsDelegationResp struct {
	Shares  string `json:"shares"`
	Balance string `json:"balance"`
}

//
type ValidatorCommission struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	Reward           sdk.DecCoin    `json:"reward"`
}
type ValidatorCommissionResp struct {
	ValidatorCommissions []ValidatorCommission `json:"validator_commissions"`
	Total                sdk.DecCoin           `json:"total"`
}
