package client

import (
	"freemasonry.cc/blockchain/cmd/config"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/chat/types"
	"freemasonry.cc/blockchain/x/pledge/keeper"
	pledgeTypese "freemasonry.cc/blockchain/x/pledge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sirupsen/logrus"
)

type ChatInfo struct {
	TxClient      *TxClient
	AccountClient *AccountClient
	ServerUrl     string
	logPrefix     string
}

type ChatClient struct {
	TxClient      *TxClient
	AccountClient *AccountClient
	ServerUrl     string
	logPrefix     string
}

type GetUserInfo struct {
	Status   int           `json:"status"`    
	Message  string        `json:"message"`   
	UserInfo QueryUserInfo `json:"user_info"` 
}

type QueryUserInfo struct {
	types.UserInfoToApp
	PledgeLevel         int64  `json:"pledge_level"`          
	GatewayProfixMobile string `json:"gateway_profix_mobile"` 
}

type QueryUserListInfo []types.AllUserInfo


func (this *ChatClient) QueryPledgeLevels(addresses []string) (data map[string]int64, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"addresses": addresses})
	res := make(map[string]int64)
	params := types.QueryPledgeLevelsParams{Addresses: addresses}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return res, err
	}

	resBytes, _, err := clientCtx.QueryWithData("custom/chat/"+types.QueryPledgeLevels, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return res, err
	}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return res, err
	}

	return res, nil
}


func (this *ChatClient) QueryUserByMobile(mobile string) (data types.AllUserInfo, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"mobile": mobile})

	params := types.QueryUserByMobileParams{Mobile: mobile}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return types.AllUserInfo{}, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/chat/"+types.QueryUserByMobile, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return types.AllUserInfo{}, err
	}

	userInfo := &types.AllUserInfo{}
	if resBytes != nil {
		err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, userInfo)
		if err != nil {
			log.WithError(err).Error("UnmarshalJSON")
			return types.AllUserInfo{}, err
		}
	}

	return *userInfo, nil
}


func (this *ChatClient) QueryUserInfos(addresses []string) (data QueryUserListInfo, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	params := types.QueryUserInfosParams{Addresses: addresses}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/chat/"+types.QueryUserInfos, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData1")
		return nil, err
	}

	userInfos := &QueryUserListInfo{}
	if resBytes != nil {
		err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, userInfos)
		if err != nil {
			log.WithError(err).Error("Unmarshal1")
			return nil, err
		}
	}

	return *userInfos, nil
}


func (this *ChatClient) QueryUserInfo(address string) (data *GetUserInfo, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"address": address})
	params := types.QueryUserInfoParams{Address: address}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/chat/"+types.QueryUserInfo, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData1")
		return nil, err
	}
	userInfo := &types.UserInfoToApp{}

	if resBytes != nil {

		err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, userInfo)
		if err != nil {
			log.WithError(err).Error("Unmarshal1")
			return nil, err
		}
	}

	
	pledgeLevelBytes, _, err := clientCtx.QueryWithData("custom/pledge/"+pledgeTypese.QueryPledgeLevel, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData2")
		return nil, err
	}
	pledgeLevel := int64(0)

	if pledgeLevelBytes != nil {
		err = clientCtx.LegacyAmino.UnmarshalJSON(pledgeLevelBytes, &pledgeLevel)
		if err != nil {
			log.WithError(err).Error("Unmarshal2")
			return nil, err
		}
	}

	params = types.QueryUserInfoParams{Address: userInfo.NodeAddress}
	bz, err = clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	gatewayFirstMobileByte, _, err := clientCtx.QueryWithData("custom/pledge/"+pledgeTypese.QueryGatewayFirstMobile, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData3")
		return nil, err
	}

	var gatewayFirstMobile string
	err = clientCtx.LegacyAmino.UnmarshalJSON(gatewayFirstMobileByte, &gatewayFirstMobile)
	if err != nil {
		log.WithError(err).Error("Unmarshal3")
		return nil, err
	}

	data = &GetUserInfo{}
	if userInfo.FromAddress == "" {
		data.Status = 0
		return
	}

	data.Status = 1
	data.UserInfo = QueryUserInfo{
		UserInfoToApp:       *userInfo,
		PledgeLevel:         pledgeLevel,
		GatewayProfixMobile: gatewayFirstMobile,
	}

	return
}


func (this *ChatClient) QueryTotalPledge(address string) (data sdk.Coin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"address": address})

	accAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return sdk.Coin{}, err
	}

	params := pledgeTypese.NewQueryDelegatorParams(accAddress)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return sdk.Coin{}, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/pledge/"+pledgeTypese.QueryDelegatorDelegations, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return sdk.Coin{}, err
	}

	allPledge := make([]keeper.QueryDelegationsResp, 0)

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &allPledge)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return sdk.Coin{}, err
	}

	if len(allPledge) == 0 {
		return sdk.NewCoin(config.BaseDenom, sdk.NewInt(0)), nil
	}

	allPledgeAmountInt := sdk.ZeroInt()
	for _, banl := range allPledge {
		allPledgeAmountInt = allPledgeAmountInt.Add(banl.Balance.Amount)
	}
	allPledgeAmountCoin := sdk.NewCoin(config.BaseDenom, allPledgeAmountInt)
	return allPledgeAmountCoin, nil
}

