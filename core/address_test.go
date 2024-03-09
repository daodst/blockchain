package abi

import (
	"github.com/ethereum/go-ethereum/common"
	"testing"
)


func TestAddressDecode(t *testing.T) {
	from := `0x0000000000000000000000002c41cb1b24a8a75e4b324bdb8eda986630de6503`
	addr := common.HexToAddress(from)
	t.Log(addr)

	to := `0x000000000000000000000000674fdee4905f5755b824f3e70f2a34099d8121fa`
	addr = common.HexToAddress(to)
	t.Log(addr)

	contractAddr := `0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef`
	addr = common.HexToAddress(contractAddr)
	t.Log(addr)
}
