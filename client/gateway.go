package client

import (
	"context"
	"encoding/hex"
	"errors"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	daoTypes "freemasonry.cc/blockchain/x/dao/types"
	daotypes "freemasonry.cc/blockchain/x/dao/types"
	"freemasonry.cc/blockchain/x/gateway/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/sirupsen/logrus"
	ttypes "github.com/tendermint/tendermint/types"
	"strconv"
	"strings"
	"time"
)

type GatewayClient struct {
	TxClient      *TxClient
	AccountClient *AccountClient
	ServerUrl     string
	logPrefix     string
}

func (this GatewayClient) QueryModuleParameters(module string) (interface{}, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	if module == "gov" {
		resp := struct {
			Height int64                        `json:"height"`
			Result govTypes.QueryParamsResponse `json:"result"`
		}{}
		resBytes1, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+module+"/params/"+govTypes.ParamDeposit, nil)
		if err != nil {
			log.WithError(err).Error("QueryWithDataWithUnwrapErr " + govTypes.ParamDeposit)
			return nil, err
		}
		if resBytes1 != nil {
			param := govTypes.DepositParams{}
			err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes1, &param)
			if err != nil {
				return nil, err
			}
			resp.Result.DepositParams = param
		}
		resBytes2, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+module+"/params/"+govTypes.ParamVoting, nil)
		if err != nil {
			log.WithError(err).Error("QueryWithDataWithUnwrapErr " + govTypes.ParamDeposit)
			return nil, err
		}
		if resBytes2 != nil {
			param := govTypes.VotingParams{}
			err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes2, &param)
			if err != nil {
				return nil, err
			}
			resp.Result.VotingParams = param
		}
		resBytes3, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+module+"/params/"+govTypes.ParamTallying, nil)
		if err != nil {
			log.WithError(err).Error("QueryWithDataWithUnwrapErr " + govTypes.ParamDeposit)
			return nil, err
		}
		if resBytes3 != nil {
			param := govTypes.TallyParams{}
			err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes3, &param)
			if err != nil {
				return nil, err
			}
			resp.Result.TallyParams = param
		}
		respByte, err := clientCtx.LegacyAmino.MarshalJSON(resp)
		if err != nil {
			return nil, err
		}
		return string(respByte), nil
	} else {
		reponseStr, err := GetRequest(this.ServerUrl, "/"+module+"/parameters")
		if err != nil {
			log.WithError(err).Error("QueryWithData")
			return "", err
		}
		reponseStr = strings.ReplaceAll(reponseStr, "\n", "")
		return reponseStr, nil
	}
}


func (this GatewayClient) QueryDelegateLastHeight(delegateAddress, validatorAddress string) (height int64, param types.Params, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := types.QueryDelegateLastHeight{DelegateAddress: delegateAddress, ValidatorAddress: validatorAddress}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+types.ModuleName+"/"+types.QueryDelegateLastHeightKey, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	if resBytes != nil {
		err = util.Json.Unmarshal(resBytes, &height)
		if err != nil {
			return
		}
	}
	paramsBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+types.ModuleName+"/"+types.QueryParams, nil)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	if paramsBytes != nil {
		param = types.Params{}
		err = clientCtx.LegacyAmino.UnmarshalJSON(paramsBytes, &param)
		if err != nil {
			return
		}
	}
	return
}


func (this GatewayClient) GatewayNumberCount(gatewayAddress, amount string) (count int64, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"gatewayAddress": gatewayAddress})

	coin := sdk.NewCoin(sdk.DefaultBondDenom, core.MustRealString2LedgerIntNoMin(amount))
	params := types.GatewayNumberCountParams{GatewayAddress: gatewayAddress, Amount: coin}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+types.ModuleName+"/"+types.QueryGatewayNumberCount, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	if resBytes != nil {
		err = util.Json.Unmarshal(resBytes, &count)
		if err != nil {
			return
		}
	}
	return
}


func (this *GatewayClient) QueryGateway(gatewayAddress, gatewayNum string) (data *types.Gateway, notFound bool, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"gatewayAddress": gatewayAddress})
	params := types.QueryGatewayInfoParams{GatewayAddress: gatewayAddress, GatewayNumIndex: gatewayNum}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+types.ModuleName+"/"+types.QueryGatewayInfo, bz)
	if err != nil {
		
		if strings.Contains(err.Error(), core.ErrGatewayNotExist.Error()) {
			notFound = true
			err = nil
		} else {
			log.WithError(err).Error("QueryWithData")
		}
		return
	}
	if resBytes != nil {
		data = new(types.Gateway)
		err = util.Json.Unmarshal(resBytes, data)
		if err != nil {
			return
		}
	}
	return
}


