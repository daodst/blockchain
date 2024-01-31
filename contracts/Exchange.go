package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/Exchange.json
	ExchangeRouterJSON []byte // nolint: golint

	ExchangeRouterContract evmtypes.CompiledContract
)

func init() {

	err := json.Unmarshal(ExchangeRouterJSON, &ExchangeRouterContract)
	if err != nil {
		panic(err)
	}

}
