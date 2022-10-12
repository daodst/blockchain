package util

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

func PricheckSign(message string, signedAddress string) (bool, string) {
	var err error
	//var addr = strings.ToLower(pubaddress) //
	//addr = string([]byte(addr)[4:]) //
	message = hex.EncodeToString([]byte(strings.ToLower(message)))
	var msg = crypto.Keccak256([]byte(message))
	sign, err := hex.DecodeString(strings.ToLower(signedAddress))
	if err != nil {
		return false, ""
	}
	recoveredPub, err := crypto.Ecrecover(msg, sign)
	if err != nil {
		return false, ""
	}
	pubKey, err := crypto.UnmarshalPubkey(recoveredPub)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	recoveredAddrstr := hex.EncodeToString(recoveredAddr[:])
	return true, recoveredAddrstr
}
