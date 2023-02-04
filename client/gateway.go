package client

import (
	"context"
	"errors"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/comm/client/rest"
	"freemasonry.cc/blockchain/x/comm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/sirupsen/logrus"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"strconv"
	"strings"
	"time"
)

type GatewayClient struct {
	ServerUrl string
	logPrefix string
}


func (this *GatewayClient) StatusInfo() (statusInfo *ctypes.ResultStatus, err error) {
	node, err := clientCtx.GetNode()
	return node.Status(context.Background())
}


func (this *GatewayClient) NetInfo() (statusInfo *ctypes.ResultNetInfo, err error) {
	node, err := clientCtx.GetNode()
	return node.NetInfo(context.Background())
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
	resBytes, _, err := clientCtx.QueryWithData("custom/comm/"+types.QueryGatewayNumberCount, bz)
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
	resBytes, _, err := clientCtx.QueryWithData("custom/comm/"+types.QueryGatewayInfo, bz)
	if err != nil {
		
		if strings.Contains(err.Error(), types.ErrGatewayNotExist.Error()) {
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
	resBytes, _, err := clientCtx.QueryWithData("custom/comm/"+types.QueryGatewayList, nil)
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
		resBytes, _, err = clientCtx.QueryWithData("custom/comm/"+types.QueryValidators, nil)
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
		dataInfo := types.GatewayListResp{}
		for _, gateway := range resp {
			for _, val := range validator {
				if val.OperatorAddress == gateway.GatewayAddress {
					dataInfo.Gateway = gateway
					dataInfo.Token = val.Tokens
					dataInfo.Online = time.Now().Unix() - val.UnbondingTime
					if val.UnbondingTime == 0 {
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
						resBytes, _, err := clientCtx.QueryWithData("custom/slashing/"+slashingTypes.QuerySigningInfo, bz)
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
						if info.StartHeight == 0 {
							info.StartHeight = 1
						}
						blockInfo, err := node.Block(context.Background(), &info.StartHeight)
						if err != nil {
							log.WithError(err).Error("node.Block")
							return nil, err
						}
						dataInfo.Online = time.Now().Unix() - blockInfo.Block.Time.Unix()
					}
					data = append(data, dataInfo)
				}
			}
		}

	}
	return
}


func (this *GatewayClient) QueryGatewayNumList() (map[string]types.GatewayNumIndex, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, _, err := clientCtx.QueryWithData("custom/comm/"+types.QueryGatewayNum, nil)
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
	resBytes, _, err := clientCtx.QueryWithData("custom/comm/"+types.QueryGatewayRedeemNum, nil)
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

//  ValidatorStatus:  0 Unbonded , 1 Unbonding , 2 Bonded , 3  , 4 
func (this *GatewayClient) ValidatorInfo() (validatorInfo *types.ValidatorInfor, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	nodeStatus, err := this.StatusInfo()
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

/**
 
*/
func (this *GatewayClient) FindValidatorByConsAddress(bech32ConsAddr string) (validator *stakingTypes.Validator, notFound bool, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	notFound = false
	consAddress, err := sdk.ConsAddressFromBech32(bech32ConsAddr)
	if err != nil {
		log.WithError(err).Error("ConsAddressFromBech32")
		err = errors.New(rest.ParseAccountError)
		return
	}
	params := types.QueryValidatorByConsAddrParams{ValidatorConsAddress: consAddress}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, notFound, err
	}

	resBytes, _, err := clientCtx.QueryWithData("custom/comm/"+types.QueryValidatorByConsAddress, bz)
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
		return 3
	} else {
		switch status.String() {
		case "BOND_STATUS_UNBONDED":
			return 0 
		case "BOND_STATUS_UNBONDING":
			return 1 
		case "BOND_STATUS_BONDED":
			return 2 
		}
	}
	return 0
}
