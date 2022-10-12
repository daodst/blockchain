package client

import (
	"bytes"
	"errors"
	"freemasonry.cc/blockchain/util"
	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"net/url"
	"strconv"
	"strings"
	"time"

	//"fmt"
	"io/ioutil"
	"net/http"
)

type StringEvent struct {
	Type       string      `json:"type,omitempty"`
	Attributes []Attribute `json:"attributes,omitempty"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

type ABCIMessageLog struct {
	MsgIndex uint16 `json:"msg_index"`
	Log      string `json:"log"`

	Events []StringEvent `json:"events"`
}

type TxDetail struct {
	Height string `json:"height"`
	Status string `json:"status"`
	Txhash string `json:"txhash"`
	Error  string `json:"error"`
}

/**
get
*/
func GetRequest(server, params string) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, //true:HTTPTCP（http1.1，）
			MaxIdleConnsPerHost: 512,  //
		},
		Timeout: time.Second * 60, //Client,、response body;Timeout
	}
	bodyReader := strings.NewReader("")
	req, err := http.NewRequest("GET", server+params, bodyReader)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}
	//fmt.Println("body",string(body))
	return string(body), err
}

/**
values  post
*/
func PostValuesRequest(server, url string, values url.Values) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, //true:HTTPTCP（http1.1，）
			MaxIdleConnsPerHost: 512,  //
		},
		Timeout: time.Second * 60, //Client,、response body;Timeout
	}
	req, err := http.NewRequest("POST", server+url, nil)
	req.PostForm = values
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//fmt.Println("body:",len(string(body)))
	if resp.StatusCode != 200 {
		if len(body) == 0 {
			return "", errors.New("" + strconv.Itoa(resp.StatusCode))
		}
		return "", errors.New(string(body))
	}
	//fmt.Println("body",string(body))
	return string(body), err
}

/**
get
*/
func PostRequest(server, url string, params []byte) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, //true:HTTPTCP（http1.1，）
			MaxIdleConnsPerHost: 512,  //
		},
		Timeout: time.Second * 60, //Client,、response body;Timeout
	}
	req, err := http.NewRequest("POST", server+url, bytes.NewReader(params))
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		if len(body) == 0 {
			return "", errors.New("" + strconv.Itoa(resp.StatusCode))
		}
		return "", errors.New(string(body))
	}

	return string(body), err
}

func PostRequestByTimeout(server, url string, params []byte, timeout time.Duration) (string, error) {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, //true:HTTPTCP（http1.1，）
			MaxIdleConnsPerHost: 512,  //
		},
		Timeout: timeout, //Client,、response body;Timeout
	}
	req, err := http.NewRequest("POST", server+url, bytes.NewReader(params))
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	//fmt.Println(resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}
	//fmt.Println("body",string(body))
	return string(body), err
}

/**
tx StdTx 
*/
//func bytesToStdTx(clientCtx client.Context, txhashBytes []byte) (*legacytx.StdTx, string, error) {
//	txhashHex := strings.ToUpper(hex.EncodeToString(txhashBytes))
//	output, err := authclient.QueryTx(clientCtx, txhashHex)
//	if err != nil {
//		fmt.Println("11")
//		return nil, "", err
//	}
//	txBytes := output.Tx.Value
//	tx, err := clientCtx.TxConfig.TxDecoder()(txBytes)
//	if err != nil {
//		fmt.Println("22")
//		return nil, "", err
//	}
//	stdTx, err := txToStdTx(clientCtx, tx)
//	if err != nil {
//		fmt.Println("33")
//		return nil, "", err
//	}
//	return stdTx, txhashHex, nil
//}

func txToStdTx(clientCtx client.Context, tx sdk.Tx) (*legacytx.StdTx, error) {
	signingTx, ok := tx.(signing.Tx)
	if !ok {
		return nil, errors.New("tx to stdtx error")
	}
	stdTx, err := clienttx.ConvertTxToStdTx(clientCtx.LegacyAmino, signingTx)
	if err != nil {
		return nil, err
	}
	return &stdTx, nil
}

//msg
func unmarshalMsg(msg sdk.Msg, obj interface{}) error {
	msgByte, err := util.Json.Marshal(msg)
	if err != nil {
		return err
	}
	return util.Json.Unmarshal(msgByte, &obj)
}
