package cli

import (
	"github.com/spf13/cobra"

	"freemasonry.cc/blockchain/x/comm/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// NewTxCmd returns a root CLI command handler for erc20 transaction commands
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "erc20 subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
	//NewConvertCoinCmd(),
	//NewConvertERC20Cmd(),
	)
	return txCmd
}