type DelegateInfo struct {
	Name    string `json:"name"`    
	Gateway string `json:"gateway"` 
	Amount  string `json:"amount"`  
}

type DelegateInfos []DelegateInfo


func (this *ChatClient) QueryAllPledges(address string) (DelegateInfos, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"address": address})

	delegateInfos := make([]DelegateInfo, 0)

	accAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		log.WithError(err).Error("invalid address")
		return delegateInfos, err
	}

	params := pledgeTypese.NewQueryDelegatorParams(accAddress)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return delegateInfos, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/pledge/"+pledgeTypese.QueryDelegatorDelegations, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return delegateInfos, err
	}

	allPledge := make([]keeper.QueryDelegationsResp, 0)

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &allPledge)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return delegateInfos, err
	}

	if len(allPledge) == 0 {
		return delegateInfos, nil
	}

	for _, banl := range allPledge {
		delegateInfos = append(delegateInfos, DelegateInfo{
			Gateway: banl.Delegation.ValidatorAddress,
			Amount:  banl.Balance.Amount.String(),
		})
	}
	return delegateInfos, nil
}

func (this *ChatClient) QueryPledgeParams() (data pledgeTypese.Params, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	resBytes, _, err := clientCtx.QueryWithData("custom/pledge/"+pledgeTypese.QueryParams, nil)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return pledgeTypese.Params{}, err
	}

	params := pledgeTypese.Params{}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return pledgeTypese.Params{}, err
	}

	return params, nil
}

func (this *ChatClient) QueryChatParams() (data types.Params, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	resBytes, _, err := clientCtx.QueryWithData("custom/chat/"+types.QueryParams, nil)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return types.Params{}, err
	}

	params := types.Params{}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return types.Params{}, err
	}

	return params, nil
}

func (this *ChatClient) QueryPrePledge(address string) (data sdk.Dec, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	param, err := clientCtx.LegacyAmino.MarshalJSON(address)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return sdk.Dec{}, err
	}

	resBytes, _, err := clientCtx.QueryWithData("custom/pledge/"+pledgeTypese.QueryPrePledge, param)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return sdk.ZeroDec(), err
	}

	amount := sdk.ZeroDec()

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &amount)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return sdk.ZeroDec(), err
	}

	return amount, nil
}

func (this *ChatClient) QueryAllCanWithdraw(fromaddress string) (data sdk.Coin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	accFromaddress, err := sdk.AccAddressFromBech32(fromaddress)
	if err != nil {
		log.WithError(err).Error("invalid address")
		return data, err
	}

	param, err := clientCtx.LegacyAmino.MarshalJSON(accFromaddress)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return data, err
	}

	resBytes, _, err := clientCtx.QueryWithData("custom/pledge/"+pledgeTypese.QueryAllCanWithdraw, param)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return data, err
	}

	res := sdk.Coin{}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return data, err
	}

	return res, nil
}


func (this *ChatClient) QueryPledgeInfo(address string) (allPledgeInfo keeper.QueryPledgeInfoResp, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"address": address})

	accAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return allPledgeInfo, err
	}

	bz, err := clientCtx.LegacyAmino.MarshalJSON(accAddress)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return allPledgeInfo, err
	}

	resBytes, _, err := clientCtx.QueryWithData("custom/pledge/"+pledgeTypese.QueryPledgeInfo, bz)

	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return allPledgeInfo, err
	}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &allPledgeInfo)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return allPledgeInfo, err
	}

	return allPledgeInfo, nil
}


func (this *ChatClient) QueryChatSendGift(fromAddress, toAddress string) (bool, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"fromAddress": fromAddress, "toAddress": toAddress})

	param := types.QueryChatSendGiftInfoParams{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
	}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(param)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return false, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/chat/"+types.QueryChatSeneGift, bz)
	if err != nil {
		log.WithError(err).Error("QueryChatSeneGift")
		return false, err
	}

	var isPay bool
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &isPay)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return false, err
	}
	return isPay, nil
}


func (this *ChatClient) QueryChatSendGifts(fromAddress string, toAddresses []string) (map[string]bool, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"fromAddress": fromAddress, "toAddresses": toAddresses})

	param := types.QueryChatSendGiftsInfoParams{
		FromAddress: fromAddress,
		ToAddresses: toAddresses,
	}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(param)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/chat/"+types.QueryChatSeneGifts, bz)
	if err != nil {
		log.WithError(err).Error("QueryChatSeneGift")
		return nil, err
	}

	var isPays map[string]bool
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &isPays)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, err
	}
	return isPays, nil
}
