package client

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func TestAccoun(t *testing.T) {
	dexAccount := "dex1tfnsctfjml9lskehnm5phvuxvpt45y20ythksh"
	ethAccount := "0xb2A559AD4E77f158da5dBF9BeA1303BEF9dB1b64"

	dex, _ := sdk.AccAddressFromBech32(dexAccount)
	eth := common.BytesToAddress(dex.Bytes())
	t.Logf("%s dex %s", dexAccount, eth.String())

	dexd := common.HexToAddress(ethAccount)
	dff := sdk.AccAddress(dexd[:20])
	t.Logf("%s dex %s", ethAccount, dff.String())

	feeAccountEth := "0xF0f4C5079BCf15a1f797326CE74aAC3375f5F693"
	feeAccountDex := "dex17xpfvakm2amg962yls6f84z3kell8c5l5s9l0c"

	dexd = common.HexToAddress(feeAccountEth)
	dff = sdk.AccAddress(dexd[:20])
	t.Logf("%s dex %s", feeAccountEth, dff.String())

	dex, _ = sdk.AccAddressFromBech32(feeAccountDex)
	eth = common.BytesToAddress(eth.Bytes())
	t.Logf("%s dex %s", feeAccountDex, eth.String())

	facc := "0x17c2bd128aaD7DD1f5b3dC31403528DcdF29863b"
	dexd = common.HexToAddress(facc)
	dff = sdk.AccAddress(dexd[:20])
	t.Logf("%s dex %s", facc, dff.String())

}
