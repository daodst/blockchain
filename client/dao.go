package client

import (
	"encoding/hex"
	"errors"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	ttypes "github.com/tendermint/tendermint/types"
)

type DaoClient struct {
	TxClient  *TxClient
	key       *SecretKey
	ServerUrl string
	logPrefix string
}

// 
func (this *DaoClient) OracleMachineUpload(address,gatewayAddress string,rates []types.ColonyRate,privateKey string, fee legacytx.StdFee) (tx ttypes.Tx, resp *core.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	msg := types.NewMsgColonyRate(address,gatewayAddress,rates)
	if err != nil {
		log.Error("NewMsgTransfer")
		return
	}
	var result *core.BaseResponse
	//
	tx, result, err = this.TxClient.SignAndSendMsg(address, privateKey, fee, "", msg)
	if err != nil {
		return
	}
	resp = new(core.BroadcastTxResponse)
	//
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return tx, resp, nil
	} else {
		//
		resp.TxHash = hex.EncodeToString(tx.Hash())
		return tx, resp, errors.New(result.Info)
	}
}