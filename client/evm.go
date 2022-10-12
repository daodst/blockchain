package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/client/evm"
	"freemasonry.cc/blockchain/core"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

//
func (this *EvmClient) NetVersion() (string, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	var res string

	rpcRes, err := this.Call("net_version", []string{})
	if err != nil {
		log.WithError(err).Error("call")
		return res, err
	}

	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, nil
}

//
func (this *EvmClient) NetListening() (bool, error) {
	var res bool
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	rpcRes, err := this.Call("net_listening", []string{})
	if err != nil {
		log.WithError(err).Error("call")
		return res, err
	}

	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, err
}

//
func (this *EvmClient) NetPeerCount() (int, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	var res int
	rpcRes, err := this.Call("net_peerCount", []string{})
	if err != nil {
		log.WithError(err).Error("call")
		return res, err
	}

	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, err
}

//
// blockNumber 
//  
//	last    @
//      @
// ----------------------------
// fullTx Tx   true   false 
func (this *EvmClient) GetBlockNumber(blockNumber string, fullTx bool) (map[string]interface{}, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	var res map[string]interface{}
	rpcRes, err := this.Call("eth_getBlockByNumber", []interface{}{blockNumber, true})
	if err != nil {
		log.WithError(err).Error("call")
		return res, err
	}

	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, err
}

//
func (this *EvmClient) BlockNumber() (hexutil.Big, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	var res hexutil.Big
	rpcRes, err := this.Call("eth_blockNumber", []interface{}{})
	if err != nil {
		log.WithError(err).Error("call")
		return res, err
	}
	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, err
}

//
func (this *EvmClient) GetBalance(addr string) (hexutil.Big, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	var res hexutil.Big
	rpcRes, err := this.Call("eth_getBalance", []string{addr, "latest"})
	if err != nil {
		log.WithField("address", addr).WithError(err).Error("call")
		return res, err
	}
	if rpcRes.Error != nil {
		log.WithFields(logrus.Fields{"code": rpcRes.Error.Code, "message": rpcRes.Error.Message, "data": rpcRes.Error.Data}).Error("rpcError")
		return res, errors.New(rpcRes.Error.Message)
	}

	err = res.UnmarshalJSON(rpcRes.Result)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("UnmarshalJSON")
		return res, err
	}
	return res, nil
}

//hash
func (this *EvmClient) GetTransactionReceipt(hash string) (map[string]interface{}, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	var res map[string]interface{}
	rpcRes, err := this.Call("eth_getTransactionReceipt", []interface{}{hash})
	if err != nil {
		log.WithField("hash", hash).WithError(err).Error("call")
		return res, err
	}
	if rpcRes.Error != nil {
		log.WithFields(logrus.Fields{"code": rpcRes.Error.Code, "message": rpcRes.Error.Message, "data": rpcRes.Error.Data}).Error("rpcError")
		return res, errors.New(rpcRes.Error.Message)
	}
	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, err
}

//Hash
func (this *EvmClient) GetTransactionByHash(hash string) (map[string]interface{}, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	var res map[string]interface{}
	rpcRes, err := this.Call("eth_getTransactionByHash", []interface{}{hash})
	if err != nil {
		log.WithField("hash", hash).WithError(err).Error("call")
		return res, err
	}
	if rpcRes.Error != nil {
		log.WithFields(logrus.Fields{"code": rpcRes.Error.Code, "message": rpcRes.Error.Message, "data": rpcRes.Error.Data}).Error("rpcError")
		return res, errors.New(rpcRes.Error.Message)
	}
	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, err
}

//
func (this *EvmClient) GetAddress() ([]hexutil.Bytes, error) {
	rpcRes, err := this.CallWithError("eth_accounts", []string{})
	if err != nil {
		return nil, err
	}
	var res []hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *EvmClient) CreateRequest(method string, params interface{}) evm.Request {
	return evm.Request{
		Version: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
}

func (this *EvmClient) CallWithError(method string, params interface{}) (*evm.Response, error) {
	req, err := json.Marshal(this.CreateRequest(method, params))
	if err != nil {
		return nil, err
	}

	var rpcRes *evm.Response
	time.Sleep(1 * time.Second)

	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", core.EvmRpcURL, bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, errors.New("Could not perform request")
	}

	decoder := json.NewDecoder(res.Body)
	rpcRes = new(evm.Response)
	err = decoder.Decode(&rpcRes)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	if rpcRes.Error != nil {
		return nil, fmt.Errorf(rpcRes.Error.Message)
	}

	return rpcRes, nil
}

func (this *EvmClient) Call(method string, params interface{}) (*evm.Response, error) {
	req, err := json.Marshal(this.CreateRequest(method, params))
	if err != nil {
		return nil, err
	}

	var rpcRes *evm.Response
	time.Sleep(1 * time.Second)

	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", this.RpcUrl, bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)
	rpcRes = new(evm.Response)
	err = decoder.Decode(&rpcRes)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return rpcRes, nil
}

type EvmClient struct {
	RpcUrl string
}
