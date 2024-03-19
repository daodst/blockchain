package client

import (
	"freemasonry.cc/blockchain/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"testing"
)


func TestEvmTokenPair(t *testing.T) {
	evmClient := NewEvmClient()
	tokenPair, err := evmClient.GetTokenPair(core.UsdtDenom)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("tokenPair:", tokenPair.String())
}

// NFT
func TestEvmNft(t *testing.T) {
	evmClient := NewEvmClient()
	addr := "dex1va8aaeystat4twpy70ns7235pxwczg0698twv4"
	contract := "0x3921f8f5876ae2B077966C10a659d65fD8883d96"
	nft, err := evmClient.GetNftInfo(addr, contract)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(":", addr)
	t.Log("nft:", nft)
}

// NFT
func TestEvmQueryContractCode(t *testing.T) {
	evmClient := NewEvmClient()
	contractAddress, err := evmClient.GetContractCode("0x282EcBadB8F4E33E67797EB880FCB5deb2420305")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("code:", contractAddress)
}

// NFT
func TestEvmQueryContractCode1(t *testing.T) {
	s := []byte{0, 56, 186, 161, 81, 12, 72, 66, 253, 236, 59, 168, 122, 58, 205, 205, 225, 180, 113, 180, 255, 17, 56, 245, 225, 166, 255, 91, 240, 35, 84, 112}
	t.Log(hexutil.Encode(s))
	aa := "0x7944dfa6Ec9db969BB0b104D64dC6C406F6e81a3"
	sss := common.HexToAddress(aa).Hash()
	t.Log(sss.Bytes())
}

// NFT
func TestEvmQueryNftContractAddress(t *testing.T) {
	evmClient := NewEvmClient()
	contractAddress, err := evmClient.GetNftContractAddress()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("nft:", contractAddress)
}


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

// eth
func TestBlockNumber(t *testing.T) {
	evmClient := NewEvmClient()
	res, err := evmClient.BlockNumber()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(res.ToInt())
}


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


func TestEthEncode(t *testing.T) {
	addBytes := ethcrypto.Keccak256([]byte("add()"))
	t.Log(hexutil.Encode(addBytes))
	//t.Log(addBytes)
	subtractBytes := ethcrypto.Keccak256([]byte("subtract()"))
	t.Log(hexutil.Encode(subtractBytes))

}

func TestQueryCrossHash(t *testing.T) {
	cl := NewEvmClient()
	isExist, err := cl.QueryCrossChainHash("dex18mfr9qdq2m2yjdvywme5nkyct7lltmeflfd6ag", "123")
	t.Log(err)
	t.Log(isExist)

}
