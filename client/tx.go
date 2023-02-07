package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/cmd/config"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	chatTypes "freemasonry.cc/blockchain/x/chat/types"
	"freemasonry.cc/blockchain/x/comm/types"
	pledgeTypes "freemasonry.cc/blockchain/x/pledge/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/client/common"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakeTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmhd "github.com/evmos/ethermint/crypto/hd"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	ttypes "github.com/tendermint/tendermint/types"
	"regexp"
	"strings"
	"time"
)

type TxClient struct {
	ServerUrl string
	logPrefix string
}

func (this *TxClient) ConvertTxToStdTx(cosmosTx sdk.Tx) (*legacytx.StdTx, error) {
	signingTx, ok := cosmosTx.(xauthsigning.Tx)
	if !ok {
		return nil, errors.New("tx to stdtx error")
	}
	stdTx, err := tx.ConvertTxToStdTx(clientCtx.LegacyAmino, signingTx)
	if err != nil {
		return nil, err
	}
	return &stdTx, nil
}

//tx tendermint tx  cosmos  tx
func (this *TxClient) TermintTx2CosmosTx(signTxs ttypes.Tx) (sdk.Tx, error) {
	return clientCtx.TxConfig.TxDecoder()(signTxs)
}

//tx
func (this *TxClient) SignTx2Bytes(signTxs xauthsigning.Tx) ([]byte, error) {
	return clientCtx.TxConfig.TxEncoder()(signTxs)
}

func (this *TxClient) SetFee(signTxs xauthsigning.Tx) ([]byte, error) {
	return clientCtx.TxConfig.TxEncoder()(signTxs)
}

//tx
func (this *TxClient) FindByByte(txhash []byte) (resultTx *ctypes.ResultTx, notFound bool, err error) {
	notFound = false
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return
	}
	resultTx, err = node.Tx(context.Background(), txhash, true)
	if err != nil {
		
		notFound = this.isTxNotFoundError(err.Error())
		if notFound {
			err = nil
		} else {
			log.WithError(err).WithField("txhash", hex.EncodeToString(txhash)).Error("node.Tx")
		}
		return
	}
	return
}

//tx,,err=nil notFound=true
//notFound 
//err    
func (this *TxClient) FindByHex(txhashStr string) (resultTx *ctypes.ResultTx, notFound bool, err error) {
	var txhash []byte
	notFound = false
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	txhash, err = hex.DecodeString(txhashStr)
	if err != nil {
		log.WithError(err).WithField("txhash", txhashStr).Error("hex.DecodeString")
		return
	}
	return this.FindByByte(txhash)
}

//tx
func (this *TxClient) isTxNotFoundError(errContent string) (ok bool) {
	errRegexp := `tx\ \([0-9A-Za-z]{64}\)\ not\ found`
	r, err := regexp.Compile(errRegexp)
	if err != nil {
		return false
	}
	if r.Match([]byte(errContent)) {
		return true
	} else {
		return false
	}
}

/**

*/
func (this *TxClient) SignAndSendMsg(address string, privateKey string, fee legacytx.StdFee, memo string, msg ...sdk.Msg) (txRes *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	seqDetail, err := this.FindAccountNumberSeq(address)
	if err != nil {
		return
	}

	
	signedTx, err := this.SignTx(privateKey, seqDetail, fee, memo, msg...)
	if err != nil {
		return
	}

	
	signPubkey, err := signedTx.GetPubKeys()
	if err != nil {
		return
	}

	signV2, _ := signedTx.GetSignaturesV2()
	senderAddrBytes := signV2[0].PubKey.Address().Bytes()
	signAddrBytes := signPubkey[0].Address().Bytes()
	if !bytes.Equal(signAddrBytes, senderAddrBytes) {
		return nil, errors.New("sign error")
	}

	
	signedTxBytes, err := this.SignTx2Bytes(signedTx)
	if err != nil {
		log.WithError(err).Error("SignTx2Bytes")
		return
	}
	
	txRes, err = this.Send(signedTxBytes)
	if txRes != nil {
		broadcastTxResponse := txRes.Data.(core.BroadcastTxResponse)
		broadcastTxResponse.SignedTxStr = hex.EncodeToString(signedTxBytes)
		txRes.Data = broadcastTxResponse
	}
	return
}

