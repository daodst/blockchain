package client

import (
	"encoding/json"
	"errors"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	chatTypes "freemasonry.cc/blockchain/x/chat/types"
	contractTypes "freemasonry.cc/blockchain/x/contract/types"
	daotypes "freemasonry.cc/blockchain/x/dao/types"
	"freemasonry.cc/blockchain/x/gateway/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakeTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func unmashalMsgSend(msgByte []byte) (sdk.Msg, error) {
	msg := bankTypes.MsgSend{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgDelegate(msgByte []byte) (sdk.Msg, error) {
	msg := stakeTypes.MsgDelegate{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDistMsgWithdrawDelegatorReward(msgByte []byte) (sdk.Msg, error) {
	msg := distributionTypes.MsgWithdrawDelegatorReward{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgVote(msgByte []byte) (sdk.Msg, error) {
	msg := govTypes.MsgVote{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgMobileTransfer(msgByte []byte) (sdk.Msg, error) {
	msg := chatTypes.MsgMobileTransfer{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgBurnGetMobile(msgByte []byte) (sdk.Msg, error) {
	msg := chatTypes.MsgBurnGetMobile{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgSetChatInfo(msgByte []byte) (sdk.Msg, error) {
	msg := SetChatInfo{}
	err := util.Json.Unmarshal(msgByte, &msg)
	if err != nil {
		return nil, err
	}

	realMas := chatTypes.MsgSetChatInfo{
		FromAddress:      msg.FromAddress,
		GatewayAddress:   msg.NodeAddress,
		AddressBook:      msg.AddressBook,
		ChatBlacklist:    msg.ChatBlacklist,
		ChatWhitelist:    msg.ChatWhitelist,
		UpdateTime:       msg.UpdateTime,
		ChatBlacklistEnc: msg.ChatBlacklistEnc,
		ChatWhitelistEnc: msg.ChatWhitelistEnc,
	}

	return &realMas, err
}

func unmashalChatTokenIssue(msgByte []byte) (sdk.Msg, error) {
	msg := contractTypes.MsgChatTokenIssue{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgAppTokenIssue(msgByte []byte) (sdk.Msg, error) {
	msg := contractTypes.MsgAppTokenIssue{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalCrossChainOut(msgByte []byte) (sdk.Msg, error) {
	msg := CrossChainOut{}
	err := util.Json.Unmarshal(msgByte, &msg)
	if err != nil {
		return nil, err
	}

	
	amountStr, err := core.DealAmount(msg.CoinAmount, msg.CoinSymbol)
	if err != nil {
		return nil, err
	}

	realMas := contractTypes.MsgCrossChainOut{
		SendAddress: msg.SendAddress,
		ToAddress:   msg.ToAddress,
		Coins:       amountStr,
		ChainType:   msg.ChainType,
		Remark:      msg.Remark,
	}

	return &realMas, err
}

func unmashalProposalParams(msgByte []byte) (sdk.Msg, error) {
	proposals := struct {
		Proposer    string `json:"proposer"`
		Deposit     string `json:"deposit"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Change      string `json:"change"`
	}{}
	err := util.Json.Unmarshal(msgByte, &proposals)
	if err != nil {
		return nil, err
	}
	var changes []proposal.ParamChange
	err = json.Unmarshal([]byte(proposals.Change), &changes)
	if err != nil {
		return nil, err
	}
	content := paramproposal.NewParameterChangeProposal(proposals.Title, proposals.Description, changes)
	deposit, err := sdk.ParseCoinsNormalized(proposals.Deposit)
	if err != nil {
		return nil, err
	}
	from, err := sdk.AccAddressFromBech32(proposals.Proposer)
	if err != nil {
		return nil, err
	}
	msg, err := govTypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return nil, err
	}
	err = msg.ValidateBasic()
	return msg, err
}

func unmashalProposalCommunity(msgByte []byte) (sdk.Msg, error) {
	proposal := struct {
		Proposer    string `json:"proposer"`
		Deposit     string `json:"deposit"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Amount      string `json:"amount"`
		Recipient   string `json:"recipient"`
	}{}
	err := util.Json.Unmarshal(msgByte, &proposal)
	if err != nil {
		return nil, err
	}
	recpAddr, err := sdk.AccAddressFromBech32(proposal.Recipient)
	if err != nil {
		return nil, err
	}
	am, err := sdk.ParseCoinsNormalized(proposal.Amount)
	if err != nil {
		return nil, err
	}
	content := distributionTypes.NewCommunityPoolSpendProposal(proposal.Title, proposal.Description, recpAddr, am)
	deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
	if err != nil {
		return nil, err
	}
	from, err := sdk.AccAddressFromBech32(proposal.Proposer)
	if err != nil {
		return nil, err
	}
	msg, err := govTypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return nil, err
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func unmashalProposalUpgrade(msgByte []byte) (sdk.Msg, error) {
	proposal := struct {
		Proposer    string `json:"proposer"`
		Deposit     string `json:"deposit"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Info        string `json:"info"`
		Height      int64  `json:"height"`
		Name        string `json:"name"`
	}{}
	err := util.Json.Unmarshal(msgByte, &proposal)
	if err != nil {
		return nil, err
	}
	plan := upgradeTypes.Plan{
		Name:   proposal.Name,
		Height: proposal.Height,
		Info:   proposal.Info,
	}

	err = UpgradeJsonValidateBasic(plan)
	if err != nil {
		return nil, err
	}
	content := upgradeTypes.NewSoftwareUpgradeProposal(proposal.Title, proposal.Description, plan)
	deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
	if err != nil {
		return nil, err
	}
	from, err := sdk.AccAddressFromBech32(proposal.Proposer)
	if err != nil {
		return nil, err
	}
	msg, err := govTypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return nil, err
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func UpgradeJsonValidateBasic(plan upgradeTypes.Plan) error {
	info, err := plan.UpgradeInfo()
	if err != nil {
		return err
	}
	if info.Gateway == nil && info.App == nil && info.Blockchain == nil {
		return errors.New("The json content is illegal")
	}
	return nil
}

/********** pledge **********/

//func unmashalPledgeMsgWithdrawValidatorCommission(msgByte []byte) (sdk.Msg, error) {
//	msg := pledgeTypes.MsgWithdrawValidatorCommission{}
//	err := util.Json.Unmarshal(msgByte, &msg)
//	return &msg, err
//}

//func unmashalMsgChangeClusterInfo(msgByte []byte) (sdk.Msg, error) {
//	msg := pledgeTypes.MsgChangeClusterInfo{}
//	err := util.Json.Unmarshal(msgByte, &msg)
//	return &msg, err
//}

//func unmashalMsgCreateDeviceCluster(msgByte []byte) (sdk.Msg, error) {
//	msg := pledgeTypes.MsgCreateDeviceCluster{}
//	err := util.Json.Unmarshal(msgByte, &msg)
//	return &msg, err
//}

//func unmashalPledgeMsgUnpledge(msgByte []byte) (sdk.Msg, error) {
//	msg := pledgeTypes.MsgUnpledge{}
//	err := util.Json.Unmarshal(msgByte, &msg)
//	return &msg, err

//}

//func unmashalPledgeMsgPledge(msgByte []byte) (sdk.Msg, error) {
//	msg := pledgeTypes.MsgPledge{}
//	err := util.Json.Unmarshal(msgByte, &msg)
//	return &msg, err
//}


//func unmashalPledgeMsgWithdrawDelegatorReward(msgByte []byte)(sdk.Msg,error){
//	msg := pledgeTypes.MsgWithdrawDelegatorReward{}
//	err := util.Json.Unmarshal(msgByte, &msg)
//	return &msg, err
//}

/********** gateway **********/

func unmashalMsgGatewayUndelegate(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgGatewayUndelegate{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgGatewayEdit(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgGatewayEdit{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgGatewayBeginRedelegate(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgGatewayBeginRedelegate{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgGatewayUpload(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgGatewayUpload{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgGatewayIndexNum(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgGatewayIndexNum{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgGatewayRegister(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgGatewayRegister{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgCreateSmartValidator(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgCreateSmartValidator{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

/*************** dao *************/

func unmashalDaoMsgCreateCluster(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgCreateCluster{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgDeleteMembers(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgDeleteMembers{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgClusterAddMembers(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterAddMembers{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgClusterChangeName(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterChangeName{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgClusterMemberExit(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterMemberExit{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgBurnToPower(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgBurnToPower{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgClusterChangeSalaryRatio(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterChangeSalaryRatio{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgClusterChangeDeviceRatio(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterChangeDeviceRatio{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgClusterChangeId(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterChangeId{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgWithdrawDeviceReward(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgWithdrawDeviceReward{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgWithdrawOwnerReward(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgWithdrawOwnerReward{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgWithdrawBurnReward(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgWithdrawBurnReward{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgThawFrozenPower(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgThawFrozenPower{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgColonyRate(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgColonyRate{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}
func unmashalMsgSubmitProposal(msgByte []byte) (sdk.Msg, error) {
	msg := group.MsgSubmitProposal{}
	err := msg.Unmarshal(msgByte)
	return &msg, err
}
func unmashalGroupMsgVote(msgByte []byte) (sdk.Msg, error) {
	msg := group.MsgVote{}
	err := msg.Unmarshal(msgByte)
	return &msg, err
}
func unmashalGroupMsgExec(msgByte []byte) (sdk.Msg, error) {
	msg := group.MsgExec{}
	err := msg.Unmarshal(msgByte)
	return &msg, err
}
