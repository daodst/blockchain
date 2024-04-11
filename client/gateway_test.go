package client

import (
	"testing"
)


func TestGasPrice(t *testing.T) {
	txClient := NewTxClient()
	gasPrice, err := txClient.QueryGasPrice()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("gas:", gasPrice)
	if gasPrice.IsZero() {
		t.Log("11111111111111111")
	}
}
func TestGatewayList(t *testing.T) {

	txClient := NewTxClient()
	gatewayclient := NewGatewayClinet(&txClient)
	gatewayList, err := gatewayclient.QueryGatewayList()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("GatewayQuota:", gatewayList[0].GatewayQuota)
	t.Log("GatewayNum:", gatewayList[0].GatewayNum)
	t.Log("GatewayName:", gatewayList[0].GatewayName)
}

func TestGatewayRegister(t *testing.T) {

	txClient := NewTxClient()
	gatewayclient := NewGatewayClinet(&txClient)
	accClient := NewAccountClient(&txClient)
	fromWallet, _ := accClient.CreateAccountFromSeed("plate vintage fortune awesome lounge mule rough unaware echo stem giraffe icon usual resource craft disease truck arm reason announce cargo word gloom bid")

	result, err := gatewayclient.GatewayRegister("dst1jxhtvfqcptq2gxmmeuh6uspr033jzzaycwckym", "http://192.168.10.136:50327", "", fromWallet.PrivateKey, "fc53f36f89e1ecc6f4978cdba1e37d95.d95", "12D3KooWSfGTq1Jy4t3LovN2d3WbKysdBUQgfCD3FhuyK4PvKJsS", "dst15zhu6jykcjwv8aeg85wunmrmjn0t82x9gyjz7k", []string{"1111111"})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Code:", result.Code)
	t.Log("TxHash:", result.TxHash)
}
func TestGatewayAddIndexNum(t *testing.T) {
	txClient := NewTxClient()
	gatewayclient := NewGatewayClinet(&txClient)
	result, err := gatewayclient.GatewayAddIndexNum("dst104gnuqqvuhfcudea3e2f7a30quvkm3tqlz7se4", "dstvaloper104gnuqqvuhfcudea3e2f7a30quvkm3tqcvwzue", "8983569bd77349789d85a485da8acf4f97fac6800fe0b4f0351fd5fdb4708e0f", []string{"1888888"})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Code:", result.Code)
	t.Log("TxHash:", result.TxHash)
}
func TestGatewayUnDelegate(t *testing.T) {
	txClient := NewTxClient()
	gatewayclient := NewGatewayClinet(&txClient)
	result, err := gatewayclient.GatewayUnDelegate("dst104gnuqqvuhfcudea3e2f7a30quvkm3tqlz7se4", "dstvaloper1zl5quqaukt4ssks8j2nr6rvz472rl67de2rfe6", "970000000000000000000fm", "8983569bd77349789d85a485da8acf4f97fac6800fe0b4f0351fd5fdb4708e0f", []string{"1234567"})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Code:", result.Code)
	t.Log("TxHash:", result.TxHash)
}
func TestQueryGatewayNumList(t *testing.T) {
	txClient := NewTxClient()
	gatewayclient := NewGatewayClinet(&txClient)
	result, err := gatewayclient.QueryGatewayNumList()
	if err != nil {
		t.Error(err)
		return
	}

	for s, val := range result {
		t.Log("key:", s)
		t.Log("value:", val)
	}
}

func TestQueryGatewayClusters(t *testing.T) {
	txClient := NewTxClient()
	gatewayclient := NewGatewayClinet(&txClient)
	result, err := gatewayclient.QueryGatewayClusters("dstvaloper10x9a5ym9hunn9splh0y3pfhakc89hxnvjwan6z")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(result)
}
