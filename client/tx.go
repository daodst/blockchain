package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	ttypes "github.com/tendermint/tendermint/types"
	evmhd "github.com/tharsis/ethermint/crypto/hd"
	"regexp"
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
		//
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
	//
	seqDetail, err := this.FindAccountNumberSeq(address)
	if err != nil {
		return
	}

	//
	signedTx, err := this.SignTx(privateKey, seqDetail, fee, memo, msg...)
	if err != nil {
		return
	}

	//，
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

	//
	signedTxBytes, err := this.SignTx2Bytes(signedTx)
	if err != nil {
		log.WithError(err).Error("SignTx2Bytes")
		return
	}
	//
	txRes, err = this.Send(signedTxBytes)
	if txRes != nil {
		broadcastTxResponse := txRes.Data.(core.BroadcastTxResponse)
		broadcastTxResponse.SignedTxStr = hex.EncodeToString(signedTxBytes)
		txRes.Data = broadcastTxResponse
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

//
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

//
func (this TxClient) GasMsg(msgType, msgs string) (sdk.Msg, error) {
	log := util.BuildLog(util.GetStructFuncName(this), util.LmChainClient)
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
	case "chat/MsgSetChatFee":
		msg := chatTypes.MsgSetChatFee{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgSetChatFee")
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
	case "chat/MsgAddressBookSave":
		msg := chatTypes.MsgAddressBookSave{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgAddressBookSave")
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
	case "chat/MsgChangeGateway":
		msg := chatTypes.MsgChangeGateway{}
		err := util.Json.Unmarshal(msgByte, &msg)
		if err != nil {
			log.WithError(err).Error("gasMsg.Unmarshal MsgChangeGateway")
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
	}
	return nil, nil
}

//
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
	//
	fee = sdk.NewCoin(config.BaseDenom, gasPrice.Amount.Mul(gasDec).TruncateInt())
	if fee.IsLT(minFee) {
		//,
		return minFee, gas, gasPrice, nil
	}
	return fee, gas, gasPrice, nil
}

//gas
func (this *TxClient) QueryGasPrice() (gasPrice sdk.DecCoins, err error) {
	log := util.BuildLog(util.GetStructFuncName(this), util.LmChainClient)
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
	txBuild.SetFeeAmount(fee.Amount) //
	txBuild.SetMemo(memo)            //
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   privKey.PubKey(),
		Data:     &sigData,
		Sequence: seqDetail.Sequence,
	}
	//
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

	//
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
		return resp, nil //
	} else {
		return resp, errors.New(resp.Info) //
	}
}
