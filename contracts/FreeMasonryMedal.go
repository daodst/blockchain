package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/FreeMasonryMedal.json
	freeMasonryMedalJSON []byte

	// FreeMasonryMedalContract is the compiled FreeMasonryMedal contract
	FreeMasonryMedalJSONContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(freeMasonryMedalJSON, &FreeMasonryMedalJSONContract)
	if err != nil {
		panic(err)
	}
}
