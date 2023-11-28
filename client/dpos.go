package client

import (
	"context"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/cosmos/cosmos-sdk/x/distribution/client/common"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"strings"
	"time"
)

type DposClient struct {
	TxClient  *TxClient
	ServerUrl string
	logPrefix string
}

func (this *DposClient) RegisterValidator(bech32DelegatorAddr, bech32ValidatorAddr string, bech32ValidatorPubkey cryptotypes.PubKey, selfDelegation sdk.Coin, desc stakingTypes.Description, commission stakingTypes.CommissionRates, minSelfDelegation sdk.Int, privateKey string, fee float64) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	_, err = sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return
	}

	
	err = commission.Validate()
	if err != nil {
		log.WithError(err).Error("commission.Validate")
		return
	}
	msg, err := stakingTypes.NewMsgCreateValidator(validatorAddr, bech32ValidatorPubkey, selfDelegation, desc, commission, minSelfDelegation)
	if err != nil {
		log.WithError(err).Error("NewMsgCreateValidator")
		return
	}
	_, resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, core.NewLedgerFee(fee), "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}

/*
*

*/
func (this *DposClient) EditorValidator(bech32ValidatorAccAddr string, desc stakingTypes.Description, newRate sdk.Dec, minSelfDelegation *sdk.Int, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	accAddr, err := sdk.AccAddressFromBech32(bech32ValidatorAccAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return
	}
	validatorAddress := sdk.ValAddress(accAddr).String()
	
	validatorInfor, err := this.FindValidatorByValAddress(validatorAddress)
	if err != nil {
		err = errors.New(core.QueryChainInforError)
		return
	}
	
	if validatorInfor.GetOperator().Empty() {
		err = stakingTypes.ErrNoValidatorFound
		return
	}
	if !minSelfDelegation.Equal(validatorInfor.MinSelfDelegation) {
		if !minSelfDelegation.GT(validatorInfor.MinSelfDelegation) {
			return nil, stakingTypes.ErrMinSelfDelegationDecreased
		}
		if minSelfDelegation.GT(validatorInfor.Tokens) {
			return nil, stakingTypes.ErrSelfDelegationBelowMinimum
		}
	} else {
		minSelfDelegation = nil
	}
	
	//err = validatorInfor.Commission.ValidateNewRate(validatorInfor.Commission.Rate, time.Now())
	if time.Now().Sub(validatorInfor.Commission.UpdateTime).Hours() < 24 {
		err = errors.New(core.ValidatorInfoError)
		return
	}
	_, err = desc.EnsureLength()
	if err != nil {
		log.WithError(err).Error("EnsureLength")
		err = errors.New(core.ValidatorDescriptionError)
		return
	}
	msg := stakingTypes.NewMsgEditValidator(validatorInfor.GetOperator(), desc, &newRate, minSelfDelegation)
	if err != nil {
		log.WithError(err).Error("NewMsgEditValidator")
		return
	}
	_, resp, err = this.TxClient.SignAndSendMsg(bech32ValidatorAccAddr, privateKey, core.NewLedgerFee(0), "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}

/*
*
 
*/
func (this *DposClient) FindValidatorByValAddress(bech32ValidatorAddr string) (validator *stakingTypes.Validator, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return
	}
	params := stakingTypes.QueryValidatorParams{ValidatorAddr: validatorAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/staking/"+stakingTypes.QueryValidator, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	validator = &stakingTypes.Validator{}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, validator)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
	}
	return
}


func (this *DposClient) UnjailValidator(bech32DelegatorAddr, bech32ValidatorAddr, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		return
	}
	
	validatorInfo, err := this.FindValidatorByValAddress(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("FindValidatorByValAddress")
		return
	}
	if !validatorInfo.Jailed {
		return
	}

	accAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return
	}
	OperatorAddress := sdk.ValAddress(accAddr).String()
	if validatorInfo.OperatorAddress != OperatorAddress {
		log.Error("validatorInfo.OperatorAddress:", validatorInfo.OperatorAddress, "|OperatorAddress:", OperatorAddress)
		return nil, core.ErrUnjailOperatorAddress
	}
	
	delegatorResponse, _, err := this.FindDelegation(bech32DelegatorAddr, bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("this.FindDelegatio", bech32DelegatorAddr, bech32ValidatorAddr)
		return
	}
	
	if delegatorResponse.Delegation.Shares.IsZero() {
		log.Error("delegatorResponse.Delegation.Shares.IsZero")
		return nil, core.ErrDelegationZero
	}

	tokens := validatorInfo.TokensFromShares(delegatorResponse.Delegation.Shares).TruncateInt()
	if tokens.LT(validatorInfo.MinSelfDelegation) {
		log.Error("tokens.LT(validatorInfo.MinSelfDelegation)")
		return nil, core.ErrDelegateAmountLtMinSelfDelegation
	}

	
	msg := slashingTypes.NewMsgUnjail(validatorAddr)
	fee := legacytx.NewStdFee(flags.DefaultGasLimit, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, sdk.NewInt(0))))
	_, resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, fee, "", msg)
	if err != nil {
		log.WithError(err).Error("this.TxClient.SignAndSendMsg")
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}

