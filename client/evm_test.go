package client

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"testing"
)

//
func TestEvmBalance(t *testing.T) {
	evmClient := NewEvmClient()
	addr := "0xA825B785B7DA1C771EA2CFA51D475E73ABDF8212"
	balance, err := evmClient.GetBalance(addr)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(":", addr)
	t.Log(":", balance.ToInt())
}

//eth
func TestBlockNumber(t *testing.T) {
	evmClient := NewEvmClient()
	res, err := evmClient.BlockNumber()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(res.ToInt())
}

//
func TestGetBlockNumber(t *testing.T) {
	evmClient := NewEvmClient()
	res, err := evmClient.GetBlockNumber("1289", true)
	if err != nil {
		t.Fatal(err)
		return
	}
	num := hexutil.MustDecodeBig(res["number"].(string))
	timestamp := hexutil.MustDecodeBig(res["timestamp"].(string))
	t.Log(":", num)
	t.Log(":", timestamp)
	t.Log(":", res["miner"])
	for k, v := range res {
		t.Log(k, ":", v)
	}
	trans := res["transactions"].([]interface{})
	for txIndex, v1 := range trans {
		mapData := v1.(map[string]interface{})
		t.Log("------------------", txIndex)
		for k, v := range mapData {
			t.Log(k, ":", v)
		}
	}
}

//
func TestGetTransactionReceipt(t *testing.T) {
	evmClient := NewEvmClient()
	addr := "0xb9cc59d6fa8056f262a9e50f4273da898a0d4f5dc6a1300bf860c3efc3580e88"
	txData, err := evmClient.GetTransactionReceipt(addr)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("hash:", addr)
	//t.Log(":",txData)
	for k, v := range txData {
		t.Log(k, ":", v)
	}
}

//
func TestGetTransactionByHash(t *testing.T) {
	evmClient := NewEvmClient()
	addr := "0xb9cc59d6fa8056f262a9e50f4273da898a0d4f5dc6a1300bf860c3efc3580e88"
	txData, err := evmClient.GetTransactionByHash(addr)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("hash:", addr)
	//t.Log(":",txData)
	for k, v := range txData {
		t.Log(k, ":", v)
		if k == "gas" {
			dd, err := hexutil.DecodeBig(v.(string))
			if err != nil {
				t.Error("", err)
			}
			t.Log("gas  ", dd.String())
		}
		if k == "gasPrice" {
			dd, err := hexutil.DecodeBig(v.(string))
			if err != nil {
				t.Error("", err)
			}
			t.Log("gasPrice ", dd.String())
		}
	}
}

//
func TestEvmAddress(t *testing.T) {
	evmClient := NewEvmClient()
	address, err := evmClient.GetAddress()
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, addr := range address {
		t.Log(":", addr)
	}

}

func TestNetPeerCount(t *testing.T) {
	evmClient := NewEvmClient()
	count, err := evmClient.NetPeerCount()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(":", count)
}

//
func TestEthEncode(t *testing.T) {
	addBytes := ethcrypto.Keccak256([]byte("add()"))
	t.Log(hexutil.Encode(addBytes))
	//t.Log(addBytes)
	subtractBytes := ethcrypto.Keccak256([]byte("subtract()"))
	t.Log(hexutil.Encode(subtractBytes))

}
