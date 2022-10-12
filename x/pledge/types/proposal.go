package types

import (
	"fmt"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeCommunityPoolSpend defines the type for a CommunityPoolSpendProposal
	ProposalTypePledgeDelegateProposal = "PledgeDelegateProposal"
)

// Assert CommunityPoolSpendProposal implements govtypes.Content at compile-time
var _ govtypes.Content = &PledgeDelegateProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypePledgeDelegateProposal)
	govtypes.RegisterProposalTypeCodec(&PledgeDelegateProposal{}, TypePledgeDelegateProposal)
}

// NewCommunityPoolSpendProposal creates a new community pool spned proposal.
//nolint:interfacer
func NewPledgeDelegateProposal(title, description string, delegate []PledgeDelegate) *PledgeDelegateProposal {
	return &PledgeDelegateProposal{title, description, delegate}
}

// GetTitle returns the title of a community pool spend proposal.
func (csp *PledgeDelegateProposal) GetTitle() string { return csp.Title }

// GetDescription returns the description of a community pool spend proposal.
func (csp *PledgeDelegateProposal) GetDescription() string { return csp.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (csp *PledgeDelegateProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (csp *PledgeDelegateProposal) ProposalType() string { return ProposalTypePledgeDelegateProposal }

// ValidateBasic runs basic stateless validity checks
func (csp *PledgeDelegateProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(csp)
	if err != nil {
		return err
	}
	if len(csp.Delegate) == 0 {
		return ErrEmptyProposalDelegate
	}

	return nil
}

func NewPledgeDelegate(address, gatewayAddress, amount string) PledgeDelegate {
	return PledgeDelegate{address, gatewayAddress, amount}
}

// String implements the Stringer interface.
func (csp PledgeDelegateProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Pledge Delegate Proposal:
  Title:       %s
  Description: %s
  Delegate:   %s
`, csp.Title, csp.Description, csp.Delegate))
	return b.String()
}
