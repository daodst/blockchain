package cmd

import (
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/core/config"
	"freemasonry.cc/blockchain/x/dao/keeper"
	"freemasonry.cc/log"
	"github.com/evmos/ethermint/encoding"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func daemonCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "secret telegram chain start",
		Run: func(cmd *cobra.Command, args []string) {
			logStoraged, _ := cmd.Flags().GetBool("log")
			keeper.EncodingConfig = encoding.MakeConfig(app.ModuleBasics)
			
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

			tendermintRepoPath := cmd.Flag("home").Value.String()

			start(tendermintRepoPath, cmd, args)
		},
	}
	cmd.Flags().Bool("log", true, "enabled log storage")
	return cmd
}


func start(cosmosRepoPath string, cmd *cobra.Command, args []string) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	logStoraged, _ := cmd.Flags().GetBool("log")

	log.Info("log storaged:", logStoraged)
	log.Info("repo path:", cosmosRepoPath)

	var err error

	log.Info("check config dir", len(cmd.Commands()))
	//config,
	if e, _ := checkSourceExist(filepath.Join(cosmosRepoPath, "config/genesis.json")); !e {
		log.Info("init chain")

		cmd.Root().SetArgs([]string{
			"init", "node1", "--chain-id", core.ChainID, "--home", cosmosRepoPath,
		})
		err = cmd.Root().Execute()
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

	//  debug | info | error
	logLevelSet, ok := os.LookupEnv("DST_CHAIN_LOGGING") //chain
	if ok {
		
		cmd.Root().SetArgs([]string{
			"start", "--log_format", "json", "--log_level", logLevelSet,
		})
	} else {
		cmd.Root().SetArgs([]string{
			"start", "--log_format", "json",
		})
	}

	log.Info("start chain")

	err = cmd.Root().Execute()
	if err != nil {
		log.WithError(err).Error("start.RunE")
		panic(err)
	}
	log.Info("exit")
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}


func checkGenesisFile(path string) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	genesisPath := filepath.Join(path, "config/genesis.json")

	
	if exist, _ := checkSourceExist(genesisPath); exist {
		return nil
	}
	err := ioutil.WriteFile(genesisPath, []byte(config.GenesisJson), os.ModePerm)
	if err != nil {
		log.WithError(err).WithField("file", genesisPath).Error("ioutil.WriteFile")
	}
	return err
}

// config.toml,
func checkConfigFile(path string) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	configPath := filepath.Join(path, "config/config.toml")

	
	if exist, _ := checkSourceExist(configPath); exist {
		return nil
	}
	err := ioutil.WriteFile(configPath, []byte(config.ConfigToml), os.ModePerm)
	if err != nil {
		log.WithError(err).WithField("file", configPath).Error("ioutil.WriteFile")
	}
	return err
}

// app.toml,
func checkAppToml(path string) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	filePath := filepath.Join(path, "config/app.toml")

	
	if exist, _ := checkSourceExist(filePath); exist {
		return nil
	}
	err := ioutil.WriteFile(filePath, []byte(config.AppToml), os.ModePerm)
	if err != nil {
		log.WithError(err).WithField("file", filePath).Error("ioutil.WriteFile")
	}
	return err
}

// client.toml,
func checkClientToml(path string) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	filePath := filepath.Join(path, "config/client.toml")

	
	if exist, _ := checkSourceExist(filePath); exist {
		return nil
	}
	err := ioutil.WriteFile(filePath, []byte(config.ClientToml), os.ModePerm)
	if err != nil {
		log.WithError(err).WithField("file", filePath).Error("ioutil.WriteFile")
	}
	return err
}

// priv_validator_state.json,
func checkValidatorStateJson(path string) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	filePath := filepath.Join(path, "config/priv_validator_state.json")

	
	if exist, _ := checkSourceExist(filePath); exist {
		return nil
	}
	err := ioutil.WriteFile(filePath, []byte(config.ValidatorStateJson), os.ModePerm)
	if err != nil {
		log.WithError(err).WithField("file", filePath).Error("ioutil.WriteFile")
	}
	return err
}

func replaceConfig(path string) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	appTomlPath := filepath.Join(path, "config/app.toml")
	if exist, err := checkSourceExist(appTomlPath); !exist {
		return err
	}
	genesisPath := filepath.Join(path, "config/genesis.json")
	if exist, err := checkSourceExist(genesisPath); !exist {
		return err
	}
	configTomlPath := filepath.Join(path, "config/config.toml")
	if exist, err := checkSourceExist(configTomlPath); !exist {
		return err
	}
	clientTomlPath := filepath.Join(path, "config/client.toml")
	if exist, err := checkSourceExist(clientTomlPath); !exist {
		return err
	}
	appTomlBuf, err := ioutil.ReadFile(appTomlPath)
	if err != nil {
		log.WithError(err).WithField("path", appTomlPath).Error("ioutil.ReadFile")
		return err
	}
	genesisBuf, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		log.WithError(err).WithField("path", genesisPath).Error("ioutil.ReadFile")
		return err
	}
	configTomlBuf, err := ioutil.ReadFile(configTomlPath)
	if err != nil {
		log.WithError(err).WithField("path", configTomlPath).Error("ioutil.ReadFile")
		return err
	}
	clientTomlBuf, err := ioutil.ReadFile(clientTomlPath)
	if err != nil {
		log.WithError(err).WithField("path", clientTomlPath).Error("ioutil.ReadFile")
		return err
	}

	appTomlContent := string(appTomlBuf)
	genesisContent := string(genesisBuf)
	configTomlContent := string(configTomlBuf)
	clientTomlContent := string(clientTomlBuf)

	
	if appTomlContent != config.AppToml {
		err = ioutil.WriteFile(appTomlPath, []byte(config.AppToml), os.ModePerm)
		if err != nil {
			log.WithError(err).WithField("path", appTomlPath).Error("ioutil.WriteFile")
			return err
		}
	}
	if genesisContent != config.GenesisJson {
		err = ioutil.WriteFile(genesisPath, []byte(config.GenesisJson), os.ModePerm)
		if err != nil {
			log.WithError(err).WithField("path", genesisPath).Error("ioutil.WriteFile")
			return err
		}
	}

	if configTomlContent != config.ConfigToml {
		err = ioutil.WriteFile(configTomlPath, []byte(config.ConfigToml), os.ModePerm)
		if err != nil {
			log.WithError(err).WithField("path", configTomlPath).Error("ioutil.WriteFile")
			return err
		}
	}

	if clientTomlContent != config.ClientToml {
		err = ioutil.WriteFile(clientTomlPath, []byte(config.ClientToml), os.ModePerm)
		if err != nil {
			log.WithError(err).WithField("path", clientTomlPath).Error("ioutil.WriteFile")
			return err
		}
	}
	return err
}


func checkSourceExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
