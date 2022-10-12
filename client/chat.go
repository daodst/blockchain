package client

import (
    "freemasonry.cc/blockchain/cmd/config"
    "freemasonry.cc/blockchain/util"
    "freemasonry.cc/blockchain/x/chat/types"
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

//
//func (this *ChatClient) Register(registerData types.RegisterData, privateKey string) (resp *core.BroadcastTxResponse, mobile int64, err error) {
//	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
//	//
//	msg := types.NewMsgRegister(registerData.FromAddress, registerData.NodeAddress, registerData.MobilePrefix, registerData.MortgageAmount)
//	if err != nil {
//		log.WithError(err).Error("NewMsgRegister")
//		return
//	}
//
//	//
//	resp, err = this.TxClient.SignAndSendMsg(registerData.FromAddress, privateKey, registerData.Fee, "", msg)
//	if err != nil {
//		return
//	}
//
//	if resp.Status == 0 {
//		return resp, mobile, errors.New(resp.Info)
//	}
//
//	//
//	params := types.QueryUserInfoParams{Address: registerData.FromAddress}
//	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
//	if err != nil {
//		log.WithError(err).Error("MarshalJSON")
//		return nil, mobile, err
//	}
//	userInfoByte, _, err := clientCtx.QueryWithData("custom/chat/"+types.QueryUserInfo, bz)
//
//	userInfo := &types.UserInfo{}
//	if userInfoByte != nil {
//		err = util.Json.Unmarshal(userInfoByte, userInfo)
//		if err != nil {
//			log.WithError(err).Error("Unmarshal")
//			return nil, mobile, err
//		}
//	}
//
//	if len(userInfo.Mobile) == 0 {
//		log.WithError(err).Error("Mobile Length")
//		return nil, mobile, errors.New("mobile Not Found")
//	}
//
//	mobileInt64, err := strconv.ParseInt(userInfo.Mobile[0], 10, 64)
//	if err == nil {
//		return nil, mobile, errors.New("mobile error")
//	}
//
//	return resp, mobileInt64, nil
//}

type GetUserInfo struct {
    Status   int            `json:"status"`    //
    Message  string         `json:"message"`   //
    UserInfo types.UserInfo `json:"user_info"` //
}

//
func (this *ChatClient) QueryUserInfo(address string) (data *GetUserInfo, err error) {
    log := util.BuildLog(util.GetStructFuncName(this), util.LmChainClient).WithFields(logrus.Fields{"address": address})
    params := types.QueryUserInfoParams{Address: address}
    bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
    if err != nil {
        log.WithError(err).Error("MarshalJSON")
        return nil, err
    }
    resBytes, _, err := clientCtx.QueryWithData("custom/chat/"+types.QueryUserInfo, bz)
    if err != nil {
        log.WithError(err).Error("QueryWithData")
        return nil, err
    }
    userInfo := &types.UserInfo{}
    if resBytes != nil {
        err = util.Json.Unmarshal(resBytes, userInfo)
        if err != nil {
            log.WithError(err).Error("Unmarshal")
            return nil, err
        }
    }

    data = &GetUserInfo{}
    if userInfo.FromAddress == "" {
        data.Status = 0
        return
    }

    data.Status = 1
    data.UserInfo = *userInfo

    //data.UserInfo.Mobile = userInfo.Mobile
    //data.UserInfo.NodeAddress = userInfo.NodeAddress
    //data.UserInfo.FromAddress = userInfo.FromAddress
    //data.UserInfo.MortgageAmount = userInfo.MortgageAmount
    //data.UserInfo.CanRedemAmount = userInfo.CanRedemAmount
    //data.UserInfo.ChatFee = userInfo.ChatFee

    return
}

//
func (this *ChatClient) QueryTotalPledge(address string) (data sdk.Coin, err error) {
    log := util.BuildLog(util.GetStructFuncName(this), util.LmChainClient).WithFields(logrus.Fields{"address": address})

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

    allPledge := pledgeTypese.DelegationResponses{}

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

func (this *ChatClient) QueryPledgeParams() (data pledgeTypese.Params, err error) {
    log := util.BuildLog(util.GetStructFuncName(this), util.LmChainClient)

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
    log := util.BuildLog(util.GetStructFuncName(this), util.LmChainClient)

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
    log := util.BuildLog(util.GetStructFuncName(this), util.LmChainClient)

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
