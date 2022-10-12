package client

import (
	"github.com/tharsis/evmos/v4/contracts"
	"testing"
)

func TestToken(t *testing.T) {
	erc20 := contracts.ERC20BurnableContract.ABI
	dd := []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA2Ncmtxd6gAAA=")
	transferEvent, err := erc20.Unpack("Transfer", dd)
	if err != nil {
		t.Error("failed to unpack transfer event", "error", err.Error())
	}

	if len(transferEvent) == 0 {
	}
}
