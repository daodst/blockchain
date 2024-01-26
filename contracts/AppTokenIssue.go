package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/AppTokenIssue.json
	appTokenIssueJSON []byte

	// FreeMasonryMedalContract is the compiled FreeMasonryMedal contract
	AppTokenIssueJSONContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(appTokenIssueJSON, &AppTokenIssueJSONContract)
	if err != nil {
		panic(err)
	}
}