func (this *GatewayClient) QueryGatewayList() (data []types.GatewayListResp, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+types.ModuleName+"/"+types.QueryGatewayList, nil)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	if resBytes != nil {
		resp := []types.Gateway{}
		err := util.Json.Unmarshal(resBytes, &resp)
		if err != nil {
			return nil, err
		}
		resBytes, _, err = util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+types.ModuleName+"/"+types.QueryValidators, nil)
		if err != nil {
			log.WithError(err).Error("QueryWithData")
			return nil, err
		}
		validator := []types.ValidatorInfo{}
		err = util.Json.Unmarshal(resBytes, &validator)
		if err != nil {
			log.WithError(err).Error("MarshalJSON")
			return nil, err
		}
		for _, gateway := range resp {
			for _, val := range validator {
				if val.OperatorAddress == gateway.GatewayAddress {
					dataInfo := types.GatewayListResp{}
					dataInfo.Gateway = gateway
					dataInfo.Token = val.Tokens
					if val.Jailed {
						data = append(data, dataInfo)
						continue
					}
					node, err := clientCtx.GetNode()
					if err != nil {
						log.WithError(err).Error("clientCtx.GetNode")
						return nil, err
					}
					params := slashingTypes.QuerySigningInfoRequest{val.ConsAddress.String()}
					bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
					if err != nil {
						log.WithError(err).Error("MarshalJSON")
						return nil, err
					}
					resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/slashing/"+slashingTypes.QuerySigningInfo, bz)
					if err != nil {
						log.WithError(err).Error("QueryWithData")
						return nil, err
					}
					info := slashingTypes.ValidatorSigningInfo{}
					err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &info)
					if err != nil {
						log.WithError(err).Error("MarshalJSON")
						return nil, err
					}

					status, err := node.Status(context.Background())
					if err != nil {
						log.WithError(err).Error("node.Status")
						return nil, err
					}
					startHeight := status.SyncInfo.LatestBlockHeight - info.IndexOffset
					if startHeight == 0 {
						startHeight = 1
					}
					blockInfo, err := node.Block(context.Background(), &startHeight)
					if err != nil {
						log.WithError(err).Error("node.Block")
						return nil, err
					}
					dataInfo.Online = time.Now().Unix() - blockInfo.Block.Time.Unix()
					data = append(data, dataInfo)
				}
			}
		}

	}
	return
}


func (this *GatewayClient) QueryGatewayNumList() (map[string]types.GatewayNumIndex, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+types.ModuleName+"/"+types.QueryGatewayNum, nil)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	data := make(map[string]types.GatewayNumIndex)
	if resBytes != nil {
		err := util.Json.Unmarshal(resBytes, &data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}


func (this *GatewayClient) QueryGatewayRedeemNumList() (map[string]types.GatewayNumIndex, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+types.ModuleName+"/"+types.QueryGatewayRedeemNum, nil)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	data := make(map[string]types.GatewayNumIndex)
	if resBytes != nil {
		err := util.Json.Unmarshal(resBytes, &data)
		if err != nil {
			log.WithError(err).Error("MarshalJSON")
			return nil, err
		}
	}
	return data, nil
}

//   ValidatorStatus:  0 Unbonded , 1 Unbonding , 2 Bonded , 3  , 4 
func (this *GatewayClient) ValidatorInfo() (validatorInfo *types.ValidatorInfor, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	blockClient := NewBlockClient()
	nodeStatus, err := blockClient.StatusInfo()
	if err != nil {
		log.WithError(err).Error("StatusInfo")
		return nil, err
	}

	
	consAddress, err := sdk.ConsAddressFromHex(nodeStatus.ValidatorInfo.Address.String())
	if err != nil {
		log.WithError(err).Error("ConsAddressFromHex")
		return nil, err
	}
	validatorInfo = &types.ValidatorInfor{
		ValidatorStatus:   "", //0  Unbonded 1 Unbonding 2 Bonded 3  4 
		ValidatorPubAddr:  nodeStatus.ValidatorInfo.PubKey.Address().String(),
		ValidatorConsAddr: consAddress.String(),
	}

	
	validator, notFound, err := this.FindValidatorByConsAddress(consAddress.String())
	if notFound {
		validatorInfo.ValidatorStatus = "4" 
		return validatorInfo, nil
	}
	if err != nil {
		log.WithError(err).Error("FindValidatorByConsAddress")
		return nil, err
	}

	validatorInfo.ValidatorStatus = strconv.Itoa(this.GetValidatorStatus(validator.Status, validator.Jailed))
	validatorInfo.ValidatorOperAddr = validator.GetOperator().String()
	accAddre := sdk.AccAddress(validator.GetOperator())
	validatorInfo.AccAddr = accAddre.String()
	return validatorInfo, nil
}

