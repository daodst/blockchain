package client

import (
	"encoding/json"
	"fmt"
	"freemasonry.cc/blockchain/client/contract"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogo/protobuf/proto"
	"github.com/tharsis/ethermint/x/evm/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	"strings"
	"testing"
)

//tx
func TestTxFind(t *testing.T) {
	txClient := NewTxClient()
	//tx
	txHash := "AC0AA54F78FDFF82327F24C55D6DBC975820E64F960CAD8B0E906873CE34D47B"
	resTx, notFound, err := txClient.FindByHex(txHash)
	if err != nil {
		t.Error("error:", err)
		return
	}
	if notFound {
		t.Log("")
		return
	}
	fmt.Println("Tx:", resTx.Tx.String())
	fmt.Println("Hash:", resTx.Hash)
	fmt.Println("Height:", resTx.Height)
	fmt.Println("Index:", resTx.Index)
	fmt.Println("tx_hash", resTx.Tx.Hash())

	fmt.Println("Proof:", resTx.Proof.Proof.String())
	fmt.Println("Proof Data Hash:", resTx.Proof.Data.Hash())
	fmt.Println("Proof Data:", resTx.Proof.Data)
	fmt.Println("Proof RootHash:", resTx.Proof.RootHash)
	fmt.Println("TxResult Event")
	for i1 := 0; i1 < len(resTx.TxResult.Events); i1++ {
		fmt.Println("TxResult Event ", i1, " type:", resTx.TxResult.Events[i1].Type)
		for i2 := 0; i2 < len(resTx.TxResult.Events[i1].Attributes); i2++ {
			fmt.Println("TxResult Event ", i1, " attribute ", i2, " key:", string(resTx.TxResult.Events[i1].Attributes[i2].Key), " value:", string(resTx.TxResult.Events[i1].Attributes[i2].Value))

			if resTx.TxResult.Events[i1].Type == "tx_log" {
				var log evmtypes.Log
				err := json.Unmarshal(resTx.TxResult.Events[i1].Attributes[i2].Value, &log)
				if err != nil {
					t.Error("log", err)
					continue
				}
				d0 := common.HexToHash(log.Topics[0])
				d1 := common.HexToAddress(log.Topics[1])
				d2 := common.HexToAddress(log.Topics[2])
				t.Log("0", d0.String())
				t.Log("1", d1.String())
				t.Log("2", d2.String())
				if log.Topics[0] == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" { //event hash

					var item contract.SmartTransfer
					contractAddress := common.HexToAddress("0xF0f4C5079BCf15a1f797326CE74aAC3375f5F693")
					// abi
					contractAbi, err := abi.JSON(strings.NewReader(contract.SmartABI))
					if err != nil {
						t.Error("abi", err)
					}
					c := bind.NewBoundContract(contractAddress, contractAbi, nil, nil, nil)
					err = c.UnpackLog(&item, "Transfer", *log.ToEthereum())
					//c.abi.UnpackIntoInterface(&item, event, log.Data);
					if err != nil {
						t.Error("", err)
					}
					t.Log("", item.Value.String())
				}
			}
		}

		/*if "tx_log" == resTx.TxResult.Events[i1].Type{
			erc20 := contracts.ERC20UsdtContract.ABI
			transferEvent, err := erc20.Unpack("Transfer", resTx.TxResult.Events[i1].Attributes[0].Value)
			if err != nil {
				t.Error("failed to unpack transfer event", "error", err.Error())
			}
			t.Log(transferEvent)
		}*/
	}
	fmt.Println("TxResult Codespace:", resTx.TxResult.Codespace)
	fmt.Println("TxResult Code:", resTx.TxResult.Code)
	fmt.Println("TxResult Info:", resTx.TxResult.Info)
	fmt.Println("TxResult Data:", string(resTx.TxResult.Data))

	tx, err := txClient.TermintTx2CosmosTx(resTx.Tx) //tx
	if err != nil {
		t.Error("TxDecoder")
		return
	}
	stdTx, err := txClient.ConvertTxToStdTx(tx) //tx
	if err != nil {
		t.Error("txToStdTx")
		return
	}
	//txmsg
	for j := 0; j < len(stdTx.Msgs); j++ {
		msgByte, err := encodingConfig.Marshaler.MarshalInterfaceJSON(stdTx.Msgs[j])
		if err != nil {
			t.Error("MarshalBinaryBare")
			return
		}
		t.Log("message name:", proto.MessageName(stdTx.Msgs[j]))
		t.Log("-------------------------")
		obj := types.MsgEthereumTx{}
		err = encodingConfig.Marshaler.UnmarshalInterfaceJSON(msgByte, &obj)
		if err != nil {
			t.Error("UnmarshalBinaryBare", err)
			return
		}
		t.Log("hash:", obj.Hash)
		t.Log("from:", obj.GetFrom())
		txData, err := types.UnpackTxData(obj.Data)
		if err != nil {
			t.Error("UnpackTxData error")
			return
		}
		t.Log("AccessList:", txData.GetAccessList())
		t.Log("to:", txData.GetTo())
		t.Log("data:", txData.GetData())
		t.Log("amount:", txData.GetValue())
		t.Log("chainID:", txData.GetChainID())
		t.Log("gas:", txData.GetGas())
		t.Log("nonce:", txData.GetNonce())
	}
}
