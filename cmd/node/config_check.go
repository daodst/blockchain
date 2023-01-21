package main

import (
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/core/config"
	"io/ioutil"
	"os"
	"path/filepath"
)


func checkGenesisFile(path string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
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

//config.toml,
func checkConfigFile(path string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
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

//app.toml,
func checkAppToml(path string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
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

//client.toml,
func checkClientToml(path string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
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

//priv_validator_state.json,
func checkValidatorStateJson(path string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
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
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
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
