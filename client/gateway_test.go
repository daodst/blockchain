package client

import "testing"

//
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
