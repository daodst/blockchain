package client


//func TestGetPoolFirstAndLastTxs(t *testing.T) {
//	poolClient := NewMempoolClient()
//	for {

//		first, last, err := poolClient.GetPoolFirstAndLastTxs()
//		if err != nil {
//			t.Error("error:", err)
//			return
//		}
//		if first != nil || last != nil {
//			t.Log("--------------------")
//		}else{
//			t.Log("wait...")
//		}
//		if first != nil {
//			t.Log("first gasLimit:", first.Fee.GetGas())
//			t.Log("first gasAmount:", first.Fee.GetAmount())
//			t.Log("first gasPrices:", first.Fee.GasPrices())
//		}
//		if last != nil {
//			t.Log("last gasLimit:", last.Fee.GetGas())
//			t.Log("last gasAmount:", last.Fee.GetAmount())
//			t.Log("last gasPrices:", last.Fee.GasPrices())
//		}
//		<-time.After(time.Second * 3)
//	}
//}


//func TestUnconfirmedTxs(t *testing.T) {
//	poolClient := NewMempoolClient()
//	for {
//		t.Log("--------------------------------")
//		list, err := poolClient.UnconfirmedTxs()
//		if err != nil {
//			t.Error("error:", err)
//			return
//		}
//		for index, tx := range list {

//			t.Log(index, " signers:", tx.GetMsgs())
//			t.Log(index, "  ", tx.Memo, tx.Fee.GetAmount(), "/", tx.Fee.GetGas(), "= ", tx.Fee.GasPrices())
//		}
//		<-time.After(time.Second * 2)
//	}
//}

////tx
//func TestTxExisted(t *testing.T) {
//	poolClient := NewMempoolClient()

//	//txhash := strings.ToUpper("1ea1adb379e987e41f19d600f29219e027b36e6a79ac96c460b867cbf902e561")

//	node, err := clientCtx.GetNode()
//	if err != nil {
//		t.Error("error:", err)
//		return
//	}

//	limit := 100

//	for {
//		t.Log("-----------------------------------------")
//		resp, err := node.UnconfirmedTxs(context.Background(), &limit)
//		if err != nil {
//			t.Error("error:", err)
//			return
//		}

//		for _, tx := range resp.Txs {
//			res, err := poolClient.TxExisted(tx)
//			if err != nil {
//				t.Error("error:", err)
//				return
//			}
//			t.Log("TxHash:", hex.EncodeToString(tx.Hash()), " TxExisted:", res.Existed)
//		}
//		<-time.After(time.Second * 2)
//	}
//}
