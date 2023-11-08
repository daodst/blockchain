package client

import (
	"encoding/hex"
	"errors"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/sirupsen/logrus"
	ttypes "github.com/tendermint/tendermint/types"
)

type ClusterClient struct {
	TxClient  *TxClient
	ServerUrl string
	logPrefix string
}

// QueryClusterInfo 
func (clusterClient ClusterClient) QueryClusterInfo(clusterId string) (types.ClusterInfo, error) {
	logs := core.BuildLog(core.GetStructFuncName(clusterClient), core.LmChainClient).WithFields(logrus.Fields{"clusterId": clusterId})

	resp := types.ClusterInfo{}

	param := []byte(clusterId)

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryClusterInfo, param)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return resp, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return resp, err
	}
	return resp, nil
}

// QueryInClusters 
func (clusterClient ClusterClient) QueryInClusters(fromAddress string) ([]types.InClusters, error) {
	logs := core.BuildLog(core.GetStructFuncName(clusterClient), core.LmChainClient).WithFields(logrus.Fields{"fromAddress": fromAddress})

	resp := make([]types.InClusters, 0)

	_, err := sdk.AccAddressFromBech32(fromAddress)
	if err != nil {
		logs.WithError(err).Error("AccAddressFromBech32")
		return resp, core.ErrAddressFormat
	}

	params := []byte(fromAddress)

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryInClusters, params)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return resp, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return resp, err
	}
	return resp, nil
}

// gas
func (clusterClient ClusterClient) QueryClusterGasReward(clusterId, member string) (sdk.DecCoins, error) {
	log := core.BuildLog(core.GetStructFuncName(clusterClient), core.LmChainClient).WithFields(logrus.Fields{"clusterId": clusterId})
	params := types.QueryClusterRewardParams{ClusterId: clusterId, Member: member}
	bz, err := util.Json.Marshal(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return sdk.DecCoins{}, err
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryClusterGasReward, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return sdk.DecCoins{}, nil
	}
	res := sdk.DecCoins{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return sdk.DecCoins{}, nil
	}
	return res, nil
}


func (clusterClient ClusterClient) QueryClusterDeviceReward(clusterId, member string) (sdk.DecCoins, error) {
	log := core.BuildLog(core.GetStructFuncName(clusterClient), core.LmChainClient).WithFields(logrus.Fields{"clusterId": clusterId})
	params := types.QueryClusterRewardParams{ClusterId: clusterId, Member: member}
	bz, err := util.Json.Marshal(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return sdk.DecCoins{}, err
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryClusterDeviceReward, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return sdk.DecCoins{}, nil
	}
	res := sdk.DecCoins{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return sdk.DecCoins{}, nil
	}
	return res, nil
}

func (clusterClient ClusterClient) QueryPersonClusterInfo(from string) (types.PersonalClusterInfo, error) {
	log := core.BuildLog(core.GetStructFuncName(clusterClient), core.LmChainClient).WithFields(logrus.Fields{"address": from})

	res := types.PersonalClusterInfo{}
	params := types.QueryPersonClusterInfoRequest{
		From: from,
	}
	bz, err := util.Json.Marshal(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return types.PersonalClusterInfo{}, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryPersonClusterInfo, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return types.PersonalClusterInfo{}, nil
	}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return types.PersonalClusterInfo{}, nil
	}

	return res, nil
}

// QueryClusterInfoById id
func (clusterClient ClusterClient) QueryClusterInfoById(chatClusterId string) (types.DeviceCluster, error) {
	log := core.BuildLog(core.GetStructFuncName(clusterClient), core.LmChainClient).WithFields(logrus.Fields{"chatClusterId": chatClusterId})

	bz := []byte(chatClusterId)
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryClusterInfoById, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return types.DeviceCluster{}, err
	}

	res := types.DeviceCluster{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
	if err != nil {
		log.WithError(err).Error("QueryClusterInfoById UnmarshalJSON Error")
		return types.DeviceCluster{}, err
	}

	return res, nil
}

func (this *ClusterClient) CreateCluster(from string, fee legacytx.StdFee, gatewayAddress, clusterId, clusterName, chatAddress string, deviceRatio, salaryRatio, burnAmount sdk.Dec, privateKey string) (tx ttypes.Tx, resp *core.BroadcastTxResponse, err error) {

	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	msg := types.NewMsgCreateCluster(from, gatewayAddress, clusterId, chatAddress, clusterName, deviceRatio, salaryRatio, burnAmount, sdk.ZeroDec())
	if err != nil {
		log.Error("NewMsgTransfer")
		return
	}
	var result *core.BaseResponse
	
	tx, result, err = this.TxClient.SignAndSendMsg(from, privateKey, fee, "txBase.Memo", msg)
	if err != nil {
		return
	}
	resp = new(core.BroadcastTxResponse)
	
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return tx, resp, nil
	} else {
		
		resp.TxHash = hex.EncodeToString(tx.Hash())
		return tx, resp, errors.New(result.Info)
	}
}

func (clusterClient ClusterClient) QueryDaoParams() (types.DaoParams, error) {
	logs := core.BuildLog(core.GetStructFuncName(clusterClient), core.LmChainClient)

	resp := types.DaoParams{}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryDaoParams, nil)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return resp, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return resp, err
	}
	return resp, nil
}