/*
*

*/
func (this *DposClient) FindDelegation(delegatorAddr, validatorAddr string) (delegation *stakingTypes.DelegationResponse, notFound bool, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	notFound = false
	params := stakingTypes.QueryDelegatorValidatorRequest{DelegatorAddr: delegatorAddr, ValidatorAddr: validatorAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, notFound, err
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/staking/"+stakingTypes.QueryDelegation, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		if strings.Contains(err.Error(), stakingTypes.ErrNoDelegation.Error()) {
			notFound = true
		}
		return nil, notFound, err
	}
	delegation = &stakingTypes.DelegationResponse{}
	err = util.Json.Unmarshal(resBytes, delegation)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}

type ValidatorInfo struct {
	stakingTypes.Validator
	StartTime           int64  `json:"start_time"`
	UnbondingTimeFormat string `json:"unbonding_time_format"`
}

/*
（）
@params status:

	BOND_STATUS_UNSPECIFIED
	BOND_STATUS_UNBONDED
	BOND_STATUS_UNBONDING
	BOND_STATUS_BONDED
*/
func (this *DposClient) QueryValidators(page, limit int, status string) ([]ValidatorInfo, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	if status == "" {
		status = stakingTypes.BondStatusBonded
	}

	params := stakingTypes.NewQueryValidatorsParams(page, limit, status)

	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("clientCtx.LegacyAmino.MarshalJSON")
		return nil, err
	}

	route := fmt.Sprintf("custom/%s/%s", stakingTypes.QuerierRoute, stakingTypes.QueryValidators)

	res, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, route, bz)
	if err != nil {
		log.WithError(err).Error("clientCtx.QueryWithData")
		return nil, err
	}

	validatorsResp := make(stakingTypes.Validators, 0)
	err = clientCtx.LegacyAmino.UnmarshalJSON(res, &validatorsResp)
	if err != nil {
		return nil, err
	}

	validatorInfos := make([]ValidatorInfo, 0)

	for _, val := range validatorsResp {

		valConsAddr, err := val.GetConsAddr()
		if err != nil {
			log.WithError(err).Error("GetConsAddr Err:" + err.Error())
			return nil, err
		}

		
		node, err := clientCtx.GetNode()
		if err != nil {
			log.WithError(err).Error("clientCtx.GetNode")
			return nil, err
		}
		slashingParams := slashingTypes.QuerySigningInfoRequest{valConsAddr.String()}
		bz, err := clientCtx.LegacyAmino.MarshalJSON(slashingParams)
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
			log.WithError(err).Error("clientCtx.LegacyAmino.UnmarshalJSON(slashingTypes.ValidatorSigningInfo)")
			return nil, err
		}
		nodeStatus, err := node.Status(context.Background())
		if err != nil {
			log.WithError(err).Error("node.Status")
			return nil, err
		}
		startHeight := nodeStatus.SyncInfo.LatestBlockHeight - info.IndexOffset
		if startHeight == 0 {
			startHeight = 1
		}
		blockInfo, err := node.Block(context.Background(), &startHeight)
		if err != nil {
			log.WithError(err).Error("node.Block")
			return nil, err
		}

		
		unbondingTimeFormat := ""
		if val.UnbondingTime.Unix() != 0 {
			unbondingTimeFormat = val.UnbondingTime.Format("2006-01-02 15:04")
		}

		ValidatorInfoEdit := ValidatorInfo{
			val,
			time.Now().Unix() - blockInfo.Block.Time.Unix(),
			unbondingTimeFormat,
		}

		ValidatorInfoEdit.ConsensusPubkey = nil
		validatorInfos = append(validatorInfos, ValidatorInfoEdit)

	}

	return validatorInfos, nil

}

