package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/ExchangeFactory.json
	ExchangeWTTJSON []byte // nolint: golint

	ExchangeWTTContract evmtypes.CompiledContract
)

func init() {

	err := json.Unmarshal(ExchangeWTTJSON, &ExchangeWTTContract)
	if err != nil {
		panic(err)
	}

}
