package types

const (
	
	QueryParams = "params"
	//NFT
	QueryNft = "nft"
	//NFT
	QueryNftContractAddress = "contract_address"
	//code
	QueryContractCode = "contract_code"
)

//NFT
type QueryNftInfoParams struct {
	Address         string `json:"address"`          
	ContractAddress string `json:"contract_address"` 
}
