package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/ExchangeFactory.json
	ExchangeFactoryJSON []byte // nolint: golint

	ExchangeFactoryContract evmtypes.CompiledContract
)

func init() {

	err := json.Unmarshal(ExchangeFactoryJSON, &ExchangeFactoryContract)
	if err != nil {
		panic(err)
	}

}
