package main

import (
	"fmt"
	cmdcfg "freemasonry.cc/blockchain/cmd/config"
	"freemasonry.cc/blockchain/cmd/node/cmd"
	scmd "freemasonry.cc/blockchain/cmd/scd/cmd"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/log"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)


func Start(cosmosRepoPath string) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
	var err error

	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	cmdcfg.SetBip44CoinType(config)
	config.Seal()
	cmdcfg.RegisterDenoms()

	log.Info("check config dir")
	//config,
	if e, _ := checkSourceExist(filepath.Join(cosmosRepoPath, "config")); !e {
		log.Info("init chain")

		
		err = chainRun(core.CommandName, "init", "node1", "--chain-id", core.ChainID, "--home", cosmosRepoPath)
		if err != nil {
			panic(err)
		}
		
		err = replaceConfig(cosmosRepoPath)
		if err != nil {
			panic(err)
		}
	}
	log.Info("check genesis.json")
	//genesis.json,,
	err = checkGenesisFile(cosmosRepoPath)
	if err != nil {
		panic(err)
	}

	log.Info("check config.toml")
	//config.toml,,
	err = checkConfigFile(cosmosRepoPath)
	if err != nil {
		panic(err)
	}

	//app.toml,,
	err = checkAppToml(cosmosRepoPath)
	if err != nil {
		panic(err)
	}

	//client.toml,,
	err = checkClientToml(cosmosRepoPath)
	if err != nil {
		panic(err)
	}

	//priv_validator_state.json,,
	err = checkValidatorStateJson(cosmosRepoPath)
	if err != nil {
		panic(err)
	}

	log.WithField("path", cosmosRepoPath).Info("chain repo")

	//logPath := filepath.Join(cosmosRepoPath, "chain.log") //cosmos

	//  debug | info | error
	logLevel := "error"
	logLevelSet, ok := os.LookupEnv("SMART_CHAIN_LOGGING") //chain
	if ok {
		logLevel = logLevelSet 
	}

	log.Info("start chain")
	os.Args = []string{core.CommandName, "start", "--log_format", "json", "--log_level", logLevel, "--home", cosmosRepoPath}
	rootCmd, _ := scmd.NewRootCmd()
	rootCmd.SetErr(os.Stdout)
	rootCmd.SetOut(os.Stdout)
	if err := svrcmd.Execute(rootCmd, cosmosRepoPath); err != nil {
		//os.Exit panic
		os.Exit(1)
	}
}

func chainRun(arg ...string) error {
	os.Args = arg
	rootCmd, _ := scmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, ""); err != nil {
		//os.Exit panic
		//os.Exit(1)
		return err
	}
	return nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start chain",
		Run: func(cmd *cobra.Command, args []string) {
			logStoraged, _ := cmd.Flags().GetBool("log")
			fmt.Println("log storaged:", logStoraged)

			
			if logStoraged {
				programPath, _ := filepath.Abs(os.Args[0])

				//.exe 
				runtimePath, _ := filepath.Split(programPath)

				logPath := filepath.Join(runtimePath, "log")

				if !PathExists(logPath) {
					err := os.Mkdir(logPath, 0644)
					if err != nil {
						panic(err)
					}
				}

				daemonLogPath := filepath.Join(logPath, "chain.log")

				log.EnableLogStorage(daemonLogPath, time.Hour*24*7, time.Hour*24) 
			}

			var tendermintRepoPath string
			home := cmd.Flag("home").Value.String()
			if home == "" {
				pwd, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				tendermintRepoPath = filepath.Join(pwd, ".scd")
			} else {
				tendermintRepoPath = home
			}

			fmt.Println("repo path:", tendermintRepoPath)

			Start(tendermintRepoPath)
		},
	}
	cmd.Flags().String("home", "", "chain repo path")
	cmd.Flags().Bool("log", true, "enabled log storage")
	return cmd
}

func main() {
	log.InitLogger(logrus.InfoLevel)

	rootCmd := &cobra.Command{
		Use:   core.CommandName,
		Short: "smart chain node",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	rootCmd.AddCommand(startCmd())
	rootCmd.AddCommand(cmd.VersionCmd())
	rootCmd.AddCommand(cmd.StatusCmd())
	rootCmd.Execute()
}
