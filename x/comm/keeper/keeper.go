package keeper

import (
	"fmt"
	"freemasonry.cc/blockchain/core"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/libs/log"
	"time"

	"freemasonry.cc/blockchain/x/comm/types"
)

// Keeper of this module maintains collections of erc20.
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	stakingKeeper *stakingKeeper.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	hooks         types.CommonHooks
}

// NewKeeper creates new instances of the erc20 Keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	stakingKeeper *stakingKeeper.Keeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramstore:    ps,
		accountKeeper: ak,
		bankKeeper:    bk,
		stakingKeeper: stakingKeeper,
		hooks:         nil,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) KVHelper(ctx sdk.Context) StoreHelper {
	store := ctx.KVStore(k.storeKey)
	return StoreHelper{
		store,
	}
}

func (k *Keeper) SetHooks(sh types.CommonHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set common hooks twice")
	}
	k.hooks = sh
	return k
}

func (k Keeper) AddressBookSet(ctx sdk.Context, fromAddress string, AddressBook []string) error {
	store := k.KVHelper(ctx)
	key := types.KeyPrefixAddressBook + fromAddress
	err := store.Set(key, AddressBook)
	if err != nil {
		return err
	}
	return nil
}

//
func (k Keeper) SetGatewayDelegateLastTime(ctx sdk.Context, delegateAddress, validatorAddress string) error {
	store := k.KVHelper(ctx)
	key := types.DelegateLastTimeKey + delegateAddress + "_" + validatorAddress
	err := store.Set(key, ctx.BlockHeight())
	if err != nil {
		return err
	}
	return nil
}

//
func (k Keeper) GetGatewayDelegateLastTime(ctx sdk.Context, delegateAddress, validatorAddress string) (int64, error) {
	store := k.KVHelper(ctx)
	key := types.DelegateLastTimeKey + delegateAddress + "_" + validatorAddress
	var lastHeight int64
	if !store.Has(key) {
		return lastHeight, nil
	}
	err := store.GetUnmarshal(key, &lastHeight)
	if err != nil {
		return 0, err
	}
	return lastHeight, nil
}

//
func (k Keeper) GatewayBonus(ctx sdk.Context, params types.Params) error {
	//
	if (ctx.BlockHeight() % params.BonusCycle) == 0 {
		//3
		index := ctx.BlockHeight() / params.BonusHalve
		amount := sdk.NewCoin(sdk.DefaultBondDenom, params.Bonus.Quo(sdk.NewInt(1<<index)))
		if amount.IsZero() {
			return nil
		}
		return k.bankKeeper.SendCoins(ctx, core.ContractGatewayBonus, core.ContractAddressFee, sdk.NewCoins(amount))
	}
	return nil
}

func (k Keeper) RedeemCheck(ctx sdk.Context, params types.Params) error {
	if ctx.BlockHeight()%params.IndexNumHeight == 0 {
		redeemMap, err := k.GetGatewayRedeemNum(ctx)
		if err != nil {
			return err
		}
		var numArray []types.GatewayNumIndex
		for _, val := range redeemMap {
			//
			if val.Status == 1 && val.Validity <= ctx.BlockHeight() {
				val.GatewayAddress = ""
				val.Status = 2
				val.Validity = 0
				numArray = append(numArray, val)
			}
		}
		//
		err = k.SetGatewayNum(ctx, numArray)
		if err != nil {
			return err
		}
		//
		err = k.GatewayRedeemNumFilter(ctx, numArray)
		if err != nil {
			return err
		}
	}

	return nil
}