type DecCoinsString struct {
	Coins []struct {
		Denom  string
		Amount string
	} `json:"coins"`
}

/*
*

*/
func (this *DposClient) QueryValCanWithdraw(accAddr string) (res distributionTypes.ValidatorAccumulatedCommission, err error) {

	//log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	accFromAddress, err := sdk.AccAddressFromBech32(accAddr)
	if err != nil {
		return res, err
	}
	valAddr := sdk.ValAddress(accFromAddress)

	// query commission
	bz, err := common.QueryValidatorCommission(clientCtx, valAddr)
	if err != nil {
		return res, err
	}

	var commission distributionTypes.ValidatorAccumulatedCommission
	clientCtx.LegacyAmino.UnmarshalJSON(bz, &commission)

	return commission, nil
}

/*
*

*/
func (this *DposClient) Delegation(bech32DelegatorAddr, bech32ValidatorAddr string, amount sdk.Coin, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return
	}
	msg := stakingTypes.NewMsgDelegate(delegatorAddr, validatorAddr, amount)
	fee := core.NewLedgerFee(0)
	_, resp, err = this.TxClient.SignAndSendMsg(msg.DelegatorAddress, privateKey, fee, "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}

//   amount:
func (this *DposClient) UnbondDelegation(bech32DelegatorAddr string, bech32ValidatorAddr string, amount sdk.Coin, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return
	}
	msg := stakingTypes.NewMsgUndelegate(delegatorAddr, validatorAddr, amount)
	fee := core.NewLedgerFee(0)
	_, resp, err = this.TxClient.SignAndSendMsg(msg.DelegatorAddress, privateKey, fee, "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) UnbondDelegationAll(bech32DelegatorAddr, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return
	}
	params := stakingTypes.NewQueryDelegatorParams(delegatorAddr)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryDelegatorDelegations, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	delegationResp := stakingTypes.DelegationResponses{}
	err = util.Json.Unmarshal(resBytes, &delegationResp)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return
	}
	if len(delegationResp) <= 0 {
		err = errors.New("delegation does not exist")
		return
	}
	msgs := []sdk.Msg{}
	for _, delegation := range delegationResp {
		valAddr, err := sdk.ValAddressFromBech32(delegation.Delegation.ValidatorAddress)
		if err != nil {
			continue
		}
		msg := stakingTypes.NewMsgUndelegate(delegatorAddr, valAddr, delegation.Balance)
		msgs = append(msgs, msg)
	}

	fee := core.NewLedgerFee(0)
	_, resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, fee, "", msgs...)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) DrawCommissionDelegationRewards(bech32DelegatorAddr, bech32ValidatorAddr, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return
	}
	bech32DelegatorValidatorAddr := sdk.ValAddress(delegatorAddr).String()

	
	delegationReward, validatorReward, err := this.RewardsPreview(bech32DelegatorAddr, bech32ValidatorAddr)
	if err != nil {
		return
	}

	//fmt.Println("validatorReward:", validatorReward, "delegationReward:", delegationReward)

	if validatorReward.IsZero() && delegationReward.IsZero() {
		err = core.ErrRewardReceive 
		return
	}
	msgs := []sdk.Msg{}

	if !delegationReward.IsZero() { 
		msg1 := distributionTypes.NewMsgWithdrawDelegatorReward(delegatorAddr, validatorAddr)
		msgs = append(msgs, msg1)
	}
	
	if bech32DelegatorValidatorAddr == bech32ValidatorAddr && !validatorReward.IsZero() { 
		
		msg2 := distributionTypes.NewMsgWithdrawValidatorCommission(validatorAddr)
		msgs = append(msgs, msg2) //tx
	}

	fee := core.NewLedgerFee(0)
	_, resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, fee, "", msgs...)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}

