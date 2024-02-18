package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/Stake.json
	StakeJSON []byte // nolint: golint

	// ERC20MinterBurnerDecimalsContract is the compiled erc20 contract
	StakeContract evmtypes.CompiledContract
)

func init() {

	err := json.Unmarshal(StakeJSON, &StakeContract)
	if err != nil {
		panic(err)
	}

}
