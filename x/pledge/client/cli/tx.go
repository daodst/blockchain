package cli

import (
	"fmt"
	"freemasonry.cc/blockchain/x/pledge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"
	"strings"
)

// NewTxCmd returns a root CLI command handler for erc20 transaction commands
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "pledge subcommands",
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

func NewPledgeDelegateProposalTxCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delegate [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a gateway delegate proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx gov submit-proposal delegate <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "delegate",
  "description": "gateway delegate",
  "delegate": [
    {
      "address": "",
      "gateway_address": "",
      "amount": "100"
    }
  ],
  "deposit": "1000fm"
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := ParsePledgeDelegateProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			content := types.NewPledgeDelegateProposal(proposal.Title, proposal.Description, proposal.Delegate.ToParamChanges())

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
}
