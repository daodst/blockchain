package cli

import (
	"freemasonry.cc/blockchain/x/contract/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"io/ioutil"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// ParseRegisterCoinProposal reads and parses a ParseRegisterCoinProposal from a file.
func ParseMetadata(cdc codec.JSONCodec, metadataFile string) (banktypes.Metadata, error) {
	metadata := banktypes.Metadata{}

	contents, err := ioutil.ReadFile(filepath.Clean(metadataFile))
	if err != nil {
		return metadata, err
	}

	if err = cdc.UnmarshalJSON(contents, &metadata); err != nil {
		return metadata, err
	}

	return metadata, nil
}

type (
	ContractsJSON []ContractJSON

	ContractJSON struct {
		ContractAddress string `json:"contract_address" yaml:"contract_address"`
		Description     string `json:"description" yaml:"description"`
		Logo            string `json:"logo" yaml:"logo"`
		Website         string `json:"website" yaml:"website"`
	}

	ContractProposalJSON struct {
		Title       string        `json:"title" yaml:"title"`
		Description string        `json:"description" yaml:"description"`
		Contract    ContractsJSON `json:"contract" yaml:"contract"`
		Deposit     string        `json:"deposit" yaml:"deposit"`
	}

	ContractProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Contract    ContractsJSON  `json:"contract" yaml:"contract"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
)

func (cj ContractJSON) ToContract() types.Contract {
	return types.NewContract(cj.ContractAddress, cj.Description, cj.Logo, cj.Website)
}

func (cj ContractsJSON) ToParamChanges() []types.Contract {
	res := make([]types.Contract, len(cj))
	for i, pc := range cj {
		res[i] = pc.ToContract()
	}
	return res
}

func ParseContractProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (ContractProposalJSON, error) {
	proposal := ContractProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
