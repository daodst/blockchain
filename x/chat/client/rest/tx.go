package rest

import (
	"errors"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/trerr"
	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/gogo/protobuf/proto"
	"io/ioutil"
	"net/http"
)

//tx，
func BroadcastTxHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := core.BuildLog("BroadcastTxHandlerFn", core.LmChainRest)
		var txBytes []byte
		if r.Body != nil {
			txBytes, _ = ioutil.ReadAll(r.Body)
		}
		baseResponse := core.BaseResponse{}
		txResponse := core.BroadcastTxResponse{}
		tx, _ := clientCtx.TxConfig.TxDecoder()(txBytes) //tx
		stdTx, err := txToStdTx(clientCtx, tx)           //tx  stdTx
		if err != nil {
			baseResponse.Info = err.Error()
			SendReponse(w, clientCtx, baseResponse)
			return
		}
		msgs := stdTx.GetMsgs() //tx
		fee := stdTx.Fee        //tx
		memo := stdTx.Memo      //tx

		log.Info("")

		//
		err = broadcastMsgCheck(msgs, fee, memo)
		if err != nil {
			errmsg := trerr.TransError(err.Error())
			baseResponse.Info = errmsg.Error()
			SendReponse(w, clientCtx, baseResponse)
			return
		}

		res, err := clientCtx.BroadcastTx(txBytes)
		if err != nil {
			baseResponse.Info = err.Error()
			SendReponse(w, clientCtx, baseResponse)
			return
		}

		if res.Code == 0 { //code=0 
			baseResponse.Status = 1
		} else {
			baseResponse.Status = 0
		}
		txResponse.Code = res.Code
		txResponse.CodeSpace = res.Codespace
		txResponse.TxHash = res.TxHash
		txResponse.Height = res.Height

		baseResponse.Info = parseErrorCode(res.Code, res.Codespace, res.RawLog)
		baseResponse.Data = txResponse
		SendReponse(w, clientCtx, baseResponse)
	}
}

//code、codeSpace、rowlog 
func parseErrorCode(code uint32, codeSpace string, rowlog string) string {
	if codeSpace == sdkErrors.RootCodespace {
		if code == sdkErrors.ErrInsufficientFee.ABCICode() { //
			return FeeIsTooLess
		} else if code == sdkErrors.ErrOutOfGas.ABCICode() { //gas
			return ErrorGasOut
		} else if code == sdkErrors.ErrUnauthorized.ABCICode() { //idaccount number 
			return ErrUnauthorized
		} else if code == sdkErrors.ErrWrongSequence.ABCICode() { //
			return ErrWrongSequence
		}
	}
	return rowlog
}

func broadcastMsgCheck(msgs []sdk.Msg, fee legacytx.StdFee, memo string) (err error) {
	for _, msg := range msgs {
		msgType := proto.MessageName(msg)
		if txHandles.HaveRegistered(msgType) { //
			msgByte, err := util.Json.Marshal(msg)
			if err != nil {
				return err
			}
			err = txHandles.Handle(msgType, msgByte, fee, memo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func txToStdTx(clientCtx client.Context, tx sdk.Tx) (*legacytx.StdTx, error) {
	signingTx, ok := tx.(signing.Tx)
	if !ok {
		return nil, errors.New("txTostdtx error")
	}
	stdTx, err := clienttx.ConvertTxToStdTx(clientCtx.LegacyAmino, signingTx)
	if err != nil {
		return nil, err
	}
	return &stdTx, nil
}
