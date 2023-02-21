package types

import (
	"fmt"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeCommunityPoolSpend defines the type for a CommunityPoolSpendProposal
	ProposalTypeContractProposal = "ContractProposal"
)

// Assert CommunityPoolSpendProposal implements govtypes.Content at compile-time
var _ govtypes.Content = &ContractProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeContractProposal)
	govtypes.RegisterProposalTypeCodec(&ContractProposal{}, ProposalTypeContractProposal)
}

// NewCommunityPoolSpendProposal creates a new community pool spned proposal.
//nolint:interfacer
func NewContractProposal(title, description string, contract []Contract) *ContractProposal {
	return &ContractProposal{title, description, contract}
}

// GetTitle returns the title of a community pool spend proposal.
func (csp *ContractProposal) GetTitle() string { return csp.Title }

// GetDescription returns the description of a community pool spend proposal.
func (csp *ContractProposal) GetDescription() string { return csp.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (csp *ContractProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (csp *ContractProposal) ProposalType() string { return ProposalTypeContractProposal }

// ValidateBasic runs basic stateless validity checks
func (csp *ContractProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(csp)
	if err != nil {
		return err
	}
	if len(csp.Contract) == 0 {
		return ErrEmptyProposalContract
	}

	return nil
}

func NewContract(address, description, logo, website string) Contract {
	return Contract{address, description, logo, website}
}

// String implements the Stringer interface.
func (csp ContractProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Contract Proposal:
  Title:       %s
  Description: %s
  Contract:   %s
`, csp.Title, csp.Description, csp.Contract))
	return b.String()
}
