package client

import "testing"


func TestQuerySendGift(t *testing.T) {

	txClient := NewTxClient()
	accClient := NewAccountClient(&txClient)
	chatClient := NewChatClient(&txClient, &accClient)

	isPay, err := chatClient.QueryChatSendGift("dex14t2m2mhhw35nk3w0qjyvdtcynul4t5tdhlhqp0", "dex1xwrc4d5p78chq04sxmd3gm0vz6tf0c02va46rt")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(isPay)
}


func TestQueryAllPledgeInfo(t *testing.T) {

	txClient := NewTxClient()
	accClient := NewAccountClient(&txClient)
	chatClient := NewChatClient(&txClient, &accClient)

	allPledgeInfo, err := chatClient.QueryPledgeInfo("dex14t2m2mhhw35nk3w0qjyvdtcynul4t5tdhlhqp0")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(allPledgeInfo)
}


func TestQueryChatInfo(t *testing.T) {

	txClient := NewTxClient()
	accClient := NewAccountClient(&txClient)
	chatClient := NewChatClient(&txClient, &accClient)

	chatInfo, err := chatClient.QueryUserInfo("dex14t2m2mhhw35nk3w0qjyvdtcynul4t5tdhlhqp0")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(chatInfo)
}


func TestQueryChatInfos(t *testing.T) {

	txClient := NewTxClient()
	accClient := NewAccountClient(&txClient)
	chatClient := NewChatClient(&txClient, &accClient)

	chatInfos, err := chatClient.QueryUserInfos([]string{"dex14t2m2mhhw35nk3w0qjyvdtcynul4t5tdhlhqp0", "dex1xwrc4d5p78chq04sxmd3gm0vz6tf0c02va46rt"})
	if err != nil {
		t.Fatal(err)
	}

	for _, userinfo := range chatInfos {
		t.Log("FromAddress:", userinfo.FromAddress)
		t.Log("NodeAddress:", userinfo.NodeAddress)
		t.Log("AddressBook:", userinfo.AddressBook)
		t.Log("ChatBlacklist:", userinfo.ChatBlacklist)
		t.Log("ChatRestrictedMode:", userinfo.ChatRestrictedMode)
		t.Log("ChatWhitelist:", userinfo.ChatWhitelist)
		t.Log("ChatSendGiftInfo:", userinfo.ChatFee)
		t.Log("ChatReceiveGiftInfo:", userinfo.Mobile)
		t.Log("ChatFee:", userinfo.UpdateTime)
		t.Log("PledgeLevel:", userinfo.PledgeLevel)                 
		t.Log("GatewayProfixMobile:", userinfo.GatewayProfixMobile) 
		t.Log("IsExist:", userinfo.IsExist)                         //0: 1:
	}
}


func TestQueryUserInfoByMobile(t *testing.T) {

	txClient := NewTxClient()
	accClient := NewAccountClient(&txClient)
	chatClient := NewChatClient(&txClient, &accClient)

	userInfo, err := chatClient.QueryUserByMobile("12345670004")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("userInfo.FromAddress:", userInfo.FromAddress)
	t.Log("userInfo.NodeAddress:", userInfo.NodeAddress)
	t.Log("userInfo.AddressBook:", userInfo.AddressBook)
	t.Log("userInfo.ChatBlacklist:", userInfo.ChatBlacklist)
	t.Log("userInfo.ChatRestrictedMode:", userInfo.ChatRestrictedMode)
	t.Log("userInfo.ChatWhitelist:", userInfo.ChatWhitelist)
	t.Log("userInfo.ChatSendGiftInfo:", userInfo.ChatFee)
	t.Log("userInfo.ChatReceiveGiftInfo:", userInfo.Mobile)
	t.Log("userInfo.ChatFee:", userInfo.UpdateTime)
	t.Log("userInfo.PledgeLevel:", userInfo.PledgeLevel)                 
	t.Log("userInfo.GatewayProfixMobile:", userInfo.GatewayProfixMobile) 
	t.Log("userInfo.IsExist:", userInfo.IsExist)                         //0: 1:
}


func TestQuerySendGifts(t *testing.T) {

	txClient := NewTxClient()
	accClient := NewAccountClient(&txClient)
	chatClient := NewChatClient(&txClient, &accClient)

	isPays, err := chatClient.QueryChatSendGifts("dex14t2m2mhhw35nk3w0qjyvdtcynul4t5tdhlhqp0", []string{"dex1xwrc4d5p78chq04sxmd3gm0vz6tf0c02va46rt"})
	if err != nil {
		t.Fatal(err)
	}

	for address, isPay := range isPays {
		t.Log("address:", address)
		t.Log("isPay:", isPay)
	}
}
