package client

import (
	"fmt"
	cmdcfg "freemasonry.cc/blockchain/cmd/config"
	"freemasonry.cc/blockchain/core"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"testing"
)

func TestDelegation(t *testing.T) {
	txClient := NewTxClient()
	dposClient := NewDposClinet(&txClient)
	am, _ := sdk.NewIntFromString("1000000000000000000000")
	result, err := dposClient.Delegation("dst1emqqtns9xdrpyjfant4wzf7zsylc59ydr0dthc", "dstvaloper1emqqtns9xdrpyjfant4wzf7zsylc59ydypaej5", sdk.NewCoin(core.GovDenom, am), "4F3BD44F642878FE2891A4D1B65941E9F10330497B6164C439DFBA841C39E276")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("data:", result.Data)
}

func TestUnbondDelegation(t *testing.T) {
	txClient := NewTxClient()
	dposClient := NewDposClinet(&txClient)
	result, err := dposClient.UnbondDelegation("dst1zl5quqaukt4ssks8j2nr6rvz472rl67d7ynmuk", "dstvaloper1zl5quqaukt4ssks8j2nr6rvz472rl67de2rfe6", sdk.NewCoin(core.GovDenom, sdk.NewInt(200)), "27632208e262fc0dd99d543c5aed8a5f2987ce3bd28287f68ff92e63ae7a27b8")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("data:", result.Data)
}

func TestRegisterValidate(t *testing.T) {

	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	cmdcfg.SetBip44CoinType(config)
	config.Seal()
	cmdcfg.RegisterDenoms()

	txClient := NewTxClient()
	blockClient := NewBlockClient()
	dposClient := NewDposClinet(&txClient)
	addr, err := sdk.AccAddressFromBech32("dst1kxcnfvtp6ep042vudctgda29cywkpzp0ttky75")
	if err != nil {
		t.Error(err)
		return
	}
	status, err := blockClient.StatusInfo()
	if err != nil {
		t.Error(err)
		return
	}

	valAddr := sdk.ValAddress(addr)
	
	selfDelegation := core.RealString2LedgerCoin("20", "nxn")
	fmt.Println(selfDelegation)

	
	minSelfDelegation := sdk.NewInt(1)
	description := stakingTypes.NewDescription("dao1", "", "", "", "")
	rate, _ := sdk.NewDecFromStr("0.99")
	maxRate, _ := sdk.NewDecFromStr("0.99")
	maxChangeRate, _ := sdk.NewDecFromStr("0.0001")
	commission := stakingTypes.NewCommissionRates(rate, maxRate, maxChangeRate) 
	pubkey, err := cryptocodec.FromTmPubKeyInterface(status.ValidatorInfo.PubKey)
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := dposClient.RegisterValidator(addr.String(), valAddr.String(), pubkey, selfDelegation, description, commission, minSelfDelegation, "349b21b102c799984d0a2d29986fee00f36abcaaa51f10bb196f44dd5a4614a6", 6)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}
