package client

import (
	"fmt"
	"testing"
)

/*
*  
 */
func TestBl1(t *testing.T) {
	client := NewBlockClient()
	block, err := client.Find(8151)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(":%d \n", block.Height)
	fmt.Printf(":%d \n", block.Time)
	fmt.Printf(":%s \n", block.Datahash)
	fmt.Printf(":%s \n", block.Apphash)
	fmt.Printf(":%s \n", block.LastCommitHash)
	fmt.Printf(":%s \n", block.ValidatorsHash)
	fmt.Printf(":%s \n", block.NextValidatorsHash)
	fmt.Printf(":%s \n", block.ProposerAddress)
	fmt.Printf("ID:%s \n", block.ChainId)
	fmt.Printf("LastResultsHash:%s \n", block.LastResultsHash)
	fmt.Printf("EvidenceHash:%s \n", block.EvidenceHash)
	fmt.Printf("LastBlockId:%s \n", block.LastBlockId)
	fmt.Printf("BlockId:%s \n", block.BlockId)

	//B258E27F2C5AF33A226A2F4F1E7DDDA20C4B4E4A513B6C6FCF34F52A49855D1B
	fmt.Println("------------------------")
	fmt.Printf(":%d \n", len(block.Txs))
	for i := 0; i < len(block.Txs); i++ {
		if i >= 0 {
			fmt.Println("------------------------")
		}
		fmt.Printf("TxTash %d:%s \n", i, block.Txs[i])
	}
}


func TestBlockResults(t *testing.T) {
	//D5E967FC8A604310371546FB0C9C1B32901D6BD2E4DC9FE1DA8CF4BE1A7F6730

	blockClient := NewBlockClient()
	var block int64 = 8150 //30069 
	events, err := blockClient.FindBlockResults(&block)
	if err != nil {
		t.Error(err)
		return
	}
	for _, k := range events {
		t.Log("type:", k.Type)
		for _, attr := range k.Attributes {
			t.Log("k:", string(attr.Key), "v:", string(attr.Value))
		}
	}
}