/*
*

delegatorReward 
validatorReward 
*/
func (this *DposClient) RewardsPreview(bech32DelegatorAddr string, bech32ValidatorAddr string) (delegatorReward sdk.Coin, validatorReward sdk.Coin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	/**
	1.
	2.
	3.1
	(ValidatorsDistributionTotal)
	3.2
	(ValidatorsDistributionReWards2)
	*/

	delegatorReward = sdk.NewCoin(core.GovDenom, sdk.NewInt(0))
	validatorReward = sdk.NewCoin(core.GovDenom, sdk.NewInt(0))
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return delegatorReward, validatorReward, errors.New(core.ParseAccountError)
	}
	validatorAddr, err := sdk.ValAddressFromBech32(bech32ValidatorAddr)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		return delegatorReward, validatorReward, errors.New(core.ParseAccountError)
	}
	params := distributionTypes.NewQueryDelegationRewardsParams(delegatorAddr, validatorAddr)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return delegatorReward, validatorReward, errors.New(core.QueryChainInforError)
	}
	
	resBytes, _, err := clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryDelegationRewards, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData error1 | ", err.Error())
		return delegatorReward, validatorReward, err
	}

	delegatorRewardDecCoins := sdk.DecCoins{} 

	err = util.Json.Unmarshal(resBytes, &delegatorRewardDecCoins)
	if err != nil {
		log.WithError(err).Error("Unmarshal error1 | ", err.Error())
		return delegatorReward, validatorReward, errors.New(core.QueryChainInforError)
	}

	if len(delegatorRewardDecCoins) != 0 {
		delegatorReward, _ = delegatorRewardDecCoins[0].TruncateDecimal() 
	}
	bech32DelegatorValidatorAddr := sdk.ValAddress(delegatorAddr).String()
	
	if bech32DelegatorValidatorAddr == bech32ValidatorAddr {
		
		resBytes, _, err = clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryValidatorCommission, bz)
		if err != nil {
			log.WithError(err).Error("QueryWithData error2 | ", err.Error())
			return delegatorReward, validatorReward, errors.New(core.QueryChainInforError)
		}
		validatorCommDecCoins := distributionTypes.ValidatorAccumulatedCommission{} 
		err = util.Json.Unmarshal(resBytes, &validatorCommDecCoins)
		if err != nil {
			log.WithError(err).Error("Unmarshal error2 | ", err.Error())
			return delegatorReward, validatorReward, errors.New(core.QueryChainInforError)
		}
		if len(validatorCommDecCoins.Commission) > 0 {
			validatorReward, _ = validatorCommDecCoins.Commission[0].TruncateDecimal() 
		}
	}

	return delegatorReward, validatorReward, err
}


