package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/cmd/stcd/cmd"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	contractTypes "freemasonry.cc/blockchain/x/contract/types"
	"freemasonry.cc/blockchain/x/gateway/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors2 "github.com/cosmos/cosmos-sdk/types/errors"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	types2 "github.com/cosmos/cosmos-sdk/x/auth/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	evmhd "github.com/evmos/ethermint/crypto/hd"
	"github.com/evmos/ethermint/encoding"
	evmoskr "github.com/evmos/evmos/v10/crypto/keyring"
	"github.com/tendermint/tendermint/crypto/tmhash"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	ttypes "github.com/tendermint/tendermint/types"
	"os"
	"regexp"
	"strings"
)

type TxClient struct {
	ServerUrl string
	logPrefix string
}


func (this TxClient) GasMsg(msgType, msgs string) (sdk.Msg, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	msgByte := []byte(msgs)
	if unmashal, ok := msgUnmashalHandles[msgType]; ok {
		msg, err := unmashal(msgByte)
		if err != nil {
			log.WithError(err).WithField("msgType", msgType).Error("gasMsg.Unmarshal MsgSend")
			return nil, err
		}
		return msg, nil
	} else {
		return nil, errors.New("unregister unmashal msg type:" + msgType)
	}
}


func (this *TxClient) GasInfo(seqDetail core.AccountNumberSeqResponse, msg ...sdk.Msg) (sdk.Coin, uint64, sdk.DecCoin, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	fee := sdk.NewCoin(core.BaseDenom, sdk.ZeroInt())
	gasPrice := sdk.NewDecCoin(core.BaseDenom, sdk.ZeroInt())
	minFee, err := sdk.ParseCoinNormalized(core.ChainDefaultFeeStr)
	if err != nil {
		return minFee, 0, gasPrice, err
	}

	
	if seqDetail.NotFound == true {
		return sdk.NewCoin(core.BaseDenom, sdk.NewInt(0)), 0, sdk.NewDecCoin(core.BaseDenom, sdk.NewInt(0)), errors2.ErrInsufficientFunds
	}

	clientFactory = clientFactory.WithSequence(seqDetail.Sequence).WithSimulateAndExecute(true)
	gasInfo, _, err := tx.CalculateGas(clientCtx, clientFactory, msg...)
	if err != nil {
		sp := strings.Split(err.Error(), ": ")
		size := len(sp)
		if size-2 >= 0 {
			err = errors.New(sp[size-2])
		} else if size-1 >= 0 {
			err = errors.New(sp[size-1])
		}
		log.WithError(err).Error("tx.CalculateGas")
		return minFee, 0, gasPrice, err
	}
	//log.Info("GasInfo:",gasInfo.GasInfo)
	//log.Info("Result:",gasInfo.Result)

	gas := gasInfo.GasInfo.GasUsed * 2
	//gas
	gasPriceDec, err := this.QueryGasPrice()
	if err != nil {
		log.WithError(err).Error("QueryGasPrice")
		return minFee, 0, gasPrice, err
	}

	if gasPriceDec.IsZero() {
		gasPriceDec = sdk.NewDecCoins(sdk.NewDecCoin(core.BaseDenom, sdk.NewInt(core.DefaultGasPrice)))
	}
	gasPrice = gasPriceDec[0]
	gasDec := sdk.NewDec(int64(gas))
	
	fee = sdk.NewCoin(core.BaseDenom, gasPrice.Amount.Mul(gasDec).TruncateInt())
	return fee, gas, gasPrice, nil
}

// gas
func (this *TxClient) QueryGasPrice() (gasPrice sdk.DecCoins, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+types.ModuleName+"/"+types.QueryGasPrice, nil)
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