/**

*/
func (this *TxClient) SignAndSendMsg2(address string, privateKey string, fee legacytx.StdFee, memo string, msg ...sdk.Msg) (txRes *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	seqDetail, err := this.FindAccountNumberSeq(address)
	if err != nil {
		return
	}

	
	signedTx, err := this.SignTx(privateKey, seqDetail, fee, memo, msg...)
	if err != nil {
		return
	}

	
	signPubkey, err := signedTx.GetPubKeys()
	if err != nil {
		return
	}

	signV2, _ := signedTx.GetSignaturesV2()
	senderAddrBytes := signV2[0].PubKey.Address().Bytes()
	signAddrBytes := signPubkey[0].Address().Bytes()
	if !bytes.Equal(signAddrBytes, senderAddrBytes) {
		return nil, errors.New("sign error")
	}

	
	signedTxBytes, err := this.SignTx2Bytes(signedTx)
	if err != nil {
		log.WithError(err).Error("SignTx2Bytes")
		return
	}
	
	txRes, err = this.Send(signedTxBytes)
	if txRes != nil {
		//broadcastTxResponse := txRes.Data.(core.BroadcastTxResponse)
		//broadcastTxResponse.SignedTxStr = hex.EncodeToString(signedTxBytes)
		//txRes.Data = broadcastTxResponse
	}
	return
}

//SequenceAccountNumber
func (this *TxClient) FindAccountNumberSeq(accountAddr string) (core.AccountNumberSeqResponse, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	baseResp := core.BaseResponse{}
	seq := core.AccountNumberSeqResponse{}
	response, err := GetRequest(this.ServerUrl, "/chat/accountNumberSeq/"+accountAddr)
	if err != nil {
		log.WithError(err).Error("GetRequest")
		return seq, err
	}
	err = json.Unmarshal([]byte(response), &baseResp)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal")
		return seq, err
	}
	if baseResp.Status == 1 {
		bytes, err := json.Marshal(baseResp.Data)
		if err != nil {
			return seq, err
		}
		err = json.Unmarshal(bytes, &seq)
		if err != nil {
			return seq, err
		}
	}
	return seq, nil
}


func (this *TxClient) Send(req []byte) (txRes *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	response, err := PostRequest(this.ServerUrl, "/chat/tx/broadcast", req)
	if err != nil {
		log.WithError(err).Error("PostRequest")
		return
	}
	txRes = &core.BaseResponse{}
	err = json.Unmarshal([]byte(response), txRes)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal")
		return
	}
	return
}


func (this TxClient) GasMsg(msgType, msgs string) (sdk.Msg, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	msgByte := []byte(msgs)
	switch msgType {
	case "cosmos-sdk/MsgSend":
		msg := bankTypes.MsgSend{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgSend")
			return nil, err
		}
		return &msg, nil
	case "pledge/MsgPledge":
		msg := pledgeTypes.MsgPledge{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgPledge")
			return nil, err
		}
		return &msg, nil
	case "pledge/MsgUnpledge":
		msg := pledgeTypes.MsgUnpledge{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgUnpledge")
			return nil, err
		}
		return &msg, nil
	case "pledge/MsgWithdrawDelegatorReward":
		msg := pledgeTypes.MsgWithdrawDelegatorReward{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgWithdrawDelegatorReward")
			return nil, err
		}
		return &msg, nil
	case "chat/MsgRegister":
		msg := chatTypes.MsgRegister{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgRegister")
			return nil, err
		}
		return &msg, nil
	case "chat/MsgSendGift":
		msg := chatTypes.MsgSendGift{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgSendGift")
			return nil, err
		}
		return &msg, nil
	case "chat/MsgChatSendGift":
		msg := chatTypes.MsgChatSendGift{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgChatSendGift")
			return nil, err
		}
		return &msg, nil
	case "chat/MsgMobileTransfer":
		msg := chatTypes.MsgMobileTransfer{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgMobileTransfer")
			return nil, err
		}
		return &msg, nil
	case "chat/MsgBurnGetMobile":
		msg := chatTypes.MsgBurnGetMobile{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgBurnGetMobile")
			return nil, err
		}
		return &msg, nil
	case "chat/MsgSetChatInfo":
		msg := chatTypes.MsgSetChatInfo{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgSetChatInfo")
			return nil, err
		}
		return &msg, nil
	case "cosmos-sdk/MsgDelegate":
		msg := stakeTypes.MsgDelegate{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgDelegate")
			return &msg, nil
		}
		return &msg, nil
	case "cosmos-sdk/MsgUndelegate":
		msg := stakeTypes.MsgUndelegate{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgUndelegate")
			return &msg, nil
		}
		return &msg, nil
	case "cosmos-sdk/MsgWithdrawDelegationReward":
		msg := distributionTypes.MsgWithdrawDelegatorReward{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgWithdrawDelegationReward")
			return &msg, nil
		}
		return &msg, nil
	case "cosmos-sdk/MsgVote":
		msg := govTypes.MsgVote{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgVote")
			return &msg, nil
		}
		return &msg, nil
	}

	return nil, nil
}


