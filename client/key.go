package client

import (
	"encoding/hex"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmhd "github.com/tharsis/ethermint/crypto/hd"

	"github.com/tyler-smith/go-bip39"
)

func NewSecretKey() *SecretKey {
	return &SecretKey{}
}

//
type SecretKey struct {
}

//
func (k *SecretKey) CreateSeedWord() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}

	// generate (english) seed words based on the entropy
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

//()
func (k *SecretKey) CreateAccountFromSeed(mnemonic string) (*CosmosWallet, error) {
	keyringAlgos := keyring.SigningAlgoList{evmhd.EthSecp256k1}
	algo, err := keyring.NewSigningAlgoFromString("eth_secp256k1", keyringAlgos)
	if err != nil {
		return nil, err
	}
	hdPath := hd.CreateHDPath(CoinType, 0, 0).String()
	bip39Passphrase := ""
	derivedPriv, err := algo.Derive()(mnemonic, bip39Passphrase, hdPath)
	if err != nil {
		return nil, err
	}
	return k.CreateAccountFromPriv(hex.EncodeToString(derivedPriv))
}

//  hex.EncodeToString 
func (k *SecretKey) CreateAccountFromPriv(priv string) (*CosmosWallet, error) {
	privKeyBytes, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	keyringAlgos := keyring.SigningAlgoList{evmhd.EthSecp256k1}
	algo, err := keyring.NewSigningAlgoFromString("eth_secp256k1", keyringAlgos)
	if err != nil {
		return nil, err
	}
	privKey := algo.Generate()(privKeyBytes)
	//t.Log(privKey.Bytes())
	//t.Log("priv:",hex.EncodeToString(privKey.Bytes()))

	address := sdk.AccAddress(privKey.PubKey().Address())
	return &CosmosWallet{
		priv:       privKey,
		PrivateKey: priv,
		PublicKey:  hex.EncodeToString(privKey.PubKey().Bytes()),
		Address:    address.String()}, nil
}

func (k *SecretKey) Sign(addr *CosmosWallet, msg []byte) ([]byte, error) {
	return addr.priv.Sign(msg)
}

const CoinType = 60

//cosmos
type CosmosWallet struct {
	Address    string        `json:"address"`
	PublicKey  string        `json:"publickey"`
	PrivateKey string        `json:"privatekey"`
	priv       types.PrivKey `json:"priv"`
}

func (this *CosmosWallet) MarshalJson() []byte {
	data, _ := json.Marshal(this)
	return data
}
