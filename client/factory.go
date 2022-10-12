package client

import (
	"bytes"
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/core"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	authType "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/pflag"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tharsis/ethermint/encoding"
)

var clientCtx client.Context

var clientFactory tx.Factory

var encodingConfig params.EncodingConfig

func NewEvmClient() EvmClient {
	return EvmClient{core.EvmRpcURL}
}

func NewTxClient() TxClient {
	return TxClient{core.ServerURL, "TxClient"}
}

func NewBlockClient() BlockClient {
	return BlockClient{core.ServerURL, "sc-BlockClient"}
}

func NewChatClient(txClient *TxClient, accClient *AccountClient) *ChatClient {
	return &ChatClient{txClient, accClient, core.ServerURL, "ChatClient"}
}

func NewAccountClient(txClient *TxClient) AccountClient {
	return AccountClient{txClient, NewSecretKey(), core.ServerURL, "AccountClient"}
}

func NewGatewayClinet() GatewayClient {
	return GatewayClient{core.ServerURL, "GatewayClient"}
}

func init() {
	encodingConfig = encoding.MakeConfig(app.ModuleBasics)

	rpcClient, err := rpchttp.New(core.RpcURL, "/websocket")
	if err != nil {
		panic("start ctx client error.")
	}

	clientCtx = client.Context{}.
		WithChainID(core.ChainID).
		WithCodec(encodingConfig.Marshaler).
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
}
