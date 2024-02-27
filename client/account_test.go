package client

import (
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"testing"
	"time"
)

/*
dst18ck9vqemdfh0na6nkvakke8qdqgxvk2gy3ex3t transfer txhash: D16DF81C7BACEA82AFE729572372570A3B72DFFA9860916FDFED6AA8455ADBC7  found: false  tx: <nil>
dst1zsqgp8jcvaljtkrs0w3uzrxltc3nafk5nz78el transfer txhash: 94C8164CB2464E3B5EB26E39D09F3D81E3FC320EDF025E873F12EB1C1E5A1357  found: false  tx: <nil>
dst16xncqn4h3tml4sjt57zwlhekryg7tx2z35xhjs transfer txhash: 933C09D1169F1C8932E7E5BF684ADEA471BA077D2C00F709FDDB0D0FA85692F5  found: false  tx: <nil>
dst1v0veq8v6vwdkgs9af74z8d4flhrawmw5rslvq4 gasPrice: 1001000.000000000000000000dst
dst1y0x0w7q2vqqsug4kjcmjz6fg5e6l65h76v20fg gasPrice: 1001000.000000000000000000dst
*/
func TestTransfer(t *testing.T) {
	txClient := NewTxClient()
	accClient := NewAccountClient(&txClient)

	sendWallet, err := accClient.CreateAccountFromSeed("imitate real clap airport husband east shove supply across point stage struggle twin wing gather rapid night inject future resource wink tide term height")
	if err != nil {
		t.Error(err)
		return
	}

	dstInt, _ := sdk.NewIntFromString("1000000000000000000000000")
	//fmInt, _ := sdk.NewIntFromString("100000000000000000000")
	data := core.TransferData{FromAddress: sendWallet.Address, ToAddress: "dst1d79gcr3f7ee4qvk3m8y8p9f0etghyslxdtjn2u", Coins: sdk.NewCoins(sdk.NewCoin(core.BaseDenom, dstInt))}
	data.Fee = core.NewLedgerFeeZero()
	_, resp, err := accClient.Transfer(data, sendWallet.PrivateKey)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("code:", resp.Code)
	t.Log("hash:", resp.TxHash)
}

func TestAccoun(t *testing.T) {
	dexAccount := "dex1tfnsctfjml9lskehnm5phvuxvpt45y20ythksh"
	ethAccount := "0xC0226D606FDb99fddd92AF3B4637f11ea3A5a36C"

	dex, _ := sdk.AccAddressFromBech32(dexAccount)
	eth := common.BytesToAddress(dex.Bytes())
	t.Logf("%s dex %s", dexAccount, eth.String())

	dexd := common.HexToAddress(ethAccount)
	dff := sdk.AccAddress(dexd[:20])
	t.Logf("%s dex %s", ethAccount, dff.String())

	feeAccountEth := "0xF0f4C5079BCf15a1f797326CE74aAC3375f5F693"
	feeAccountDex := "dex17xpfvakm2amg962yls6f84z3kell8c5l5s9l0c"

	dexd = common.HexToAddress(feeAccountEth)
	dff = sdk.AccAddress(dexd[:20])
	t.Logf("%s dex %s", feeAccountEth, dff.String())

	dex, _ = sdk.AccAddressFromBech32(feeAccountDex)
	eth = common.BytesToAddress(eth.Bytes())
	t.Logf("%s dex %s", feeAccountDex, eth.String())

	facc := "0x17c2bd128aaD7DD1f5b3dC31403528DcdF29863b"
	dexd = common.HexToAddress(facc)
	dff = sdk.AccAddress(dexd[:20])
	t.Logf("%s dex %s", facc, dff.String())

}


