package client

import (
	"fmt"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/sirupsen/logrus"
)

type AccountClient struct {
	TxClient  *TxClient
	key       *SecretKey
	ServerUrl string
	logPrefix string
}

type Account struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Address  string `json:"address"`
	Pubkey   string `json:"pubkey"`
	Mnemonic string `'json:"mnemonic"`
}

func (this *Account) Print() {
	fmt.Printf("Name:\t %s \n", this.Name)
	fmt.Printf("Address:\t %s \n", this.Address)
	fmt.Printf("Type:\t\t %s \n", this.Type)
	fmt.Printf("Pubkey:\t\t %s \n", this.Pubkey)
	fmt.Printf("Menmonic:\t\t %s \n", this.Mnemonic)
}

type AccountList struct {
	Accounts []Account
}

//
func (this *AccountClient) GetAllAccounts() (accounts []string, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	reponseStr, _, err := clientCtx.QueryWithData("custom/auth/accounts", []byte{})
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(reponseStr, &accounts)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON2")
		return
	}
	return
}

//token
func (this *AccountClient) FindAccountBalances(accountAddr string, height string) (coins sdk.Coins, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"acc": accountAddr, "height": height})
	url := "/bank/balances/" + accountAddr
	if height != "" {
		url += "?height=" + height
	}
	reponseStr, err := GetRequest(this.ServerUrl, url)
	if err != nil {
		log.Error("GetRequest")
		return
	}
	var resp = rest.ResponseWithHeight{}
	err = clientCtx.LegacyAmino.UnmarshalJSON([]byte(reponseStr), &resp)
	if err != nil {
		log.Error("UnmarshalJSON1")
		return
	}
	//var ledgerCoins sdk.Coins
	err = clientCtx.LegacyAmino.UnmarshalJSON(resp.Result, &coins)
	if err != nil {
		log.Error("UnmarshalJSON2")
		return
	}

	return
}

//token
func (this *AccountClient) FindAccountBalance(accountAddr string, denom, height string) (coin sdk.Coin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"acc": accountAddr, "denom": denom, "height": height})
	url := "/bank/balances/" + accountAddr + "?denom=" + denom
	if height != "" {
		url += "&height=" + height
	}
	reponseStr, err := GetRequest(this.ServerUrl, url)
	if err != nil {
		log.Error("GetRequest")
		return
	}

	var resp = rest.ResponseWithHeight{}
	err = clientCtx.LegacyAmino.UnmarshalJSON([]byte(reponseStr), &resp)
	if err != nil {
		log.Error("UnmarshalJSON1")
		return
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resp.Result, &coin)
	if err != nil {
		log.Error("UnmarshalJSON2")
		return
	}
	//realCoins = core.MustLedgerCoin2RealCoin(coin)
	return
}

func (this *AccountClient) FindBalanceByRpc(accountAddr string, denom string) (realCoins core.RealCoin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"acc": accountAddr, "denom": denom})

	req := banktypes.QueryBalanceRequest{Address: accountAddr, Denom: denom}

	reqBytes, _ := clientCtx.LegacyAmino.MarshalJSON(req)

	reponseStr, _, err := clientCtx.QueryWithData("custom/bank/balance", reqBytes)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	var coin sdk.Coin
	err = clientCtx.LegacyAmino.UnmarshalJSON(reponseStr, &coin)
	if err != nil {
		log.Error("UnmarshalJSON2")
		return
	}
	realCoins = core.MustLedgerCoin2RealCoin(coin)
	return
}

//  hex.EncodeToString 
func (this *AccountClient) CreateAccountFromPriv(priv string) (*CosmosWallet, error) {
	return this.key.CreateAccountFromPriv(priv)
}

//()
func (this *AccountClient) CreateAccountFromSeed(seed string) (acc *CosmosWallet, err error) {
	return this.key.CreateAccountFromSeed(seed)
}

//
func (this *AccountClient) CreateSeedWord() (mnemonic string, err error) {
	return this.key.CreateSeedWord()
}
