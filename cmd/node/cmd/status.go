package cmd

import (
	"fmt"
	"freemasonry.cc/blockchain/util"
	"github.com/spf13/cobra"
	"time"
)

func StatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "lookup status info",
		Run: func(cmd *cobra.Command, args []string) {
			status, err := util.HttpGet("http://127.0.0.1:26657/status", time.Second*3)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(status)
			}
		},
	}
	return cmd
}
