package types

const (
	TypeContractProposal = "contract/ContractProposal"
)

type NftInfo struct {
	TokenId       int64  `json:"tokenId"`
	CreateAddress string `json:"createAddress"`
	Level         int64  `json:"level"`
	CreateTime    int64  `json:"createTime"`
}
