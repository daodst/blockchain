package client

import (
	"encoding/json"
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/core/abi/erc20"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"strings"
	"testing"
)

// tx
func TestTxFind(t *testing.T) {
	txClient := NewTxClient()
	//tx
	txHash := "B74F7FF6D0E9483092EE465317989F53A9AA489D6D674E4C1DB74A0815592396"
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

				if log.Topics[0] == erc20.FunctionHashTransfer { //event hash
					from := common.HexToAddress(log.Topics[1])
					to := common.HexToAddress(log.Topics[2])

					// abi
					transferABI, err := abi.JSON(strings.NewReader(erc20.TransferEventAbi))
					if err != nil {
						t.Error("abi", err)
					}
					amount, err := transferABI.Unpack("Transfer", log.Data)
					if err != nil {
						t.Error("", err)
					}
					fmt.Println(":", from.Hex())
					fmt.Println(":", to.Hex())
					fmt.Println(":", amount)
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

	//tx, err := txClient.TermintTx2CosmosTx(resTx.Tx) //tx
	//if err != nil {
	//	t.Error("TxDecoder")
	//	return
	//}
	//stdTx, err := txClient.ConvertTxToStdTx(tx) //tx
	//if err != nil {
	//	t.Error("txToStdTx")
	//	return
	//}
	////txmsg
	//for j := 0; j < len(stdTx.Msgs); j++ {
	//	msgByte, err := encodingConfig.Marshaler.MarshalInterfaceJSON(stdTx.Msgs[j])
	//	if err != nil {
	//		t.Error("MarshalBinaryBare")
	//		return
	//	}
	//	t.Log("message name:", proto.MessageName(stdTx.Msgs[j]))
	//	t.Log("-------------------------")
	//	obj := types.MsgEthereumTx{}
	//	err = encodingConfig.Marshaler.UnmarshalInterfaceJSON(msgByte, &obj)
	//	if err != nil {
	//		t.Error("UnmarshalBinaryBare", err)
	//		return
	//	}
	//	t.Log("hash:", obj.Hash)
	//	t.Log("from:", obj.GetFrom())
	//	txData, err := types.UnpackTxData(obj.Data)
	//	if err != nil {
	//		t.Error("UnpackTxData error")
	//		return
	//	}
	//	t.Log("AccessList:", txData.GetAccessList())
	//	t.Log("to:", txData.GetTo())
	//	t.Log("data:", txData.GetData())
	//	t.Log("amount:", txData.GetValue())
	//	t.Log("chainID:", txData.GetChainID())
	//	t.Log("gas:", txData.GetGas())
	//	t.Log("nonce:", txData.GetNonce())
	//}
}

func TestFindProposal(t *testing.T) {
	txClient := NewTxClient()
	proposal, err := txClient.QueryProposer(1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(proposal.GetContent().ProposalType())
	t.Log(proposal.GetContent().(*upgradetypes.SoftwareUpgradeProposal).Plan.Info)
}

func TestUnjail(t *testing.T) {
	txClient := NewTxClient()
	dposClient := NewDposClinet(&txClient)
	a, err := dposClient.UnjailValidator("dst1kxcnfvtp6ep042vudctgda29cywkpzp0ttky75", "dstvaloper1kxcnfvtp6ep042vudctgda29cywkpzp0v9xkmc", "349b21b102c799984d0a2d29986fee00f36abcaaa51f10bb196f44dd5a4614a6")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(err)
	t.Log(a)
}
func TestUnjail1(t *testing.T) {
	txClient := NewTxClient()
	a, err := txClient.ChatTokenIssue("dex1khy3drtnyq4un5q2tft9pvq2g6l2a483h03js9", "87858601d5ca0a15980fd2e5e852efe48fc85de37379def0fc5f3ce003d2b192")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(err)
	t.Log(a)

}

func TestVal2Acc(t *testing.T) {
	valS := "dexvaloper1khy3drtnyq4un5q2tft9pvq2g6l2a48332zydc"
	val, err := sdk.ValAddressFromBech32(valS)
	if err != nil {
		t.Error(err)
	}

	acc := sdk.AccAddress(val)

	t.Log(common.BytesToAddress(acc.Bytes()))

}

func TestQuery1(t *testing.T) {

	txclient := NewTxClient()
	accclient := NewAccountClient(&txclient)
	chatClient := NewChatClient(&txclient, &accclient)

	resm, err := chatClient.QueryUserInfo("dex1kyrm5eydks4vxxhrc2h2p7c6gk6apu84zjhyax")
	t.Log(resm)
	t.Log(err)
}

func TestGasInfo(t *testing.T) {
	txclient := NewTxClient()
	seq, err := txclient.FindAccountNumberSeq("dex1acxdjanpgkeheh5pk8xq6rqhpjc407d3dke0dt")
	if err != nil {
		t.Error(err)
		return
	}
	msg, err := txclient.GasMsg("proposal_params", `{"proposer":"dex1acxdjanpgkeheh5pk8xq6rqhpjc407d3dke0dt","deposit":"100fm","title":"","description":"","change":"[{\"subspace\":\"staking\",\"key\":\"MaxValidators\",\"value\":\"110\"}]"}`)
	if err != nil {
		t.Error(err)
		return
	}
	fee, gas, _, err := txclient.GasInfo(seq, msg)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(fee)
	t.Log(gas)
}

func TestTxDecoder(t *testing.T) {
	txclient := NewTxClient()
	msgs := []sdk.Msg{}
	msg := types.NewMsgClusterAddMembers("dst1kxcnfvtp6ep042vudctgda29cywkpzp0ttky75", "sfsdafasdfasdfasdfasdfasdfsadf", nil)
	msgs = append(msgs, msg)
	seq := core.AccountNumberSeqResponse{}
	signedTx, err := txclient.SignTx("349b21b102c799984d0a2d29986fee00f36abcaaa51f10bb196f44dd5a4614a6", seq, core.NewLedgerFeeZero(), "", msgs...)
	if err != nil {
		t.Error(err)
		return
	}
	
	txBytes, err := txclient.SignTx2Bytes(signedTx)
	if err != nil {
		t.Error(err)
		return
	}
	tx, err := clientCtx.TxConfig.TxDecoder()(txBytes)
	if err != nil {
		t.Error(err)
		return
	}
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		t.Error("Tx must be a FeeTx")
		return
	}

	msgs = feeTx.GetMsgs()
	for _, ms := range msgs {
		msgByte, err := util.Json.Marshal(ms)
		if err != nil {
			t.Error(err)
			return
		}
		obj := types.MsgClusterAddMembers{}
		err = util.Json.Unmarshal(msgByte, &obj)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(obj.FromAddress)
		t.Log(obj.ClusterId)
	}
}
