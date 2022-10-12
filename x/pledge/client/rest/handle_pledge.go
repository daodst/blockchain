package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

// RegisterHandlerFn 
func PledgeHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	return nil
}

func UnpledgeHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	return nil
}

func PledgeReceiveHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {

	return nil
}
