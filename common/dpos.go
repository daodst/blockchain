package common

func DposCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dpos",
		Short: "dpos create and unjail",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(dposCreateCmd())
	return cmd
}