func (this *TxClient) GetTranserTypeConfig() map[string]string {
	return core.GetTranserTypeConfig()
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

// tx tendermint tx  cosmos  tx
func (this *TxClient) TermintTx2CosmosTx(signTxs ttypes.Tx) (sdk.Tx, error) {
	return clientCtx.TxConfig.TxDecoder()(signTxs)
}

// tx
func (this *TxClient) SignTx2Bytes(signTxs xauthsigning.Tx) ([]byte, error) {
	return clientCtx.TxConfig.TxEncoder()(signTxs)
}

func (this *TxClient) SetFee(signTxs xauthsigning.Tx) ([]byte, error) {
	return clientCtx.TxConfig.TxEncoder()(signTxs)
}

// tx
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

// tx,,err=nil notFound=true
// notFound 
// err    
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

// tx
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

/*
*

*/
func (this *TxClient) SignAndSendMsg(address string, privateKey string, fee legacytx.StdFee, memo string, msg ...sdk.Msg) (tx ttypes.Tx, txRes *core.BaseResponse, err error) {
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
		return nil, nil, errors.New("sign error")
	}

	
	tx, err = this.SignTx2Bytes(signedTx)
	if err != nil {
		log.WithError(err).Error("SignTx2Bytes")
		return
	}
	log.Info("txhash:", hex.EncodeToString(tmhash.Sum(tx)))
	
	txRes, err = this.Send(tx)
	if txRes != nil {
		broadcastTxResponse := core.BroadcastTxResponse{}
		bytes, _ := json.Marshal(txRes.Data)
		err = json.Unmarshal(bytes, &broadcastTxResponse)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return
		}
		broadcastTxResponse.SignedTxStr = hex.EncodeToString(tx)
		txRes.Data = broadcastTxResponse
	}
	return
}

// SequenceAccountNumber
func (this *TxClient) FindAccountNumberSeq(accountAddr string) (core.AccountNumberSeqResponse, error) {
	//log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	//baseResp := core.BaseResponse{}
	//seq := core.AccountNumberSeqResponse{}
	//response, err := GetRequest(this.ServerUrl, "/chat/accountNumberSeq/"+accountAddr)
	//if err != nil {
	//	log.WithError(err).Error("GetRequest")
	//	return seq, err
	//}
	//err = json.Unmarshal([]byte(response), &baseResp)
	//if err != nil {
	//	log.WithError(err).Error("json.Unmarshal")
	//	return seq, err
	//}
	//if baseResp.Status == 1 {
	//	bytesM, err := json.Marshal(baseResp.Data)
	//	if err != nil {
	//		return seq, err
	//	}
	//	err = json.Unmarshal(bytesM, &seq)
	//	if err != nil {
	//		return seq, err
	//	}
	//}
	//return c, nil
	seq := core.AccountNumberSeqResponse{}
	accAddr, err := sdk.AccAddressFromBech32(accountAddr)
	if err != nil {
		return seq, err
	}
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	clientCtx = clientCtx.WithBroadcastMode(flags.BroadcastSync).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithCodec(encodingConfig.Codec).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types2.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome).
		WithKeyringOptions(evmoskr.Option()).
		WithViper(cmd.EnvPrefix).
		WithLedgerHasProtobuf(true)
	accountNumber, sequence, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, accAddr)
	if err != nil {
		seq.NotFound = true
		return seq, err
	}

	seq.AccountNumber = accountNumber
	seq.Sequence = sequence
	return seq, err

}


func (this *TxClient) Send(req []byte) (txRes *core.BaseResponse, err error) {
	//log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	var txBytes []byte
	if req != nil {
		txBytes = req
	}
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	clientCtx = clientCtx.WithBroadcastMode(flags.BroadcastBlock).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithCodec(encodingConfig.Codec).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types2.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithKeyringOptions(evmoskr.Option()).
		WithViper(cmd.EnvPrefix).
		WithLedgerHasProtobuf(true)
	res, err := clientCtx.BroadcastTx(txBytes)

	//fmt.Println("broadcast err:")
	//fmt.Println(err)
	//fmt.Println(res)

	txRes = &core.BaseResponse{}
	if err != nil {
		sp := strings.Split(err.Error(), ": ")
		txRes.Info = sp[len(sp)-1]
		return
	}

	txResponse := core.BroadcastTxResponse{}
	txResponse.Code = res.Code
	txResponse.CodeSpace = res.Codespace
	txResponse.TxHash = res.TxHash
	txResponse.Height = res.Height

	if res.Code == 0 { //code=0 
		txRes.Status = 1
	} else {
		txRes.Status = 0
	}

	txRes.Info = ParseErrorCode(res.Code, res.Codespace, res.RawLog)
	txRes.Data = txResponse

	return
}