func TestCreateAccount(t *testing.T) {
	txClient := NewTxClient()
	accountClient := NewAccountClient(&txClient)

	sendWallet, err := accountClient.CreateAccountFromSeed("chimney meadow crop economy merit fitness receive penalty source crumble arena eager february sun end dinner link pulse thing observe then wreck eight toe")
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 5; i++ {
		seed, _ := accountClient.CreateSeedWord()
		toWallet, _ := accountClient.CreateAccountFromSeed(seed)
		t.Log(toWallet.Address, seed)

		dstInt, _ := sdk.NewIntFromString("100000000000000000")
		data := core.TransferData{FromAddress: sendWallet.Address, ToAddress: toWallet.Address, Coins: sdk.NewCoins(sdk.NewCoin(core.BaseDenom, dstInt))}
		data.Fee = core.NewLedgerFeeZero()
		_, resp, err := accountClient.Transfer(data, sendWallet.PrivateKey)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("transfer code:", resp.Code)
		<-time.After(time.Second * 6)
	}
}

func TestFindMainTokenBalances(t *testing.T) {
	txClient := NewTxClient()
	accountClient := NewAccountClient(&txClient)

	s := "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet"

	b, err := accountClient.FindMainTokenBalances(s)

	//t.Log(err)
	t.Log(b)
	t.Log(s)
	seed, _ := accountClient.CreateSeedWord()
	t.Log(seed)
	acc, err := accountClient.CreateAccountFromSeed("grab cage fine peace library gun waste industry need mention trim absent eager excite timber magic medal clock ritual remind flower divert hurt razor")
	if err != nil {
		t.Error(err)
		return
	}
	balances, err := accountClient.FindAccountBalances(acc.Address, "")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("balances:", balances)
	balances, err = accountClient.FindAccountBalances("dst1yzfgx4v3gw2d63prr206hq0w570xev4qglrteg", "")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("balances:", balances)
	t.Log(acc.Address)

}

func gasEstimate() sdk.Dec {
	//gasPrice := sdk.NewDec(0)
	//poolClient := NewMempoolClient()
	//first, _, err := poolClient.GetPoolFirstAndLastTxs()
	//if err != nil {
	//	return gasPrice
	//}
	
	////for _, tx := range txs {
	////	t.Log("gas:", tx.Fee.GetGas())
	////	t.Log("gasAmount:", tx.Fee.GetAmount())
	////	t.Log("gasPrices:", )
	////}
	//rate := sdk.NewDec(10000000)
	//addStep := rate.Mul(sdk.NewDec(1000))
	//if first != nil {
	//	for _, coin := range first.Fee.GasPrices() {
	//		if coin.Denom == "dst" {
	//			gasPrice = gasPrice.Add(coin.Amount.Mul(rate))
	//		}
	//	}
	//	gasPrice = gasPrice.Add(addStep) 
	//}
	//return gasPrice
	return sdk.Dec{}
}


