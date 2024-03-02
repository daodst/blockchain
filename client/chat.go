package client

import (
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/chat/types"
	daotypes "freemasonry.cc/blockchain/x/dao/types"
	gatewayTypes "freemasonry.cc/blockchain/x/gateway/types"
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
	types.UserInfo
	PledgeLevel         int64  `json:"pledge_level"`          
	GatewayProfixMobile string `json:"gateway_profix_mobile"` 
}

type QueryUserListInfo []types.AllUserInfo


func (this *ChatClient) QueryUserByMobile(mobile string) (data types.AllUserInfo, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"mobile": mobile})

	params := types.QueryUserByMobileParams{Mobile: mobile}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return types.AllUserInfo{}, err
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/chat/"+types.QueryUserByMobile, bz)
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

	log.Info("QueryUserInfos+++++++++++++++++++")
	log.Info("addresses:")
	for _, v := range addresses {
		log.Info(v)
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/chat/"+types.QueryUserInfos, bz)
	if err != nil {
		log.Info("error:", err.Error())
		log.WithError(err).Error("QueryWithData1")
		return nil, err
	}
	log.Info("QueryUserInfos---------------------")
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


func (this *ChatClient) QueryUsersChatInfo(addresses []string) (data []types.CustomInfo, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	params := types.QueryUserInfosParams{Addresses: addresses}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/chat/"+types.QueryUsersChatInfo, bz)
	if err != nil {
		log.Info("error:", err.Error())
		log.WithError(err).Error("QueryWithData1")
		return nil, err
	}
	log.Info("QueryUserInfos---------------------")
	userInfos := []types.CustomInfo{}
	if resBytes != nil {
		err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &userInfos)
		if err != nil {
			log.WithError(err).Error("Unmarshal1")
			return nil, err
		}
	}

	return userInfos, nil
}


func (this *ChatClient) QueryUserInfo(address string) (data *GetUserInfo, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"address": address})
	params := types.QueryUserInfoParams{Address: address}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/chat/"+types.QueryUserInfo, bz)
	if err != nil {
		if err.Error() == "user not found" {
			
			chatParamsResBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/chat/"+types.QueryParams, nil)
			if err != nil {
				return nil, err
			}
			var chatParams types.Params
			err = clientCtx.LegacyAmino.UnmarshalJSON(chatParamsResBytes, &chatParams)
			if err != nil {
				log.WithError(err).Error("Unmarshal1")
				return nil, err
			}

			//data.UserInfo.UserInfo.FromAddress = address
			data = &GetUserInfo{
				Status:   1,
				UserInfo: QueryUserInfo{},
			}
			return data, nil
		}

		log.Info("error :" + err.Error())
		log.WithError(err).Error("QueryWithData1")
		return nil, err
	}
	userInfo := &types.UserInfo{}

	if resBytes != nil {

		err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, userInfo)
		if err != nil {
			log.WithError(err).Error("Unmarshal1")
			return nil, err
		}
	}

	
	burnLevelBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+daotypes.QueryBurnLevels, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData2")
		return nil, err
	}

	pledgeLevelInfo := make(map[string]int64)
	burnLevel := int64(0)

	if burnLevelBytes != nil {
		err = clientCtx.LegacyAmino.UnmarshalJSON(burnLevelBytes, &pledgeLevelInfo)
		if err != nil {
			log.WithError(err).Error("Unmarshal2")
			return nil, err
		}
	}

	if _, ok := pledgeLevelInfo[address]; ok {
		burnLevel = pledgeLevelInfo[address]
	}

	gatewayInfo := new(gatewayTypes.Gateway)
	if userInfo.NodeAddress != "" {
		gparams := gatewayTypes.QueryGatewayInfoParams{GatewayAddress: userInfo.NodeAddress}
		bz, err = clientCtx.LegacyAmino.MarshalJSON(gparams)
		if err != nil {
			log.WithError(err).Error("MarshalJSON")
			return nil, err
		}
		gatewayFirstMobileByte, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/gateway/"+gatewayTypes.QueryGatewayInfo, bz)
		if err != nil {
			log.WithError(err).Error("QueryWithData3")
			return nil, err
		}

		err = clientCtx.LegacyAmino.UnmarshalJSON(gatewayFirstMobileByte, &gatewayInfo)
		if err != nil {
			log.WithError(err).Error("Unmarshal3")
			return nil, err
		}
	}

	data = &GetUserInfo{}

	data.Status = 1
	data.UserInfo = QueryUserInfo{
		UserInfo:            *userInfo,
		PledgeLevel:         burnLevel,
		GatewayProfixMobile: gatewayInfo.GatewayNum[0].NumberIndex,
	}

	return
}


func (this *ChatClient) QueryChatGain(address string) (sdk.Int, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"address": address})
	params := types.QueryUserInfoParams{Address: address}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return sdk.ZeroInt(), err
	}

	
	pledgeLevelBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+daotypes.QueryBurnLevel, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData2")
		return sdk.ZeroInt(), err
	}

	burnLevelInfo := daotypes.BurnLevel{}

	if pledgeLevelBytes != nil {
		err = clientCtx.LegacyAmino.UnmarshalJSON(pledgeLevelBytes, &burnLevelInfo)
		if err != nil {
			log.WithError(err).Error("Unmarshal2")
			return sdk.ZeroInt(), err
		}
	}

	return burnLevelInfo.AddPercent, nil
}

type DelegateInfo struct {
	Name    string `json:"name"`    
	Gateway string `json:"gateway"` 
	Amount  string `json:"amount"`  
}

type DelegateInfos []DelegateInfo


func (this *ChatClient) QueryAddrByChatAddr(chatAddr string) (string, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"address": chatAddr})
	params := chatAddr
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return "", err
	}

	chatAddrBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/chat/"+types.QueryAddrByChatAddr, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData2")
		return "", err
	}

	fromAddr := ""

	if chatAddrBytes != nil {
		err = clientCtx.LegacyAmino.UnmarshalJSON(chatAddrBytes, &fromAddr)
		if err != nil {
			log.WithError(err).Error("Unmarshal2")
			return "", err
		}
	}

	return fromAddr, nil
}

// QueryBurnLevels 
func (this *ChatClient) QueryBurnLevels(addresses []string) (data map[string]int64, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"addresses": addresses})
	res := make(map[string]int64)
	params := types.QueryPledgeLevelsParams{Addresses: addresses}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return res, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+daotypes.QueryBurnLevels, bz)
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