/*
*
 
*/
func (this *GatewayClient) FindValidatorByConsAddress(bech32ConsAddr string) (validator *stakingTypes.Validator, notFound bool, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	notFound = false
	consAddress, err := sdk.ConsAddressFromBech32(bech32ConsAddr)
	if err != nil {
		log.WithError(err).Error("ConsAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return
	}
	params := types.QueryValidatorByConsAddrParams{ValidatorConsAddress: consAddress}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, notFound, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+types.ModuleName+"/"+types.QueryValidatorByConsAddress, bz)
	if err != nil {
		
		if strings.Contains(err.Error(), stakingTypes.ErrNoValidatorFound.Error()) {
			notFound = true
			err = nil 
		} else {
			log.WithError(err).Error("QueryWithData")
		}
		return nil, notFound, err
	}
	validator = &stakingTypes.Validator{}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, validator)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
	}
	return
}


func (this *GatewayClient) GetValidatorStatus(status stakingTypes.BondStatus, jailed bool) int {
	if jailed {
		return 5
	} else {
		switch status.String() {
		case "BOND_STATUS_UNSPECIFIED":
			return 0
		case "BOND_STATUS_UNBONDED":
			return 1 
		case "BOND_STATUS_UNBONDING":
			return 2 
		case "BOND_STATUS_BONDED":
			return 3 
		}
	}
	return 0
}

// key
func (this *GatewayClient) GatewayUpload(address string, data []byte, privateKey string) (resp *core.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	_, err = sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	msg := types.NewMsgGatewayUpload(address, data)
	if err != nil {
		log.Error("NewMsgTransfer")
		return
	}
	
	_, result, err := this.TxClient.SignAndSendMsg(address, privateKey, core.NewLedgerFeeFromGas(flags.DefaultGasLimit, 0), "", msg)
	if err != nil {
		return
	}
	
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		resp = new(core.BroadcastTxResponse)
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return resp, nil
	} else {
		
		return resp, errors.New(result.Info)
	}
}

func (this *GatewayClient) GetGatewayExistUpload(valAddress string) (bool, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	val, err := sdk.ValAddressFromBech32(valAddress)
	if err != nil {
		log.WithError(err).Error("ConsAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return false, err
	}
	params := val.String()
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return false, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+types.ModuleName+"/"+types.QueryGatewayUpload, bz)
	if err != nil {
		return false, err
	}

	if resBytes == nil {
		return false, nil
	}

	return true, nil
}

func (this *GatewayClient) GetGatewayUpload(valAddress string) (*types.GatewayUploadData, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	val, err := sdk.ValAddressFromBech32(valAddress)
	if err != nil {
		log.WithError(err).Error("ConsAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return nil, err
	}
	params := val.String()
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+types.ModuleName+"/"+types.QueryGatewayUpload, bz)
	if err != nil {
		return nil, err
	}

	res := types.GatewayUploadData{}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
	}

	return &res, nil
}


func (this *GatewayClient) GatewayRegister(address, gatewayUrl, delegation, privateKey, packageName, peerId, machineAddress string, indexNum []string) (resp *core.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	_, err = sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	msg := types.NewMsgGatewayRegister(address, gatewayUrl, delegation, packageName, peerId, machineAddress, indexNum)
	if err != nil {
		log.WithError(err).Error("NewMsgGatewayRegister")
		return
	}
	
	tx, result, err := this.TxClient.SignAndSendMsg(address, privateKey, core.NewLedgerFeeFromGas(flags.DefaultGasLimit, 0), "", msg)
	if err != nil {
		return
	}
	
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		resp = new(core.BroadcastTxResponse)
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return resp, nil
	} else {
		
		resp.TxHash = hex.EncodeToString(tx.Hash())
		return resp, errors.New(result.Info)
	}
}


