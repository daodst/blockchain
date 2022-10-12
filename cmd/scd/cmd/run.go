package cmd

import (
	"freemasonry.cc/blockchain/app"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"os"
)

func ChainRun(arg ...string) error {
	os.Args = arg
	rootCmd, _ := NewRootCmd()
	//rootCmd.SetOut(os.Stdout)
	//rootCmd.SetErr(os.Stderr)
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		//os.Exit panic
		//os.Exit(1)
		return err
	}
	return nil
}
