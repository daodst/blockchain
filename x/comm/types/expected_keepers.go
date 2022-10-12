package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingHooks event hooks for staking validator object (noalias)
type StakingHooks interface {
	AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress)
}

type CommonHooks interface {
	AfterCreateGateway(ctx sdk.Context, validator stakingTypes.Validator)
}
