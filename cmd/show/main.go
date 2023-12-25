package main

import (
	"encoding/hex"
	"fmt"
	"freemasonry.cc/blockchain/client"
	"github.com/tendermint/tendermint/privval"
	"os"
	"path/filepath"
)


func main() {
	programPath, _ := filepath.Abs(os.Args[0])


	txClient := client.NewTxClient()
	accClient := client.NewAccountClient(&txClient)
	gatewayClient := client.NewGatewayClinet(&txClient)
	//.exe 
	executePath, _ := filepath.Split(programPath)
	chainRepo := filepath.Join(executePath,".stcd")
	pvKeyFile := filepath.Join(filepath.Join(chainRepo, "config"), "priv_validator_key.json")
	filePV := privval.LoadFilePVEmptyState(pvKeyFile, "")
	fmt.Println("validator priv key:", hex.EncodeToString(filePV.Key.PrivKey.Bytes()))
	validatorInfo, err := gatewayClient.ValidatorInfo()
	if err == nil {
		fmt.Println("gatewayAddr:", validatorInfo.ValidatorOperAddr)
	}
	acc,err := accClient.CreateAccountFromPriv(hex.EncodeToString(filePV.Key.PrivKey.Bytes()))
	if err == nil {
		fmt.Println("machineAddress:", acc.Address)
	}
}
