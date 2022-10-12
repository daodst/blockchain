package cli

import (
	"freemasonry.cc/blockchain/x/pledge/types"
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
	PledgeDelegatesJSON []PledgeDelegateJSON

	PledgeDelegateJSON struct {
		Address        string `json:"address" yaml:"address"`
		GatewayAddress string `json:"gateway_address" yaml:"gateway_address"`
		Amount         string `json:"amount" yaml:"amount"`
	}

	DelegateProposalJSON struct {
		Title       string              `json:"title" yaml:"title"`
		Description string              `json:"description" yaml:"description"`
		Delegate    PledgeDelegatesJSON `json:"delegate" yaml:"delegate"`
		Deposit     string              `json:"deposit" yaml:"deposit"`
	}

	PledgeDelegateProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string              `json:"title" yaml:"title"`
		Description string              `json:"description" yaml:"description"`
		Delegate    PledgeDelegatesJSON `json:"delegate" yaml:"delegate"`
		Proposer    sdk.AccAddress      `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins           `json:"deposit" yaml:"deposit"`
	}
)

func (pcj PledgeDelegateJSON) ToPlegdeDelegate() types.PledgeDelegate {
	return types.NewPledgeDelegate(pcj.Address, pcj.GatewayAddress, pcj.Amount)
}

func (pcj PledgeDelegatesJSON) ToParamChanges() []types.PledgeDelegate {
	res := make([]types.PledgeDelegate, len(pcj))
	for i, pc := range pcj {
		res[i] = pc.ToPlegdeDelegate()
	}
	return res
}

func ParsePledgeDelegateProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (DelegateProposalJSON, error) {
	proposal := DelegateProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