func (this *TxClient) GasInfo(seqDetail core.AccountNumberSeqResponse, msg ...sdk.Msg) (sdk.Coin, uint64, sdk.DecCoin, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
	fee := sdk.NewCoin(config.BaseDenom, sdk.ZeroInt())
	gasPrice := sdk.NewDecCoin(config.BaseDenom, sdk.ZeroInt())
	minFee, err := sdk.ParseCoinNormalized(core.ChainDefaultFeeStr)
	if err != nil {
		return minFee, 0, gasPrice, err
	}
	clientFactory = clientFactory.WithSequence(seqDetail.Sequence)
	gasInfo, _, err := tx.CalculateGas(clientCtx, clientFactory, msg...)
	if err != nil {
		log.WithError(err).Error("tx.CalculateGas")
		return minFee, 0, gasPrice, err
	}

	gas := gasInfo.GasInfo.GasUsed * 2
	//gas
	gasPriceDec, err := this.QueryGasPrice()
	if err != nil {
		log.WithError(err).Error("QueryGasPrice")
		return minFee, 0, gasPrice, err
	}

	if gasPriceDec.IsZero() {
		return minFee, gas, gasPrice, nil
	}
	gasPrice = gasPriceDec[0]
	gasDec := sdk.NewDec(int64(gas))
	
	fee = sdk.NewCoin(config.BaseDenom, gasPrice.Amount.Mul(gasDec).TruncateInt())
	if fee.IsLT(minFee) {
		
		return minFee, gas, gasPrice, nil
	}
	return fee, gas, gasPrice, nil
}

//gas
func (this *TxClient) QueryGasPrice() (gasPrice sdk.DecCoins, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, _, err := clientCtx.QueryWithData("custom/comm/"+types.QueryGasPrice, nil)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	if resBytes != nil {
		err := util.Json.Unmarshal(resBytes, &gasPrice)
		if err != nil {
			return nil, err
		}
	}
	return
}

/**
tx
*/
func (this *TxClient) SignTx(privateKey string, seqDetail core.AccountNumberSeqResponse, fee legacytx.StdFee, memo string, msgs ...sdk.Msg) (xauthsigning.Tx, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	privKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		log.WithError(err).Error("hex.DecodeString")
		return nil, err
	}
	keyringAlgos := keyring.SigningAlgoList{evmhd.EthSecp256k1}
	algo, err := keyring.NewSigningAlgoFromString("eth_secp256k1", keyringAlgos)
	if err != nil {
		return nil, err
	}
	privKey := algo.Generate()(privKeyBytes)
	//gas,gas
	if fee.Gas == flags.DefaultGasLimit {
		_, gas, _, err := this.GasInfo(seqDetail, msgs...)
		if err != nil {
			log.WithError(err).Error("CulGas")
			return nil, core.Errformat(err)
		}
		log.WithField("gas", gas).Info("CulGas:")
		fee.Gas = gas
	}
	signMode := clientCtx.TxConfig.SignModeHandler().DefaultMode()
	signerData := xauthsigning.SignerData{
		ChainID:       clientCtx.ChainID,
		AccountNumber: seqDetail.AccountNumber,
		Sequence:      seqDetail.Sequence,
	}
	txBuild, err := tx.BuildUnsignedTx(clientFactory, msgs...)
	if err != nil {
		log.WithError(err).Error("tx.BuildUnsignedTx")
		return nil, err
	}
	txBuild.SetGasLimit(fee.Gas)     //gas
	txBuild.SetFeeAmount(fee.Amount) 
	txBuild.SetMemo(memo)            
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   privKey.PubKey(),
		Data:     &sigData,
		Sequence: seqDetail.Sequence,
	}
	
	if err := txBuild.SetSignatures(sig); err != nil {
		log.WithError(err).Error("SetSignatures")
		return nil, err
	}
	signV2, err := tx.SignWithPrivKey(signMode, signerData, txBuild, privKey, clientCtx.TxConfig, seqDetail.Sequence)
	if err != nil {
		log.WithError(err).Error("SignWithPrivKey")
		return nil, err
	}
	err = txBuild.SetSignatures(signV2)
	if err != nil {
		log.WithError(err).Error("SetSignatures")
		return nil, err
	}

	signedTx := txBuild.GetTx()
	//fmt.Println("getSigners:",signedTx.GetSigners())
	return signedTx, nil
}