type SetChatInfo struct {
	FromAddress        string `json:"from_address,omitempty" yaml:"from_address"`
	NodeAddress        string `json:"node_address,omitempty" yaml:"node_address"`
	AddressBook        string `json:"address_book,omitempty" yaml:"address_book"`
	ChatBlacklist      string `json:"chat_blacklist,omitempty" yaml:"chat_blacklist"`
	ChatRestrictedMode string `json:"chat_restricted_mode,omitempty" yaml:"chat_limit"`
	ChatFeeAmount      string `json:"chat_fee_amount" yaml:"chat_fee_amount"`
	ChatFeeCoinSymbol  string `json:"chat_fee_coin_symbol" yaml:"chat_fee_coin_symbol"`
	ChatWhitelist      string `json:"chat_whitelist,omitempty" yaml:"chat_whitelist"`
	UpdateTime         int64  `json:"update_time,omitempty" yaml:"update_time"`
	ChatBlacklistEnc   string `json:"chat_blacklist_enc,omitempty" yaml:"chat_blacklist_enc"`
	ChatWhitelistEnc   string `json:"chat_whitelist_enc,omitempty" yaml:"chat_whitelist_enc"`
	Remarks            string `json:"remarks,omitempty" yaml:"remarks"`
}

type CrossChainOut struct {
	SendAddress string `json:"send_address"`
	ToAddress   string `json:"to_address"`
	CoinAmount  string `json:"coin_amount"`
	CoinSymbol  string `json:"coin_symbol"`
	ChainType   string `json:"chain_type"`
	Remark      string `json:"remark"`
}

/*
*
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
		feeAmount, gas, _, err := this.GasInfo(seqDetail, msgs...)
		if err != nil {
			log.WithError(err).Error("CulGas")
			return nil, util.ErrFilter(err)
		}
		log.WithField("gas", gas).Info("CulGas:")
		fee.Gas = gas
		fee.Amount = sdk.NewCoins(feeAmount)
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

func (this *TxClient) QueryProposer(proposerId uint64) (resp *govTypes.Proposal, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := govTypes.QueryProposalParams{ProposalID: proposerId}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/gov/"+govTypes.QueryProposal, bz)
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

func (this *TxClient) ChatTokenIssue(fromAddress, privateKey string) (resp *core.BaseResponse, err error) {
	msg := contractTypes.NewMsgChatTokenIssue(
		fromAddress,
		"testname",
		"zzz",
		"1000000000000",
		"6",
		"150000",
		"300000",
		"130000",
		"200000",
		"670000",
		"100000",
		"10000000000",
		"http://img",
	)
	fee := legacytx.NewStdFee(flags.DefaultGasLimit, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, sdk.NewInt(1))))
	_, resp, err = this.SignAndSendMsg(fromAddress, privateKey, fee, "", msg)
	if err != nil {
		return
	}
	if resp.Status == 1 {
		return resp, nil 
	} else {
		return resp, errors.New(resp.Info) 
	}
}

func ParseErrorCode(code uint32, codeSpace string, rowlog string) string {
	if codeSpace == sdkErrors.RootCodespace {
		if code == sdkErrors.ErrInsufficientFee.ABCICode() { 
			return core.FeeIsTooLess
		} else if code == sdkErrors.ErrOutOfGas.ABCICode() { //gas
			return core.ErrorGasOut
		} else if code == sdkErrors.ErrUnauthorized.ABCICode() { //idaccount number 
			return core.ErrUnauthorized
		} else if code == sdkErrors.ErrWrongSequence.ABCICode() { 
			return core.ErrWrongSequence
		}
	}

	//（grpc）
	sp := strings.Split(rowlog, ": ")
	return sp[len(sp)-1]
}
