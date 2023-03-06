package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

func TestA(t *testing.T) {
	t.Log(sdk.NewInt(123).String())
}
