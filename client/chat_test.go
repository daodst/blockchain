package client

import (
	"testing"
)


func TestQueryChatInfo(t *testing.T) {

	txClient := NewTxClient()
	accClient := NewAccountClient(&txClient)
	chatClient := NewChatClient(&txClient, &accClient)

	chatInfo, err := chatClient.QueryUserInfo("dex1nwp7ggg09px3ese9hy3a52fcrljlexu5xp5mx7")
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
		t.Log("ChatWhitelist:", userinfo.ChatWhitelist)
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

	userInfo, err := chatClient.QueryUserByMobile("123456789")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("userInfo.FromAddress:", userInfo.FromAddress)
	t.Log("userInfo.NodeAddress:", userInfo.NodeAddress)
	t.Log("userInfo.AddressBook:", userInfo.AddressBook)
	t.Log("userInfo.ChatBlacklist:", userInfo.ChatBlacklist)
	t.Log("userInfo.ChatWhitelist:", userInfo.ChatWhitelist)
	t.Log("userInfo.ChatReceiveGiftInfo:", userInfo.Mobile)
	t.Log("userInfo.ChatFee:", userInfo.UpdateTime)
	t.Log("userInfo.PledgeLevel:", userInfo.PledgeLevel)                 
	t.Log("userInfo.GatewayProfixMobile:", userInfo.GatewayProfixMobile) 
	t.Log("userInfo.IsExist:", userInfo.IsExist)                         //0: 1:
}

func TestQueryAddrByChatAddr(t *testing.T) {

	txClient := NewTxClient()
	accClient := NewAccountClient(&txClient)
	chatClient := NewChatClient(&txClient, &accClient)

	res, err := chatClient.QueryAddrByChatAddr("dex13zkyjmlseg5ms3lxhq7nwt5jlw8yzc2d3hrh0v")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}