func TestBatchTransfer(t *testing.T) {
	//txClient := NewTxClient()
	//accountClient := NewAccountClient(&txClient)
	//poolClient := NewMempoolClient()
	
	//go func() {
	//	fee := core.NewLedgerFee(0.00001) //gas
	//	fee.Gas = 10000000                //gas 
	//	for {
	//		<-time.After(time.Second * 2)
	//		sendWallet, err := accountClient.CreateAccountFromSeed("reunion road design filter sea symptom bird suggest region elder neutral paper")
	//		if err != nil {
	//			t.Error(err)
	//			return
	//		}
	//		t.Log("1 address:", sendWallet.Address, "gasPrice:", fee.GasPrices())
	//		dstInt, _ := sdk.NewIntFromString("1")
	//		data := core.TransferData{FromAddress: sendWallet.Address, ToAddress: "dst1xsst5cr8p4fpwsrx70t2jxuuekgplymwyy0f8m", Coins: sdk.NewCoins(sdk.NewCoin(core.BaseDenom, dstInt))}
	//		data.Fee = fee
	//		tx, resp, err := accountClient.Transfer(data, sendWallet.PrivateKey)
	//		if err != nil {
	//			t.Error("1", " txhash:", resp.TxHash, " Transfer error:", err)
	//			go func() {
	//				for {
	//					if txRes, err := poolClient.TxExisted(tx); err == nil {
	//						if txRes.Existed == 1 {
	//							t.Error("1", " txhash:", resp.TxHash, " existed in txpool")
	//						} else {
	//							t.Error("1", " txhash:", resp.TxHash, " not existed in txpool")
	//						}
	//					} else {
	//						break
	//					}
	//					<-time.After(time.Second * 2)
	//				}
	//			}()
	
	//			break
	//			//gas
	//			gasPrice := gasEstimate()
	//			if !gasPrice.IsZero() {
	//				fee.Amount[0].Amount = gasPrice.TruncateInt() //gas
	//			}
	//			continue
	//		}
	//		txRes2, notFound, err := txClient.FindByHex(resp.TxHash)
	//		if err != nil {
	//			t.Error("1 FindByHex error:", err)
	//			continue
	//		}
	//		t.Log("1", sendWallet.Address, "transfer txhash:", resp.TxHash, " found:", !notFound, " tx:", txRes2)
	
	//	}
	//}()
	
	//go func() {
	//	fee := core.NewLedgerFee(0.00002) //gas
	//	fee.Gas = 10000000                //gas 
	//	for {
	//		<-time.After(time.Second * 2)
	//		sendWallet, err := accountClient.CreateAccountFromSeed("pulse movie useful tumble modify false stem alone scissors canoe mad shield")
	//		if err != nil {
	//			t.Error(err)
	//			return
	//		}
	//		t.Log("2 address:", sendWallet.Address, "gasPrice:", fee.GasPrices())
	//		dstInt, _ := sdk.NewIntFromString("1")
	//		data := core.TransferData{FromAddress: sendWallet.Address, ToAddress: "dst1xsst5cr8p4fpwsrx70t2jxuuekgplymwyy0f8m", Coins: sdk.NewCoins(sdk.NewCoin(core.BaseDenom, dstInt))}
	//		data.Fee = fee
	//		tx, resp, err := accountClient.Transfer(data, sendWallet.PrivateKey)
	//		if err != nil {
	//			t.Error("2", " txhash:", resp.TxHash, " Transfer error:", err)
	//			go func() {
	//				for {
	//					if txRes, err := poolClient.TxExisted(tx); err == nil {
	//						if txRes.Existed == 1 {
	//							t.Error("2", " txhash:", resp.TxHash, " existed in txpool")
	//						} else {
	//							t.Error("2", " txhash:", resp.TxHash, " not existed in txpool")
	//						}
	//					} else {
	//						break
	//					}
	//					<-time.After(time.Second * 2)
	//				}
	//			}()
	//			break
	//			//gas
	//			gasPrice := gasEstimate()
	//			if !gasPrice.IsZero() {
	//				fee.Amount[0].Amount = gasPrice.TruncateInt() //gas
	//			}
	//			continue
	//		}
	//		txRes2, notFound, err := txClient.FindByHex(resp.TxHash)
	//		if err != nil {
	//			t.Error("2 FindByHex error:", err)
	//			continue
	//		}
	//		t.Log("2", sendWallet.Address, "transfer txhash:", resp.TxHash, " found:", !notFound, " tx:", txRes2)
	//	}
	//}()
	
	//go func() {
	//	fee := core.NewLedgerFee(0.00003) //gas
	//	fee.Gas = 10000000                //gas 
	//	for {
	//		<-time.After(time.Second * 2)
	//		sendWallet, err := accountClient.CreateAccountFromSeed("reveal wrong age sample issue expire give fish tennis coffee broccoli alone")
	//		if err != nil {
	//			t.Error(err)
	//			return
	//		}
	//		dstInt, _ := sdk.NewIntFromString("1")
	//		t.Log("3 address:", sendWallet.Address, "gasPrice:", fee.GasPrices())
	//		data := core.TransferData{FromAddress: sendWallet.Address, ToAddress: "dst1xsst5cr8p4fpwsrx70t2jxuuekgplymwyy0f8m", Coins: sdk.NewCoins(sdk.NewCoin(core.BaseDenom, dstInt))}
	//		data.Fee = fee
	//		_, resp, err := accountClient.Transfer(data, sendWallet.PrivateKey)
	//		if err != nil {
	//			t.Error("3 Transfer error:", err)
	
	//			//gas
	//			gasPrice := gasEstimate()
	//			if !gasPrice.IsZero() {
	//				fee.Amount[0].Amount = gasPrice.TruncateInt() //gas
	//			}
	//			continue
	//		}
	//		tx, notFound, err := txClient.FindByHex(resp.TxHash)
	//		if err != nil {
	//			t.Error("3 FindByHex error:", err)
	//			continue
	//		}
	//		t.Log("3", sendWallet.Address, "transfer txhash:", resp.TxHash, " found:", !notFound, " tx:", tx)
	//	}
	//}()
	
	//go func() {
	//	fee := core.NewLedgerFee(0.00004) //gas
	//	fee.Gas = 10000000                //gas 
	//	for {
	//		<-time.After(time.Second * 2)
	//		sendWallet, err := accountClient.CreateAccountFromSeed("umbrella bounce token weapon garden miracle crouch casino toddler source remove bitter")
	//		if err != nil {
	//			t.Error(err)
	//			return
	//		}
	//		t.Log("4 address:", sendWallet.Address, "gasPrice:", fee.GasPrices())
	//		dstInt, _ := sdk.NewIntFromString("1")
	//		data := core.TransferData{FromAddress: sendWallet.Address, ToAddress: "dst1xsst5cr8p4fpwsrx70t2jxuuekgplymwyy0f8m", Coins: sdk.NewCoins(sdk.NewCoin(core.BaseDenom, dstInt))}
	//		data.Fee = fee
	//		_, resp, err := accountClient.Transfer(data, sendWallet.PrivateKey)
	//		if err != nil {
	//			t.Error("4 Transfer error:", err)
	
	//			//gas
	//			gasPrice := gasEstimate()
	//			if !gasPrice.IsZero() {
	//				fee.Amount[0].Amount = gasPrice.TruncateInt() //gas
	//			}
	//			continue
	//		}
	//		tx, notFound, err := txClient.FindByHex(resp.TxHash)
	//		if err != nil {
	//			t.Error("4 FindByHex error:", err)
	//			continue
	//		}
	//		t.Log("4", sendWallet.Address, "transfer txhash:", resp.TxHash, " found:", !notFound, " tx:", tx)
	//	}
	//}()
	
	//go func() {
	//	fee := core.NewLedgerFee(0.00005) //gas
	//	fee.Gas = 10000000                //gas 
	//	for {
	//		<-time.After(time.Second * 2)
	//		sendWallet, err := accountClient.CreateAccountFromSeed("learn glad coyote picnic kiwi install manual cruel drink market stove empty")
	//		if err != nil {
	//			t.Error("error:", err)
	//			return
	//		}
	//		t.Log("5 address:", sendWallet.Address, "gasPrice:", fee.GasPrices())
	//		dstInt, _ := sdk.NewIntFromString("1")
	//		data := core.TransferData{FromAddress: sendWallet.Address, ToAddress: "dst1xsst5cr8p4fpwsrx70t2jxuuekgplymwyy0f8m", Coins: sdk.NewCoins(sdk.NewCoin(core.BaseDenom, dstInt))}
	//		data.Fee = fee
	//		_, resp, err := accountClient.Transfer(data, sendWallet.PrivateKey)
	//		if err != nil {
	//			t.Error("5 Transfer error:", err)
	
	//			//gas
	//			gasPrice := gasEstimate()
	//			if !gasPrice.IsZero() {
	//				fee.Amount[0].Amount = gasPrice.TruncateInt() //gas
	//			}
	//			continue
	//		}
	//		tx, notFound, err := txClient.FindByHex(resp.TxHash)
	//		if err != nil {
	//			t.Error("5 FindByHex error:", err)
	//			continue
	//		}
	//		t.Log("5", sendWallet.Address, "transfer txhash:", resp.TxHash, " found:", !notFound, " tx:", tx)
	//	}
	//}()
	//select {}
}

func TestMsgJson(t *testing.T) {
	t.Log()
}