func (this *GatewayClient) GatewayEdit(address, gatewayUrl, privateKey string) (resp *core.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	_, err = sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	msg := types.NewMsgGatewayEdit(address, gatewayUrl)
	if err != nil {
		log.WithError(err).Error("NewMsgGatewayRegister")
		return
	}
	
	tx, result, err := this.TxClient.SignAndSendMsg(address, privateKey, core.NewLedgerFeeFromGas(flags.DefaultGasLimit, 0), "", msg)
	if err != nil {
		return
	}
	
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		resp = new(core.BroadcastTxResponse)
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return resp, nil
	} else {
		
		resp.TxHash = hex.EncodeToString(tx.Hash())
		return resp, errors.New(result.Info)
	}
}


func (this *GatewayClient) GatewayAddIndexNum(address, validatorAddress, privateKey string, indexNum []string) (resp *core.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	_, err = sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	msg := types.NewMsgGatewayIndexNum(address, validatorAddress, indexNum)
	if err != nil {
		log.WithError(err).Error("NewMsgGatewayIndexNum")
		return
	}
	
	tx, result, err := this.TxClient.SignAndSendMsg(address, privateKey, core.NewLedgerFeeFromGas(flags.DefaultGasLimit, 0), "", msg)
	if err != nil {
		return
	}
	
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		resp = new(core.BroadcastTxResponse)
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return resp, nil
	} else {
		
		resp.TxHash = hex.EncodeToString(tx.Hash())
		return resp, errors.New(result.Info)
	}
}


func (this *GatewayClient) GatewayUnDelegate(address, validatorAddress, undelegation, privateKey string, indexNum []string) (resp *core.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	_, err = sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	undelegateCoin, err := sdk.ParseCoinNormalized(undelegation)
	if err != nil {
		return nil, err
	}
	msg := types.NewMsgGatewayUndelegation(address, validatorAddress, undelegateCoin, indexNum)
	if err != nil {
		log.WithError(err).Error("NewMsgGatewayUndelegation")
		return
	}
	
	tx, result, err := this.TxClient.SignAndSendMsg(address, privateKey, core.NewLedgerFeeFromGas(flags.DefaultGasLimit, 0), "", msg)
	if err != nil {
		return
	}
	
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		resp = new(core.BroadcastTxResponse)
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return resp, nil
	} else {
		
		resp.TxHash = hex.EncodeToString(tx.Hash())
		return resp, errors.New(result.Info)
	}
}


func (this *GatewayClient) QueryGatewayClusters(gatewayAddress string) (resp map[string]map[string]daotypes.ClusterDeviceMember, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"gatewayAddress": gatewayAddress})
	params := daoTypes.QueryGatewayClustersParams{
		GatewayAddress: gatewayAddress,
	}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+daoTypes.ModuleName+"/"+daoTypes.QueryGatewayClusters, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData QueryuGatewayClusters")
		return
	}
	if resBytes != nil {
		err = util.Json.Unmarshal(resBytes, &resp)
		if err != nil {
			return
		}

		return resp, nil
	}
	return
}


func (this *GatewayClient) ColonyRate(address, gatewayAddress, privateKey string, rate []daotypes.ColonyRate) (resp *core.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	_, err = sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	msg := daotypes.NewMsgColonyRate(address, gatewayAddress, rate)
	if err != nil {
		log.Error("NewMsgColonyRate")
		return
	}
	
	tx, result, err := this.TxClient.SignAndSendMsg(address, privateKey, core.NewLedgerFeeFromGas(flags.DefaultGasLimit, 0), "", msg)
	if err != nil {
		return
	}
	
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		resp = new(core.BroadcastTxResponse)
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return resp, nil
	} else {
		
		resp.TxHash = hex.EncodeToString(tx.Hash())
		return resp, errors.New(result.Info)
	}
}

// MultiMsgBroadcast 
func (this *GatewayClient) MultiMsgBroadcast(multiMsg []sdk.Msg, privateKey string) (tx ttypes.Tx, resp *core.BroadcastTxResponse, err error) {
	//logs := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	from, err := this.AccountClient.CreateAccountFromPriv(privateKey)
	if err != nil {
		return nil, nil, err
	}

	var result *core.BaseResponse
	
	txBase := core.TxBase{
		Fee:  legacytx.NewStdFee(0, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, sdk.ZeroInt()))),
		Memo: "",
	}

	tx, result, err = this.TxClient.SignAndSendMsg(from.Address, privateKey, txBase.Fee, txBase.Memo, multiMsg...)
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
