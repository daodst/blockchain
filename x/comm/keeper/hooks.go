package keeper

import (
	"freemasonry.cc/blockchain/x/comm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) AfterDelegationModified(ctx sdk.Context, addr sdk.AccAddress, valAddr sdk.ValAddress) {
	//，
	if valAddr.String() == sdk.ValAddress(addr).String() {
		//
		gateway, err := k.GetGatewayInfo(ctx, valAddr.String())
		if err != nil {
			if err == types.ErrGatewayNotExist {
				return
			}
			panic(err)
		}
		//
		delegation, found := k.stakingKeeper.GetDelegation(ctx, addr, valAddr)
		if !found {
			return
		}
		params := k.GetParams(ctx)
		num := delegation.Shares.QuoInt(params.MinDelegate)
		gateway.GatewayQuota = num.TruncateInt64()
		//
		err = k.UpdateGatewayInfo(ctx, *gateway)
		if err != nil {
			panic(err)
		}
		//
		err = k.SetGatewayDelegateLastTime(ctx, addr.String(), valAddr.String())
		if err != nil {
			panic(err)
		}
	}
}

// Hooks wrapper struct for slashing keeper
type Hooks struct {
	k Keeper
}

var _ types.StakingHooks = Hooks{}

var _ types.CommonHooks = Keeper{}

func (k Keeper) AfterCreateGateway(ctx sdk.Context, validator stakingTypes.Validator) {
	if k.hooks != nil {
		k.hooks.AfterCreateGateway(ctx, validator)
	}
}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

//
func (h Hooks) AfterDelegationModified(ctx sdk.Context, addr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.k.AfterDelegationModified(ctx, addr, valAddr)
}

func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)          {}
func (h Hooks) AfterValidatorRemoved(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)         {}
func (h Hooks) AfterValidatorCreated(_ sdk.Context, _ sdk.ValAddress)                            {}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)  {}
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                          {}
func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}
func (h Hooks) BeforeValidatorSlashed(_ sdk.Context, _ sdk.ValAddress, _ sdk.Dec)                {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
