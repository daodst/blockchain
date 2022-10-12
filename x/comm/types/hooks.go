package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type MultiCommonHooks []CommonHooks

func NewMultiStakingHooks(hooks ...CommonHooks) MultiCommonHooks {
	return hooks
}

func (h MultiCommonHooks) AfterCreateGateway(ctx sdk.Context, validator stakingTypes.Validator) {
	for i := range h {
		h[i].AfterCreateGateway(ctx, validator)
	}
}
