package client

import (
	"freemasonry.cc/blockchain/x/contract/client/cli"
	"freemasonry.cc/blockchain/x/contract/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewContractProposalTxCmd, rest.ProposalRESTHandler)