func (this *DposClient) DrawCommissionDelegationRewardsAll(bech32DelegatorAddr, privateKey string) (resp *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return
	}
	totalRewards, ValAddressArr, err := this.RewardsPreviewAll(bech32DelegatorAddr)
	if err != nil {
		err = errors.New(core.ParseAccountError)
		return
	}
	if totalRewards.Total.IsZero() && len(ValAddressArr.ValidatorCommissions) <= 0 {
		err = errors.New("There is no reward to receive") 
		return
	}
	msgs := []sdk.Msg{}

	if !totalRewards.Total.IsZero() { 
		for _, reward := range totalRewards.Rewards {
			//log.Debug(":", reward.Reward.String())
			if reward.Reward.IsZero() {
				continue
			}
			if reward.Reward[0].Amount.LT(sdk.NewDec(1)) {
				continue
			}
			validatorAddr, err := sdk.ValAddressFromBech32(reward.ValidatorAddress)
			if err != nil {
				continue
			}
			msg1 := distributionTypes.NewMsgWithdrawDelegatorReward(delegatorAddr, validatorAddr)
			msgs = append(msgs, msg1)
		}
	}
	
	if len(ValAddressArr.ValidatorCommissions) > 0 { 
		
		for _, valAddress := range ValAddressArr.ValidatorCommissions {
			msg2 := distributionTypes.NewMsgWithdrawValidatorCommission(valAddress.ValidatorAddress)
			msgs = append(msgs, msg2) //tx
		}
	}

	fee := core.NewLedgerFee(0)
	_, resp, err = this.TxClient.SignAndSendMsg(bech32DelegatorAddr, privateKey, fee, "", msgs...)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}


func (this *DposClient) RewardsPreviewAll(bech32DelegatorAddr string) (distributionTypes.QueryDelegatorTotalRewardsResponse, stakingTypes.ValidatorCommissionResp, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	totalRewards := distributionTypes.QueryDelegatorTotalRewardsResponse{}
	var ValAddressArr stakingTypes.ValidatorCommissionResp
	var valAddrcommission stakingTypes.ValidatorCommission
	ValAddressArr.Total = sdk.DecCoins{}
	delegatorAddr, err := sdk.AccAddressFromBech32(bech32DelegatorAddr)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		err = errors.New(core.ParseAccountError)
		return totalRewards, ValAddressArr, nil
	}
	params := distributionTypes.NewQueryDelegatorParams(delegatorAddr)
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON error1 | ", err.Error())
		return totalRewards, ValAddressArr, nil
	}
	
	resBytes, _, err := clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryDelegatorTotalRewards, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData error1 | ", err.Error())
		return totalRewards, ValAddressArr, nil
	}

	err = util.Json.Unmarshal(resBytes, &totalRewards)
	if err != nil {
		log.WithError(err).Error("Unmarshal error1 | ", err.Error())
		return totalRewards, ValAddressArr, nil
	}
	for _, reward := range totalRewards.Rewards {
		validatorAddr, err := sdk.ValAddressFromBech32(reward.ValidatorAddress)
		if err != nil {
			log.WithError(err).Error("ValAddressFromBech32 error2 | ", err.Error())
			continue
		}
		params := distributionTypes.NewQueryDelegationRewardsParams(delegatorAddr, validatorAddr)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			log.WithError(err).Error("MarshalJSON error2 | ", err.Error())
			continue
		}
		bech32DelegatorValidatorAddr := sdk.ValAddress(delegatorAddr).String()
		
		if bech32DelegatorValidatorAddr == validatorAddr.String() {
			
			resBytes, _, err = clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryValidatorCommission, bz)
			if err != nil {
				log.WithError(err).Error("QueryWithData error2 | ", err.Error())
				continue
			}
			validatorComm := distributionTypes.ValidatorAccumulatedCommission{} 
			err = util.Json.Unmarshal(resBytes, &validatorComm)
			if err != nil {
				log.WithError(err).Error("Unmarshal error2 | ", err.Error())
				continue
			}
			if validatorComm.Commission.IsZero() {
				continue
			}
			valAddrcommission.ValidatorAddress = validatorAddr
			valAddrcommission.Reward = validatorComm.Commission
			for _, coin := range validatorComm.Commission {
				ValAddressArr.Total = ValAddressArr.Total.Add(coin)
			}
			ValAddressArr.ValidatorCommissions = append(ValAddressArr.ValidatorCommissions, valAddrcommission)
		}

	}
	return totalRewards, ValAddressArr, nil
}
