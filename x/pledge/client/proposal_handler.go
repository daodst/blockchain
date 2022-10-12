package client

import (
	"freemasonry.cc/blockchain/x/pledge/client/cli"
	"freemasonry.cc/blockchain/x/pledge/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewPledgeDelegateProposalTxCmd, rest.ProposalRESTHandler)