//
func (k Keeper) createValidator(ctx sdk.Context, delegatorAddress sdk.AccAddress, validatorAddress sdk.ValAddress, msg types.MsgCreateSmartValidator, delegation sdk.Coin) error {
	pk, err := ParseBech32ValConsPubkey(msg.PubKey)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pk)
	}

	if _, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return stakingTypes.ErrValidatorPubKeyExists
	}
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	if delegation.Denom != bondDenom {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", delegation.Denom, bondDenom,
		)
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return err
	}
	validator, err := stakingTypes.NewValidator(validatorAddress, pk, msg.Description)
	if err != nil {
		return err
	}
	commission := stakingTypes.NewCommissionWithTime(
		msg.Commission.Rate, msg.Commission.MaxRate,
		msg.Commission.MaxChangeRate, ctx.BlockHeader().Time,
	)
	validator, err = validator.SetInitialCommission(commission)
	//
	validator.MinSelfDelegation = msg.MinSelfDelegation
	k.stakingKeeper.SetValidator(ctx, validator)
	k.stakingKeeper.SetValidatorByConsAddr(ctx, validator)
	k.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)
	k.stakingKeeper.AfterValidatorCreated(ctx, validator.GetOperator())
	//
	_, err = k.stakingKeeper.Delegate(ctx, delegatorAddress, delegation.Amount, stakingTypes.Unbonded, validator, true)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingTypes.EventTypeCreateValidator,
			sdk.NewAttribute(stakingTypes.AttributeKeyValidator, validator.GetOperator().String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, delegation.String()),
		),
	})
	return nil
}

//dpos
func (k Keeper) delegate(ctx sdk.Context, delegatorAddress sdk.AccAddress, validatorAddress sdk.ValAddress, validator stakingTypes.Validator, coin sdk.Coin) error {
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	if coin.Denom != bondDenom {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", coin.Denom, bondDenom,
		)
	}
	//dpos
	newShares, err := k.stakingKeeper.Delegate(ctx, delegatorAddress, coin.Amount, stakingTypes.Unbonded, validator, true)
	if err != nil {
		return err
	}
	//
	err = k.SetGatewayDelegateLastTime(ctx, delegatorAddress.String(), validatorAddress.String())
	if err != nil {
		return err
	}
	coins := k.bankKeeper.GetAllBalances(ctx, delegatorAddress)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingTypes.EventTypeDelegate,
			sdk.NewAttribute(stakingTypes.AttributeKeyValidator, validatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, coin.String()),
			sdk.NewAttribute(stakingTypes.AttributeKeyNewShares, newShares.String()),    //
			sdk.NewAttribute(stakingTypes.AttributeKeyDelegatorBalance, coins.String()), //
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, stakingTypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, delegatorAddress.String()),
		),
	})
	return nil
}

//dpos,
func (k Keeper) Undelegate(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, validator stakingTypes.Validator, sharesAmount sdk.Dec) (time.Time, sdk.Int, error) {

	if k.stakingKeeper.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
		return time.Time{}, sdk.ZeroInt(), stakingTypes.ErrMaxUnbondingDelegationEntries
	}
	returnAmount, err := k.stakingKeeper.Unbond(ctx, delAddr, valAddr, sharesAmount)
	if err != nil {
		return time.Time{}, sdk.ZeroInt(), err
	}
	//，
	if validator.GetOperator().String() == sdk.ValAddress(delAddr).String() {
		params := k.GetParams(ctx)
		//
		lastTime, err := k.GetGatewayDelegateLastTime(ctx, delAddr.String(), valAddr.String())
		if err != nil {
			return time.Time{}, sdk.ZeroInt(), err
		}
		diff := ctx.BlockHeight() - lastTime
		if diff < params.RedeemFeeHeight { //10%
			fees := returnAmount.ToDec().Mul(params.RedeemFee)
			returnAmount = returnAmount.Sub(fees.RoundInt())
			coin := sdk.NewCoin(sdk.DefaultBondDenom, fees.RoundInt())
			err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, stakingTypes.BondedPoolName, authtypes.FeeCollectorName, sdk.NewCoins(coin))
			if err != nil {
				panic(err)
			}
		}
	}
	// transfer the validator tokens to the not bonded pool
	if validator.IsBonded() {
		k.bondedTokensToNotBonded(ctx, returnAmount)
	}

	completionTime := ctx.BlockHeader().Time.Add(k.stakingKeeper.UnbondingTime(ctx))
	ubd := k.stakingKeeper.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
	k.stakingKeeper.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, returnAmount, nil
}

func (k Keeper) bondedTokensToNotBonded(ctx sdk.Context, tokens sdk.Int) {
	coins := sdk.NewCoins(sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), tokens))
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, stakingTypes.BondedPoolName, stakingTypes.NotBondedPoolName, coins); err != nil {
		panic(err)
	}
}
