package cli

import (
	"fmt"
	"freemasonry.cc/blockchain/x/contract/types"
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
		Short:                      "contract subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand()
	return txCmd
}

func NewContractProposalTxCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "contract [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a contract info proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx gov submit-proposal delegate <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "delegate",
  "description": "gateway delegate",
  "contract": [
    {
      "contract_address": "",
      "description": "",
      "logo": "100",
      "website" : "",
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
			proposal, err := ParseContractProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			content := types.NewContractProposal(proposal.Title, proposal.Description, proposal.Contract.ToParamChanges())

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
