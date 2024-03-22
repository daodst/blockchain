package client

import (
	"bytes"
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/dao/keeper"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authType "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/evmos/ethermint/encoding"
	"github.com/spf13/pflag"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

var clientCtx client.Context
var ClientCtx *client.Context

var clientFactory tx.Factory

var encodingConfig params.EncodingConfig

func NewEvmClient() EvmClient {
	return EvmClient{core.EvmRpcURL}
}

func NewTxClient() TxClient {
	return TxClient{core.ServerURL, "TxClient"}
}

//func NewMempoolClient() MempoolClient {
//	return MempoolClient{core.ServerURL, "TxClient"}
//}

func NewBlockClient() BlockClient {
	return BlockClient{core.ServerURL, "sc-BlockClient"}
}

func NewChatClient(txClient *TxClient, accClient *AccountClient) *ChatClient {
	return &ChatClient{txClient, accClient, core.ServerURL, "ChatClient"}
}

func NewAccountClient(txClient *TxClient) AccountClient {
	return AccountClient{txClient, NewSecretKey(), core.ServerURL, "AccountClient"}
}

func NewGatewayClinet(txClient *TxClient) GatewayClient {
	acountClient := NewAccountClient(txClient)
	return GatewayClient{txClient, &acountClient, core.ServerURL, "GatewayClient"}
}

func NewDposClinet(txClient *TxClient) DposClient {
	return DposClient{txClient, core.ServerURL, "DposClient"}
}

func NewClusterClient(txClient *TxClient) ClusterClient {
	return ClusterClient{
		TxClient:  txClient,
		ServerUrl: core.ServerURL,
		logPrefix: "ClusterClient",
	}
}

func init() {
	encodingConfig = encoding.MakeConfig(app.ModuleBasics)
	keeper.EncodingConfig = encodingConfig
	rpcClient, err := rpchttp.New(core.RpcURL, "/websocket")
	if err != nil {
		panic("start ctx client error.")
	}

	clientCtx = client.Context{}.
		WithChainID(core.ChainID).
		WithCodec(encodingConfig.Codec).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithOffline(true).
		WithNodeURI(core.RpcURL).
		WithClient(rpcClient).
		WithAccountRetriever(authType.AccountRetriever{})

	//cfg := network.DefaultConfig()

	//clientCtx = clientCtx.WithLegacyAmino(cfg.LegacyAmino)

	flags := pflag.NewFlagSet("chat", pflag.ContinueOnError)

	flagErrorBuf := new(bytes.Buffer)

	flags.SetOutput(flagErrorBuf)

	//gas  gas 
	clientFactory = tx.NewFactoryCLI(clientCtx, flags)
	clientFactory.WithChainID(core.ChainID).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithTxConfig(clientCtx.TxConfig)
	ClientCtx = &clientCtx
}


func MsgToStruct(msg sdk.Msg, obj interface{}) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	msgByte, err := encodingConfig.Amino.Marshal(msg)
	if err != nil {
		log.WithError(err).Error("MarshalBinaryBare")
		return err
	}
	err = encodingConfig.Amino.Unmarshal(msgByte, obj)
	if err != nil {
		log.WithError(err).Error("UnmarshalBinaryBare")
		return err
	}
	return nil
}
