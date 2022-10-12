package client

import (
	"encoding/hex"
	cmdcfg "freemasonry.cc/blockchain/cmd/config"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmhd "github.com/tharsis/ethermint/crypto/hd"
	"testing"
)

func TestPrivImport(t *testing.T) {
	derivedPriv, _ := hex.DecodeString("dfcdf5b4376189b7df123b74e7f22e7cebe4cf0a41779501d8a6678877c7d28b")
	keyringAlgos := keyring.SigningAlgoList{evmhd.EthSecp256k1}
	algo, err := keyring.NewSigningAlgoFromString("eth_secp256k1", keyringAlgos)
	if err != nil {
		t.Error(err)
		return
	}
	privKey := algo.Generate()(derivedPriv)
	//t.Log(privKey.Bytes())
	t.Log("priv:", hex.EncodeToString(privKey.Bytes()))

	address := sdk.AccAddress(privKey.PubKey().Address())
	t.Log("address:", address)
	t.Log("address:", address.String())
}

func TestSeedImport(t *testing.T) {
	keyringAlgos := keyring.SigningAlgoList{evmhd.EthSecp256k1}
	algo, err := keyring.NewSigningAlgoFromString("eth_secp256k1", keyringAlgos)
	if err != nil {
		t.Error(err)
		return
	}
	hdPath := hd.CreateHDPath(CoinType, 0, 0).String()
	bip39Passphrase := ""
	mnemonic := "idle hamster sword stamp primary clay bright reduce collect tackle host slam they rookie few globe jealous tongue draw pencil useless blood flavor fluid"
	derivedPriv, err := algo.Derive()(mnemonic, bip39Passphrase, hdPath)
	if err != nil {
		t.Error(err)
		return
	}
	privKey := algo.Generate()(derivedPriv)
	//t.Log(privKey.Bytes())
	t.Log("priv:", hex.EncodeToString(privKey.Bytes()))

	address := sdk.AccAddress(privKey.PubKey().Address())
	t.Log("address:", address)
	t.Log("address:", address.String())
}

func TestAccCreate1(t *testing.T) {

	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	cmdcfg.SetBip44CoinType(config)

	sk := NewSecretKey()

	seed := "idle hamster sword stamp primary clay bright reduce collect tackle host slam they rookie few globe jealous tongue draw pencil useless blood flavor fluid"
	wallet, err := sk.CreateAccountFromSeed(seed)
	if err != nil {
		t.Error(err)
		return
	}
	/**
	  key_test.go:30: priv: dfcdf5b4376189b7df123b74e7f22e7cebe4cf0a41779501d8a6678877c7d28b
	  key_test.go:33: address: 8B0D9AE2BA9D3E8D7E113CB76DB234E6F2AAA680
	  key_test.go:34: address: dex13vxe4c46n5lg6ls38jmkmv35ume24f5qwa3mth
	*/
	t.Log(wallet.Address)
	t.Log(wallet.PublicKey)
	t.Log(wallet.PrivateKey)

	//uid := "zhangsan"
	//kb := keyring.NewInMemory(evmkr.Option())
	//keyringAlgos :=  keyring.SigningAlgoList{evmhd.EthSecp256k1}
	//algo, err := keyring.NewSigningAlgoFromString("eth_secp256k1", keyringAlgos)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//
	//hdPath := hd.CreateHDPath(CoinType, 0, 0).String()
	//info, err := kb.NewAccount(uid, seed, "", hdPath, algo)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//priv,_ := kb.ExportPrivKeyArmor(uid,"")
	//
	//privKey, _, err := crypto.UnarmorDecryptPrivKey(priv, "")
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//t.Log("priv:",hex.EncodeToString(privKey.Bytes()))
	//
	//t.Log("pbk:",info.GetPubKey().Address())
	//t.Log(info.GetAddress().String())
	//signText,pbk,err := kb.Sign(info.GetName(),[]byte("hello world"))
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//t.Log("sign:",signText)
	//t.Log("pbk:",pbk.Address().String())

}
