package core

import (
	cmdcfg "freemasonry.cc/blockchain/cmd/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

func TestAccount(t *testing.T) {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	cmdcfg.SetBip44CoinType(config)

	t.Log(":", ContractAddressFee.String()) //dex17xpfvakm2amg962yls6f84z3kell8c5l5s9l0c

	t.Log(":", ContractAddressBank.String()) //dex1gwqac243g2z3vryqsev6acq965f9ttwhajlh6d

	t.Log(":", ContractAddressDistribution.String()) //dex1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8rkzrd6

	t.Log(":", ContractAddressStakingBonded.String()) //dex1fl48vsnmsdzcv85q5d2q4z5ajdha8yu33j0saj

	t.Log(":", ContractAddressStakingNotBonded.String()) //dex1tygms3xhhs3yv487phx3dw4a95jn7t7l9jnptx

	t.Log("IBC:", ContractAddressIbcTransfer.String()) //dex1yl6hdjhmkf37639730gffanpzndzdpmh2ksknx

	t.Log(":", ContractGatewayBonus.String()) //dex1rq7kxcflutn2tq8cx9tvlvr6clcl4slkskt92r

}
