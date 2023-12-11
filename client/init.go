package client

import (
	"encoding/json"
	"fmt"
	cmdcfg "freemasonry.cc/blockchain/cmd/config"
	daotypes "freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)


var msgUnmashalHandles map[string]func(msgByte []byte) (sdk.Msg, error)

// tx
func registerUnmashalHandles(msgType string, callback func(msgByte []byte) (sdk.Msg, error), msgJson func() ([]byte, error)) {
	if msgJson != nil {
		jjson, _ := msgJson()
		fmt.Println(msgType, "=", string(jjson))
	}

	msgUnmashalHandles[msgType] = callback
}

func init() {
	// set the address prefixes
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	// TODO fix
	// if err := cmdcfg.EnableObservability(); err != nil {
	// 	panic(err)
	// }
	cmdcfg.SetBip44CoinType(config)
	config.Seal()

	msgUnmashalHandles = make(map[string]func(data []byte) (sdk.Msg, error))

	registerUnmashalHandles("cosmos-sdk/MsgWithdrawDelegationReward", unmashalDistMsgWithdrawDelegatorReward, func() ([]byte, error) {
		return json.Marshal(&distributionTypes.MsgWithdrawDelegatorReward{DelegatorAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet", ValidatorAddress: "dstvaloper1emqqtns9xdrpyjfant4wzf7zsylc59ydypaej5"})
	})
	registerUnmashalHandles("cosmos-sdk/MsgSend", unmashalMsgSend, nil)
	registerUnmashalHandles("cosmos-sdk/MsgDelegate", unmashalMsgDelegate, nil)
	registerUnmashalHandles("cosmos-sdk/MsgVote", unmashalMsgVote, nil)

	registerUnmashalHandles("dao/MsgColonyRate", unmashalDaoMsgColonyRate, nil)

	
	registerUnmashalHandles("dao/MsgCreateCluster", unmashalDaoMsgCreateCluster, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgCreateCluster{
			FromAddress:  "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			GateAddress:  "dstvaloper1emqqtns9xdrpyjfant4wzf7zsylc59ydypaej5",
			ClusterId:    "123455",
			DeviceRatio:  sdk.NewDec(1),
			SalaryRatio:  sdk.NewDec(1),
			BurnAmount:   sdk.NewDec(1),
			ChatAddress:  "",
			ClusterName:  "",
			FreezeAmount: sdk.NewDec(1),
			Metadata:     ""})
	})

	
	registerUnmashalHandles("dao/MsgClusterAddMembers", unmashalDaoMsgClusterAddMembers, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgClusterAddMembers{
			FromAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			ClusterId:   "123455",
			Members:     []daotypes.Members{{MemberAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet", IndexNum: "", ChatAddress: ""}}})
	})

	
	registerUnmashalHandles("dao/MsgDeleteMembers", unmashalDaoMsgDeleteMembers, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgDeleteMembers{
			FromAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			ClusterId:   "123455",
			Members:     []string{"dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet", "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet"}})
	})

	
	registerUnmashalHandles("dao/MsgClusterChangeName", unmashalDaoMsgClusterChangeName, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgClusterChangeName{
			FromAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			ClusterId:   "123455",
			ClusterName: ""})
	})

	
	registerUnmashalHandles("dao/MsgClusterMemberExit", unmashalDaoMsgClusterMemberExit, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgClusterMemberExit{
			FromAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			ClusterId:   "123455"})
	})

	
	registerUnmashalHandles("dao/MsgBurnToPower", unmashalDaoMsgBurnToPower, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgBurnToPower{
			FromAddress:     "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			ToAddress:       "dstvaloper1emqqtns9xdrpyjfant4wzf7zsylc59ydypaej5",
			BurnAmount:      sdk.NewDec(1),
			UseFreezeAmount: sdk.NewDec(1),
			GatewayAddress:  "dstvoloperxxxxxxxxxxxxxxxx",
			ChatAddress:     "",
			ClusterId:       "123455"})
	})

	
	registerUnmashalHandles("dao/MsgClusterChangeDeviceRatio", unmashalDaoMsgClusterChangeDeviceRatio, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgClusterChangeDeviceRatio{
			FromAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			DeviceRatio: sdk.NewDec(1),
			ClusterId:   "123455"})
	})

	
	registerUnmashalHandles("dao/MsgClusterChangeSalaryRatio", unmashalDaoMsgClusterChangeSalaryRatio, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgClusterChangeSalaryRatio{
			FromAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			SalaryRatio: sdk.NewDec(1),
			ClusterId:   "123455"})
	})

	// id
	registerUnmashalHandles("dao/MsgClusterChangeId", unmashalDaoMsgClusterChangeId, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgClusterChangeId{
			FromAddress:  "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			NewClusterId: "id",
			ClusterId:    "123455"})
	})

	
	registerUnmashalHandles("dao/MsgWithdrawOwnerReward", unmashalDaoMsgWithdrawOwnerReward, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgWithdrawOwnerReward{
			ClusterId: "!zOnfR8GEhgtDGYP5:1111111.nxn",
			Address:   "dstxxxxxxxxxxxxxxxxx",
		})
	})
	
	registerUnmashalHandles("dao/MsgWithdrawBurnReward", unmashalDaoMsgWithdrawBurnReward, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgWithdrawBurnReward{
			ClusterId:     "!zOnfR8GEhgtDGYP5:1111111.nxn",
			MemberAddress: "dstxxxxxxxxxxxxxxxxx",
		})
	})

	
	registerUnmashalHandles("dao/MsgWithdrawDeviceReward", unmashalDaoMsgWithdrawDeviceReward, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgWithdrawDeviceReward{
			ClusterId:     "!zOnfR8GEhgtDGYP5:1111111.nxn",
			MemberAddress: "dstxxxxxxxxxxxxxxxxx",
		})
	})

	
	registerUnmashalHandles("dao/MsgThawFrozenPower", unmashalDaoMsgThawFrozenPower, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgThawFrozenPower{
			FromAddress:    "dstxxxxxxxxxxxxxxxxx",
			ClusterId:      "!zOnfR8GEhgtDGYP5:1111111.nxn",
			ThawAmount:     sdk.MustNewDecFromStr("12.34"),
			GatewayAddress: "dstveloperxxxxxxxxxxxxxxxxxxxxxxxx",
			ChatAddress:    "dstxxxxxxxxxxxxxxxxx",
		})
	})

	//registerUnmashalHandles("pledge/MsgPledge",unmashalPledgeMsgPledge,nil)   
	//registerUnmashalHandles("pledge/MsgUnpledge",unmashalPledgeMsgUnpledge,nil)  //ji
	//registerUnmashalHandles("pledge/MsgWithdrawDelegatorReward",unmashalPledgeMsgWithdrawDelegatorReward,nil) 
	//registerUnmashalHandles("pledge/MsgCreateDeviceCluster",unmashalMsgCreateDeviceCluster,nil)  
	//registerUnmashalHandles("pledge/MsgChangeClusterInfo",unmashalMsgChangeClusterInfo,nil) 
	//registerUnmashalHandles("pledge/MsgWithdrawValidatorCommission",unmashalPledgeMsgWithdrawValidatorCommission,nil) 

	registerUnmashalHandles("chat/MsgMobileTransfer", unmashalMsgMobileTransfer, nil)
	registerUnmashalHandles("chat/MsgBurnGetMobile", unmashalMsgBurnGetMobile, nil)
	registerUnmashalHandles("chat/MsgSetChatInfo", unmashalMsgSetChatInfo, nil)

	registerUnmashalHandles("gateway/MsgCreateSmartValidator", unmashalMsgCreateSmartValidator, nil)     
	registerUnmashalHandles("gateway/MsgGatewayRegister", unmashalMsgGatewayRegister, nil)               
	registerUnmashalHandles("gateway/MsgGatewayIndexNum", unmashalMsgGatewayIndexNum, nil)               
	registerUnmashalHandles("gateway/MsgGatewayUndelegate", unmashalMsgGatewayUndelegate, nil)           
	registerUnmashalHandles("gateway/MsgGatewayBeginRedelegate", unmashalMsgGatewayBeginRedelegate, nil) 
	registerUnmashalHandles("gateway/MsgGatewayUpload", unmashalMsgGatewayUpload, nil)                   //key
	registerUnmashalHandles("gateway/MsgGatewayEdit", unmashalMsgGatewayEdit, nil)                       

	registerUnmashalHandles("contract/ChatTokenIssue", unmashalChatTokenIssue, nil)
	registerUnmashalHandles("contract/AppTokenIssue", unmashalMsgAppTokenIssue, nil)
	registerUnmashalHandles("contract/CrossChainOut", unmashalCrossChainOut, nil)
	registerUnmashalHandles("proposal_upgrade", unmashalProposalUpgrade, nil)
	registerUnmashalHandles("proposal_community", unmashalProposalCommunity, nil)
	registerUnmashalHandles("proposal_params", unmashalProposalParams, nil)
	registerUnmashalHandles("cosmos-sdk/group/MsgSubmitProposal", unmashalMsgSubmitProposal, nil)
	registerUnmashalHandles("cosmos-sdk/group/MsgVote", unmashalGroupMsgVote, nil)
	registerUnmashalHandles("cosmos-sdk/group/MsgExec", unmashalGroupMsgExec, nil)
}
