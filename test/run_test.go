package test

import (
	"encoding/base64"
	"freemasonry.cc/blockchain/client"
	cmdcfg "freemasonry.cc/blockchain/cmd/config"
	"freemasonry.cc/blockchain/cmd/scd/cmd"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/comm/keeper"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
	"time"
)

var (
	txClient  = client.NewTxClient()
	accClient = client.NewAccountClient(&txClient)
)

func TestRun(t *testing.T) {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	cmdcfg.SetBip44CoinType(config)
	config.Seal()
	cmdcfg.RegisterDenoms()
	cmd.ChainRun("scd", "start", "--log_format", "json", "--log_level", "debug", "--home", "G:\\chatTest\\.scd")
}

func TestRun2(t *testing.T) {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	cmdcfg.SetBip44CoinType(config)
	config.Seal()
	cmdcfg.RegisterDenoms()
	ss, err := sdk.ParseCoinsNormalized("0.00001tt")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ss.String())
	decCoin, err := sdk.ParseDecCoin("0.1tt")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(decCoin)

}

func TestRun1(t *testing.T) {
	addr, _ := sdk.AccAddressFromBech32("dex1f6w59g7y33j7m4fg64nf2c9es6t9ktqez6mt8v")
	validatorInfoPubKeyBytes, err := base64.StdEncoding.DecodeString("SsEVVOCqw9e7+WyuWPtfZC3FMYt9vIg8YQdB3eCBaLw=")
	if err != nil {
		t.Error(err)
		return
	}
	pbk := ed25519.PubKey(validatorInfoPubKeyBytes) //ed25519
	pubkey, err := cryptocodec.FromTmPubKeyInterface(pbk)
	if err != nil {
		t.Error(err)
		return
	}
	des := stakingTypes.Description{"", "", "", "", ""}
	comm := stakingTypes.CommissionRates{sdk.NewDec(1), sdk.OneDec(), sdk.OneDec()}
	aa, err := stakingTypes.NewMsgCreateValidator(sdk.ValAddress(addr), pubkey, sdk.NewCoin("fm", sdk.NewInt(100)), des, comm, sdk.NewInt(1))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(aa.Pubkey.String())
	t.Log(7900 % 100)
}

func TestRegisterValidate(t *testing.T) {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	cmdcfg.SetBip44CoinType(config)
	config.Seal()
	cmdcfg.RegisterDenoms()
	//
	createAccAddr := "dex1f6w59g7y33j7m4fg64nf2c9es6t9ktqez6mt8v"
	createAcc, err := sdk.AccAddressFromBech32(createAccAddr)
	if err != nil {
		t.Fatal(err)
		return
	}
	//
	createAccAddrPrivateKey := "EF69C29AFA4C4418DDF3829264FDEDA504A2D93F8E51967D0C770F08349E3567"
	wallet, err := accClient.CreateAccountFromPriv(createAccAddrPrivateKey)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(":", wallet.Address, wallet.PublicKey)

	//
	validatorAddress := sdk.ValAddress(createAcc).String()
	t.Log("validatorAddress:", validatorAddress)

	validatorName := "test" //

	validatorPubkeyBase64 := "SsEVVOCqw9e7+WyuWPtfZC3FMYt9vIg8YQdB3eCBaLw="

	bech32ValidatorPubkey, err := keeper.ParseBech32ValConsPubkey(validatorPubkeyBase64)
	if err != nil {
		t.Fatal(err)
		return
	}

	RegisterValdator(t, wallet, bech32ValidatorPubkey, validatorName)
}

//
func RegisterValdator(t *testing.T, wallet *client.CosmosWallet, validatorPubkey cryptotypes.PubKey, nodeName string) {
	//
	createAccAddr := wallet.Address
	createAcc, err := sdk.AccAddressFromBech32(createAccAddr)
	if err != nil {
		t.Fatal(err)
		return
	}
	//
	createAccAddrPrivateKey := wallet.PrivateKey

	//
	validatorAddress := sdk.ValAddress(createAcc).String()
	t.Log("validatorAddress:", validatorAddress)

	//
	selfDelegation := core.RealString2LedgerCoin("100", "fm")

	//
	minSelfDelegation := sdk.NewInt(6000)

	name := nodeName //
	identity := "w"
	website := "w"
	contact := "w"
	remark := "w"

	description := stakingTypes.NewDescription(name, identity, website, contact, remark)

	rate, _ := sdk.NewDecFromStr("0.99")
	maxRate, _ := sdk.NewDecFromStr("0.99")
	maxChangeRate, _ := sdk.NewDecFromStr("0.0001")
	commission := stakingTypes.NewCommissionRates(rate, maxRate, maxChangeRate) //

	//bech32ValidatorPubkey, err := nodeClient.ParseBech32ValConsPubkey(validatorPubkeyBase64)
	//if err != nil {
	//	t.Fatal(err)
	//	return
	//}

	resp, err := txClient.RegisterValidator(
		createAccAddr, validatorAddress, validatorPubkey, selfDelegation,
		description, commission, minSelfDelegation, createAccAddrPrivateKey, 0.001)
	//
	if err == nil {
		txResp := resp.Data.(core.BroadcastTxResponse)
		t.Log("TxHash:", txResp.TxHash)
		for {
			result, _, err := txClient.FindByHex(txResp.TxHash)
			if err != nil {
				time.Sleep(time.Second)
				t.Log("。。。")
				continue
			}
			if result.TxResult.Code == 0 {
				t.Log(":", result.Height)
				t.Log("")
				for i1 := 0; i1 < len(result.TxResult.Events); i1++ {
					t.Log("TxResult Event ", i1, " type:", result.TxResult.Events[i1].Type)
					for i2 := 0; i2 < len(result.TxResult.Events[i1].Attributes); i2++ {
						t.Log("TxResult Event ", i1, " attribute ", i2, " key:", string(result.TxResult.Events[i1].Attributes[i2].Key), " value:", string(result.TxResult.Events[i1].Attributes[i2].Value))
					}
				}
				t.Log(result.TxResult.Events)
				break
			} else {
				t.Log("")
				t.Log(result.TxResult.Log)
				break
			}
		}

	} else {
		t.Log(":", err.Error())
	}
}
