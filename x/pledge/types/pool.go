package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// names used as root for pool module accounts:
//
// - NotBondedPool -> "not_bonded_tokens_pool"
//
// - BondedPool -> "bonded_tokens_pool"
const (
	NotBondedPoolName = "pledge_not_bonded_tokens_pool"
	BondedPoolName    = "pledge_bonded_tokens_pool"
)

// NewPool creates a new Pool instance used for queries
func NewPool(notBonded, bonded sdk.Int) Pool {
	return Pool{
		NotBondedTokens: notBonded,
		BondedTokens:    bonded,
	}
}

//
//type Pool struct {
//	NotBondedTokens sdk.Int `json:"not_bonded_tokens"`
//	BondedTokens    sdk.Int `json:"bonded_tokens"`
//}
