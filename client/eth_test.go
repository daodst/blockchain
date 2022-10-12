package client

import (
	"context"
	"fmt"
	"freemasonry.cc/blockchain/client/contract"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
	"testing"
)

func TestLogFromToParam(t *testing.T) {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}

	ethAddress := common.HexToAddress("0x903b928D7213c539B5da7121C51F23EdA0012831")
	topicHash := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	var topics [][]common.Hash
	var topicsub []common.Hash
	topicsub = append(topicsub, topicHash)
	topics = append(topics, topicsub)

	contractAbi, err := abi.JSON(strings.NewReader(contract.SmartABI))
	if err != nil {
		panic(err)
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{ethAddress},
		Topics:    topics,
		FromBlock: big.NewInt(5250),
		ToBlock:   big.NewInt(5255),
	}
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		panic(err)
	}

	for _, vLog := range logs {
		var item contract.SmartTransfer

		contractAddress := common.HexToAddress("0x903b928D7213c539B5da7121C51F23EdA0012831")
		c := bind.NewBoundContract(contractAddress, contractAbi, nil, nil, nil)
		err := c.UnpackLog(&item, "Transfer", vLog)
		if err != nil {
			panic(err)
		}
		t.Log("", vLog.BlockNumber)
		t.Log("index", vLog.Index)
		t.Log("txhash", vLog.TxHash)
		t.Log("", item.From.String())
		t.Log("", item.To.String())
		t.Log("", item.Value.String())

		fmt.Println("********", vLog.BlockNumber, vLog.Index, vLog.TxHash.String(), item.From.String(), item.To.String(), item.Value.String(), item.Raw)
	}
}

func TestEth(t *testing.T) {
	/*client,err := geth.NewEthereumClient("http://127.0.0.1:8545")
	if err != nil{
		t.Error("",err)
	}
	context := geth.NewContext()
	filterQuery := geth.NewFilterQuery()
	filterQuery.SetFromBlock(geth.NewBigInt(5250))
	filterQuery.SetToBlock(geth.NewBigInt(5255))
	logs,err := client.FilterLogs(context,filterQuery)
	if err != nil{
		t.Error("",err)
	}
	for i := 0; i < logs.Size(); i++ {
		log, _ := logs.Get(i)
		t.Log("", log)
		topics := log.GetTopics()
		d0, _ := topics.Get(0)

		//d0:= common.HexToHash(log.GetTopics())
		d1, _ := topics.Get(1)
		d2, _ := topics.Get(2)
		t.Log("0", d0.String())
		t.Log("1", d1.String())
		t.Log("2", d2.String())
		if d0.String() == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" { //event hash

			var item contract.SmartTransfer
			contractAddress := common.HexToAddress("0xF0f4C5079BCf15a1f797326CE74aAC3375f5F693")
			// abi
			contractAbi, err := abi.JSON(strings.NewReader(contract.SmartABI))
			if err != nil {
				t.Error("abi", err)
			}
			c := bind.NewBoundContract(contractAddress, contractAbi, nil, nil, nil)
			//log.GetData()
			logedd := *log.(types.Log)
			err = c.UnpackLog(&item, "Transfer", logedd)
			//c.abi.UnpackIntoInterface(&item, event, log.Data);
			if err != nil {
				t.Error("", err)
			}
			t.Log("", item.Value.String())
		}
	}*/
}
