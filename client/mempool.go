package client

import (
	"context"
	"encoding/hex"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/tendermint/tendermint/crypto/tmhash"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	ttypes "github.com/tendermint/tendermint/types"
)

//
type MempoolClient struct {
	ServerUrl string
	logPrefix string
}

func (this *MempoolClient) TxExisted(tx ttypes.Tx) (res *ctypes.ResultTxExisted, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return res, err
	}
	return node.TxExisted(context.Background(), tx)
}

//tx，，gas
func (this *MempoolClient) CheckTx(tx ttypes.Tx) (resultTx *ctypes.ResultCheckTx, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return resultTx, err
	}
	return node.CheckTx(context.Background(), tx)
}

// 
func (this *MempoolClient) UnconfirmedTxs() (list []*legacytx.StdTx, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return list, err
	}
	limit := 100
	resp, err := node.UnconfirmedTxs(context.Background(), &limit)
	if err != nil {
		log.WithError(err).Error("UnconfirmedTxs")
		return list, err
	}
	for _, tx := range resp.Txs {
		cosmosTx, err := termintTx2CosmosTx(tx)
		if err != nil {
			log.WithError(err).Error("TermintTx2CosmosTx1")
			return list, err
		}
		stdTx, err := convertTxToStdTx(cosmosTx)
		if err != nil {
			log.WithError(err).Error("ConvertTxToStdTx1")
			return list, err
		}
		stdTx.Memo = hex.EncodeToString(tmhash.Sum(tx))
		list = append(list, stdTx)
	}
	return list, nil
}

// 
func (this *MempoolClient) GetPoolFirstAndLastTxs() (frist *legacytx.StdTx, last *legacytx.StdTx, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return frist, last, err
	}
	resp, err := node.PoolFirstAndLastTxs(context.Background())
	if err != nil {
		log.WithError(err).Error("UnconfirmedTxs")
		return frist, last, err
	}
	var cosmosTx sdk.Tx
	if resp.First != nil {
		cosmosTx, err = termintTx2CosmosTx(resp.First)
		if err != nil {
			log.WithError(err).Error("TermintTx2CosmosTx1")
			return frist, last, err
		}
		frist, err = convertTxToStdTx(cosmosTx)
		if err != nil {
			log.WithError(err).Error("ConvertTxToStdTx1")
			return frist, last, err
		}
	}
	if resp.Last != nil {
		cosmosTx, err = termintTx2CosmosTx(resp.Last)
		if err != nil {
			log.WithError(err).Error("TermintTx2CosmosTx2")
			return frist, last, err
		}
		last, err = convertTxToStdTx(cosmosTx)
		if err != nil {
			log.WithError(err).Error("ConvertTxToStdTx2")
			return frist, last, err
		}
	}

	return frist, last, nil
}