func (this *TxClient) RegisterValidator(bech32DelegatorAddr, bech32ValidatorAddr string, bech32ValidatorPubkey cryptotypes.PubKey, selfDelegation sdk.Coin, desc stakingTypes.Description, commission stakingTypes.CommissionRates, minSelfDelegation sdk.Int, privateKey string, fee float64) (resp *core.BaseResponse, err error) {
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
	resp, err = this.SignAndSendMsg(bech32DelegatorAddr, privateKey, core.NewLedgerFee(fee), "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}

func (this *TxClient) QueryProposer(proposerId uint64) (resp *govTypes.Proposal, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := govTypes.QueryProposalParams{ProposalID: proposerId}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/gov/"+govTypes.QueryProposal, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	var res govTypes.Proposal
	if resBytes != nil {
		err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
		if err != nil {
			return
		}
		resp = &res
	}
	return
}


func (this *TxClient) UnjailValidator(bech32DelegatorAddr, bech32ValidatorAddr, privateKey string) (resp *core.BaseResponse, err error) {
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
		return
	}
	
	delegatorResponse, notFound, err := this.FindDelegation(bech32DelegatorAddr, bech32ValidatorAddr)
	if err != nil {
		if notFound {
		} else {
		}
		return
	}
	
	if delegatorResponse.Delegation.Shares.IsZero() {
		return
	}

	tokens := validatorInfo.TokensFromShares(delegatorResponse.Delegation.Shares).TruncateInt()
	if tokens.LT(validatorInfo.MinSelfDelegation) {
		return
	}

	
	msg := slashingTypes.NewMsgUnjail(validatorAddr)
	fee := legacytx.NewStdFee(flags.DefaultGasLimit, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1))))
	resp, err = this.SignAndSendMsg2(bech32DelegatorAddr, privateKey, fee, "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
	return nil, nil
}

/**
 
*/
func (this *TxClient) FindValidatorByValAddress(bech32ValidatorAddr string) (validator *stakingTypes.Validator, err error) {
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
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryValidator, bz)
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

/**

*/
func (this *TxClient) FindDelegation(delegatorAddr, validatorAddr string) (delegation *stakingTypes.DelegationResponse, notFound bool, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	notFound = false
	params := stakingTypes.QueryDelegatorValidatorRequest{DelegatorAddr: delegatorAddr, ValidatorAddr: validatorAddr}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, notFound, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryDelegation, bz)
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
func (this *TxClient) QueryValidators(page, limit int, status string) ([]ValidatorInfo, error) {
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

	res, _, err := clientCtx.QueryWithData(route, bz)
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

		resBytes, _, err := clientCtx.QueryWithData("custom/slashing/"+slashingTypes.QuerySigningInfo, bz)
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
		if info.StartHeight == 0 {
			info.StartHeight = 1
		}
		blockInfo, err := node.Block(context.Background(), &info.StartHeight)
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

/**

*/
func (this *TxClient) QueryValCanWithdraw(accAddr string) (res distributionTypes.ValidatorAccumulatedCommission, err error) {

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
