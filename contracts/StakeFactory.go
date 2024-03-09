package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/StakeFactory.json
	stakeFactoryJSON []byte

	// FreeMasonryMedalContract is the compiled FreeMasonryMedal contract
	StakeFactoryJSONContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(stakeFactoryJSON, &StakeFactoryJSONContract)
	if err != nil {
		panic(err)
	}
}
