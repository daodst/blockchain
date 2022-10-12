package keeper

import (
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/pledge/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	types2 "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/libs/log"
	tmstrings "github.com/tendermint/tendermint/libs/strings"
	"time"
)

// Keeper of this module maintains collections of erc20.
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	hooks            types.PledgeHooks
	AccountKeeper    types.AccountKeeper
	BankKeeper       types.BankKeeper
	CommKeeper       types.CommKeeper
	FeeCollectorName string
}

// NewKeeper creates new instances of the erc20 Keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	ck types.CommKeeper,
	feeCollectorName string,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	// ensure bonded and not bonded module accounts are set
	if addr := ak.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	return Keeper{
		storeKey:         storeKey,
		cdc:              cdc,
		paramstore:       ps,
		AccountKeeper:    ak,
		BankKeeper:       bk,
		hooks:            nil,
		CommKeeper:       ck,
		FeeCollectorName: feeCollectorName,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) KVHelper(ctx sdk.Context) storeHelper {
	store := ctx.KVStore(k.storeKey)
	return storeHelper{
		store,
	}
}

func (k Keeper) BondDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyBondDenom, &res)
	return
}

func (k Keeper) GetAllChatPledge(ctx sdk.Context, delAddr sdk.AccAddress) (types.DelegationResponses, error) {
	delegations := k.GetAllDelegatorDelegations(ctx, delAddr)
	delegationResps, err := DelegationsToDelegationResponses(ctx, k, delegations)
	if err != nil {
		return types.DelegationResponses{}, err
	}
	return delegationResps, nil
}

// Delegate performs a delegation, set/update everything necessary within the store.
// tokenSrc indicates the bond status of the incoming funds.
func (k Keeper) Delegate(
	ctx sdk.Context, fromAddress, delAddr sdk.AccAddress, bondAmt sdk.Int,
	validator types.Validator,
) (newShares sdk.Dec, err error) {

	log := core.BuildLog(core.GetFuncName(), core.LmChainPledgeKeeper)

	//
	params := k.GetParams(ctx)
	delegations := k.GetAllDelegatorDelegations(ctx, delAddr)
	delegationResps, err := DelegationsToDelegationResponses(ctx, k, delegations)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	pledgeTotal := sdk.ZeroInt()
	if len(delegationResps) > 0 {
		for _, resp := range delegationResps {
			pledgeTotal = pledgeTotal.Add(resp.Balance.Amount)
		}
	}

	log.Debug("+++++++++++++++++++++++++")
	log.Debug(bondAmt)
	log.Debug(pledgeTotal)
	log.Debug(params.MinMortgageCoin.Amount)
	log.Debug("--------------------------")
	if pledgeTotal.Add(bondAmt).LT(params.MinMortgageCoin.Amount) {
		return sdk.ZeroDec(), errors.New("pledge total amount less than min pledge amount")
	}

	// In some situations, the exchange rate becomes invalid, e.g. if
	// Validator loses all tokens due to slashing. In this case,
	// make all future delegations invalid.

	if validator.InvalidExRate() {
		return sdk.ZeroDec(), types.ErrDelegatorShareExRateInvalid
	}

	// Get or create the delegation object
	delegation, found := k.GetDelegation(ctx, delAddr, validator.GetOperator())
	if !found {
		delegation = types.NewDelegation(delAddr, validator.GetOperator(), sdk.ZeroDec())
	}

	// call the appropriate hook if present
	if found {
		//k.BeforeDelegationSharesModified(ctx, delAddr, validator.GetOperator())

		val := k.Validator(ctx, validator.GetOperator())
		del := k.Delegation(ctx, delAddr, validator.GetOperator())
		_, err = k.withdrawDelegationRewards(ctx, val, del)
		if err != nil {
			log.WithError(err).Error("withdrawDelegationRewards")
			return sdk.ZeroDec(), err
		}

	} else {
		//k.BeforeDelegationCreated(ctx, delAddr, validator.GetOperator())
		val := k.Validator(ctx, validator.GetOperator())
		k.IncrementValidatorPeriod(ctx, val)
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return sdk.ZeroDec(), err
	}

	sendName := types.BondedPoolName

	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), bondAmt))
	if err := k.BankKeeper.DelegateCoinsFromAccountToModuleForPledge(ctx, fromAddress, delegatorAddress, sendName, coins); err != nil {
		log.WithError(err).Error("DelegateCoinsFromAccountToModuleForPledge")
		return sdk.Dec{}, err
	}

	_, newShares = k.AddValidatorTokensAndShares(ctx, validator, bondAmt)

	// Update delegation
	delegation.Shares = delegation.Shares.Add(newShares)
	k.SetDelegation(ctx, delegation)

	//
	err = k.SetPledgeSum(ctx, delegation.DelegatorAddress, bondAmt)
	if err != nil {
		log.WithError(err).Error("SetPledgeSum")
		return sdk.Dec{}, err
	}

	// Call the after-modification hook

	//todo hook？
	//k.AfterDelegationModified(ctx, delegatorAddress, delegation.GetValidatorAddr())
	k.initializeDelegation(ctx, validator.GetOperator(), delAddr)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDelegate,
			sdk.NewAttribute(types.EventPledgeFromAddress, fromAddress.String()),
			sdk.NewAttribute(types.EventPledgeToAddress, validator.OperatorAddress),
			sdk.NewAttribute(types.EventPledgeAmount, bondAmt.String()),
			sdk.NewAttribute(types.EventPledgeDenom, k.BondDenom(ctx)),
			sdk.NewAttribute(types.EventPledgeFromBalance, k.BankKeeper.GetBalance(ctx, fromAddress, k.BondDenom(ctx)).Amount.String()),
		),
	)

	return newShares, nil
}

func (k Keeper) SetPledgeSum(ctx sdk.Context, delegatorAddress string, amount sdk.Int) error {
	//logs := core.BuildLog(core.GetFuncName(), core.LmChainPledgeKeeper)
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.PledgeDelegateSumKey + delegatorAddress)

	if store.Has(key) {
		oldAmountByte := store.Get(key)
		oldAmount, ok := sdk.NewIntFromString(string(oldAmountByte))
		if !ok {
			return errors.New("error set pledge sum amount")
		}
		newAmount := oldAmount.Add(amount)
		store.Set(key, []byte(newAmount.String()))
	} else {
		store.Set(key, []byte(amount.String()))
	}

	return nil
}

func (k Keeper) GetPledgeSum(ctx sdk.Context, delegatorAddress string) (getAmount sdk.Int, err error) {
	//logs := core.BuildLog(core.GetFuncName(), core.LmChainPledgeKeeper)
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.PledgeDelegateSumKey + delegatorAddress)

	if store.Has(key) {
		AmountByte := store.Get(key)
		getAmount, ok := sdk.NewIntFromString(string(AmountByte))
		if !ok {
			return getAmount, errors.New("error get pledge sum amount")
		}
		return getAmount, nil
	} else {
		return sdk.ZeroInt(), nil
	}
}

// GetDelegation returns a specific delegation.
func (k Keeper) GetDelegation(ctx sdk.Context,
	delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation types.Delegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDelegationKey(delAddr, valAddr)

	value := store.Get(key)
	if value == nil {
		return delegation, false
	}

	delegation = types.MustUnmarshalDelegation(k.cdc, value)

	return delegation, true
}

func (k Keeper) SetDelegation(ctx sdk.Context, delegation types.Delegation) {
	delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		panic(err)
	}

	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(types.GetDelegationKey(delegatorAddress, delegation.GetValidatorAddr()), b)
}

// DequeueAllMatureUBDQueue returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue.
func (k Keeper) DequeueAllMatureUBDQueue(ctx sdk.Context, currTime time.Time) (matureUnbonds []types.DVPair) {
	store := ctx.KVStore(k.storeKey)

	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	unbondingTimesliceIterator := k.UBDQueueIterator(ctx, ctx.BlockHeader().Time)
	defer unbondingTimesliceIterator.Close()

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := types.DVPairs{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshal(value, &timeslice)

		matureUnbonds = append(matureUnbonds, timeslice.Pairs...)

		store.Delete(unbondingTimesliceIterator.Key())
	}

	return matureUnbonds
}

func (k Keeper) UBDQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UnbondingQueueKey,
		sdk.InclusiveEndBytes(types.GetUnbondingDelegationTimeKey(endTime)))
}

// CompleteUnbonding completes the unbonding of all mature entries in the
// retrieved unbonding delegation object and returns the total unbonding balance
// or an error upon failure.
func (k Keeper) CompleteUnbonding(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	ubd, found := k.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, types.ErrNoUnbondingDelegation
	}

	bondDenom := k.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time

	delegatorAddress, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(ctxTime) {
			ubd.RemoveEntry(int64(i))
			i--

			// track undelegation only when remaining or truncated shares are non-zero
			if !entry.Balance.IsZero() {
				amt := sdk.NewCoin(bondDenom, entry.Balance)
				if err := k.BankKeeper.UndelegateCoinsFromModuleToAccount(
					ctx, types.NotBondedPoolName, delegatorAddress, sdk.NewCoins(amt),
				); err != nil {
					return nil, err
				}

				balances = balances.Add(amt)
			}
		}
	}

	// set the unbonding delegation or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveUnbondingDelegation(ctx, ubd)
	} else {
		k.SetUnbondingDelegation(ctx, ubd)
	}

	return balances, nil
}

// GetUnbondingDelegation returns a unbonding delegation.
func (k Keeper) GetUnbondingDelegation(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress,
) (ubd types.UnbondingDelegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDKey(delAddr, valAddr)
	value := store.Get(key)

	if value == nil {
		return ubd, false
	}

	ubd = types.MustUnmarshalUBD(k.cdc, value)

	return ubd, true
}

// RemoveUnbondingDelegation removes the unbonding delegation object and associated index.
func (k Keeper) RemoveUnbondingDelegation(ctx sdk.Context, ubd types.UnbondingDelegation) {
	delegatorAddress, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		panic(err)
	}

	store := ctx.KVStore(k.storeKey)
	addr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	key := types.GetUBDKey(delegatorAddress, addr)
	store.Delete(key)
	store.Delete(types.GetUBDByValIndexKey(delegatorAddress, addr))
}

// SetUnbondingDelegationEntry adds an entry to the unbonding delegation at the given addresses. It creates the unbonding delegation if it does not exist.
func (k Keeper) SetUnbondingDelegationEntry(
	ctx sdk.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	creationHeight int64, minTime time.Time, balance sdk.Int,
) types.UnbondingDelegation {
	ubd, found := k.GetUnbondingDelegation(ctx, delegatorAddr, validatorAddr)
	if found {
		ubd.AddEntry(creationHeight, minTime, balance)
	} else {
		ubd = types.NewUnbondingDelegation(delegatorAddr, validatorAddr, creationHeight, minTime, balance)
	}

	k.SetUnbondingDelegation(ctx, ubd)

	return ubd
}

// SetUnbondingDelegation sets the unbonding delegation and associated index.
func (k Keeper) SetUnbondingDelegation(ctx sdk.Context, ubd types.UnbondingDelegation) {
	delegatorAddress, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalUBD(k.cdc, ubd)
	addr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	key := types.GetUBDKey(delegatorAddress, addr)
	store.Set(key, bz)
	store.Set(types.GetUBDByValIndexKey(delegatorAddress, addr), []byte{}) // index, store empty bytes
}

// DequeueAllMatureRedelegationQueue returns a concatenated list of all the
// timeslices inclusively previous to currTime, and deletes the timeslices from
// the queue.
func (k Keeper) DequeueAllMatureRedelegationQueue(ctx sdk.Context, currTime time.Time) (matureRedelegations []types.DVVTriplet) {
	store := ctx.KVStore(k.storeKey)

	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	redelegationTimesliceIterator := k.RedelegationQueueIterator(ctx, ctx.BlockHeader().Time)
	defer redelegationTimesliceIterator.Close()

	for ; redelegationTimesliceIterator.Valid(); redelegationTimesliceIterator.Next() {
		timeslice := types.DVVTriplets{}
		value := redelegationTimesliceIterator.Value()
		k.cdc.MustUnmarshal(value, &timeslice)

		matureRedelegations = append(matureRedelegations, timeslice.Triplets...)

		store.Delete(redelegationTimesliceIterator.Key())
	}

	return matureRedelegations
}

// RedelegationQueueIterator returns all the redelegation queue timeslices from
// time 0 until endTime.
func (k Keeper) RedelegationQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.RedelegationQueueKey, sdk.InclusiveEndBytes(types.GetRedelegationTimeKey(endTime)))
}

// CompleteRedelegation completes the redelegations of all mature entries in the
// retrieved redelegation object and returns the total redelegation (initial)
// balance or an error upon failure.
func (k Keeper) CompleteRedelegation(
	ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress,
) (sdk.Coins, error) {
	red, found := k.GetRedelegation(ctx, delAddr, valSrcAddr, valDstAddr)
	if !found {
		return nil, types.ErrNoRedelegation
	}

	bondDenom := k.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time

	// loop through all the entries and complete mature redelegation entries
	for i := 0; i < len(red.Entries); i++ {
		entry := red.Entries[i]
		if entry.IsMature(ctxTime) {
			red.RemoveEntry(int64(i))
			i--

			if !entry.InitialBalance.IsZero() {
				balances = balances.Add(sdk.NewCoin(bondDenom, entry.InitialBalance))
			}
		}
	}

	// set the redelegation or remove it if there are no more entries
	if len(red.Entries) == 0 {
		k.RemoveRedelegation(ctx, red)
	} else {
		k.SetRedelegation(ctx, red)
	}

	return balances, nil
}

// GetRedelegation returns a redelegation.
func (k Keeper) GetRedelegation(ctx sdk.Context,
	delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress) (red types.Redelegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetREDKey(delAddr, valSrcAddr, valDstAddr)

	value := store.Get(key)
	if value == nil {
		return red, false
	}

	red = types.MustUnmarshalRED(k.cdc, value)

	return red, true
}

// RemoveRedelegation removes a redelegation object and associated index.
func (k Keeper) RemoveRedelegation(ctx sdk.Context, red types.Redelegation) {
	delegatorAddress, err := sdk.AccAddressFromBech32(red.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	store := ctx.KVStore(k.storeKey)
	valSrcAddr, err := sdk.ValAddressFromBech32(red.ValidatorSrcAddress)
	if err != nil {
		panic(err)
	}
	valDestAddr, err := sdk.ValAddressFromBech32(red.ValidatorDstAddress)
	if err != nil {
		panic(err)
	}
	redKey := types.GetREDKey(delegatorAddress, valSrcAddr, valDestAddr)
	store.Delete(redKey)
	store.Delete(types.GetREDByValSrcIndexKey(delegatorAddress, valSrcAddr, valDestAddr))
	store.Delete(types.GetREDByValDstIndexKey(delegatorAddress, valSrcAddr, valDestAddr))
}

// SetRedelegation set a redelegation and associated index.
func (k Keeper) SetRedelegation(ctx sdk.Context, red types.Redelegation) {
	delegatorAddress, err := sdk.AccAddressFromBech32(red.DelegatorAddress)
	if err != nil {
		panic(err)
	}

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalRED(k.cdc, red)
	valSrcAddr, err := sdk.ValAddressFromBech32(red.ValidatorSrcAddress)
	if err != nil {
		panic(err)
	}
	valDestAddr, err := sdk.ValAddressFromBech32(red.ValidatorDstAddress)
	if err != nil {
		panic(err)
	}
	key := types.GetREDKey(delegatorAddress, valSrcAddr, valDestAddr)
	store.Set(key, bz)
	store.Set(types.GetREDByValSrcIndexKey(delegatorAddress, valSrcAddr, valDestAddr), []byte{})
	store.Set(types.GetREDByValDstIndexKey(delegatorAddress, valSrcAddr, valDestAddr), []byte{})
}

// Set the last total validator power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{Int: power})
	store.Set(types.LastTotalPowerKey, bz)
}

func (k Keeper) WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	logs := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	val := k.Validator(ctx, valAddr)
	if val == nil {
		return nil, types.ErrNoValidatorDistInfo
	}

	del := k.Delegation(ctx, delAddr, valAddr)
	if del == nil {
		return nil, types.ErrEmptyDelegationDistInfo
	}

	// withdraw rewards
	rewards, err := k.withdrawDelegationRewards(ctx, val, del)
	if err != nil {
		return nil, err
	}

	logs.Info("")
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.ChatWithDrawEventType,
			sdk.NewAttribute(types.ChatWithDrawEventTypeFromAddress, delAddr.String()),
			sdk.NewAttribute(types.ChatWithDrawEventTypeModuleAddress, k.AccountKeeper.GetModuleAddress(types.ModuleName).String()),
			sdk.NewAttribute(types.ChatWithDrawEventTypeFromBalance, k.BankKeeper.GetAllBalances(ctx, delAddr).String()),
			sdk.NewAttribute(types.ChatWithDrawEventTypeAmount, rewards.AmountOf(core.BaseDenom).String()),
			sdk.NewAttribute(types.ChatWithDrawEventTypeDenom, core.BaseDenom),
		),
	)

	// reinitialize the delegation
	k.initializeDelegation(ctx, valAddr, delAddr)
	return rewards, nil
}

func (k Keeper) CreateValidator(ctx sdk.Context, smsg types2.MsgCreateValidator) error {
	//logs := core.BuildLog(core.GetFuncName(), core.LmChainPledgeKeeper)

	msg := types.MsgCreateValidator{}
	msg.Description = types.Description(smsg.Description)
	msg.Commission = types.CommissionRates(smsg.Commission)
	msg.MinSelfDelegation = smsg.MinSelfDelegation
	msg.DelegatorAddress = smsg.DelegatorAddress
	msg.ValidatorAddress = smsg.ValidatorAddress
	msg.Pubkey = smsg.Pubkey
	msg.Value = smsg.Value

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return err
	}

	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetValidator(ctx, valAddr); found {
		return types.ErrValidatorOwnerExists
	}

	pk, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pk)
	}

	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return types.ErrValidatorPubKeyExists
	}

	bondDenom := k.BondDenom(ctx)
	if msg.Value.Denom != bondDenom {
		msg.Value.Denom = bondDenom
	}

	//if _, err := msg.Description.EnsureLength(); err != nil {
	//	return err
	//}

	cp := ctx.ConsensusParams()
	if cp != nil && cp.Validator != nil {
		if !tmstrings.StringInSlice(pk.Type(), cp.Validator.PubKeyTypes) {
			return sdkerrors.Wrapf(
				types.ErrValidatorPubKeyTypeNotSupported,
				"got: %s, expected: %s", pk.Type(), cp.Validator.PubKeyTypes,
			)
		}
	}

	validator, err := types.NewValidator(valAddr, pk, msg.Description)
	if err != nil {
		return err
	}
	//commission := types.NewCommissionWithTime(
	//	msg.Commission.Rate, msg.Commission.MaxRate,
	//	msg.Commission.MaxChangeRate, ctx.BlockHeader().Time,
	//)
	//
	//validator, err = validator.SetInitialCommission(commission)
	//if err != nil {
	//	return err
	//}

	//delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	//if err != nil {
	//	return err
	//}

	//validator.MinSelfDelegation = msg.MinSelfDelegation

	k.SetValidator(ctx, validator)
	//k.SetValidatorByConsAddr(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)

	// call the after-creation hook
	//todo hook
	//k.AfterValidatorCreated(ctx, validator.GetOperator())

	k.InitializeValidator(ctx, validator)

	//ivalidator := k.Validator(ctx, valAddr)
	//consPk, err := ivalidator.ConsPubKey()
	//if err != nil {
	//	return err
	//}
	//k.AddPubkey(ctx, consPk)

	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	// NOTE source will always be from a wallet which are unbonded
	//_, err = k.Keeper.Delegate(ctx, delegatorAddress, msg.Value.Amount, types.Unbonded, validator, true)
	//if err != nil {
	//	return err
	//}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateValidator,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Value.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	})

	return nil
}

func (k Keeper) AddPubkey(ctx sdk.Context, pubkey cryptotypes.PubKey) error {
	bz, err := k.cdc.MarshalInterface(pubkey)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	key := types.AddrPubkeyRelationKey(pubkey.Address())
	store.Set(key, bz)
	return nil
}

//
func (k Keeper) SetPledgeDelegate(ctx sdk.Context, delegate map[string]types.PledgeDelegate) error {
	storeHelper := k.KVHelper(ctx)
	err := storeHelper.Set(types.PledgeDelegateKey, delegate)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) GetPledgeDelegate(ctx sdk.Context, address string) (sdk.Dec, error) {
	storeHelper := k.KVHelper(ctx)
	delDec := sdk.ZeroDec()
	if !storeHelper.Has(types.PledgeDelegateKey) {
		return delDec, nil
	}
	delegateMap := make(map[string]types.PledgeDelegate)
	err := storeHelper.GetUnmarshal(types.PledgeDelegateKey, &delegateMap)
	if err != nil {
		return delDec, err
	}

	if _, ok := delegateMap[address]; ok {
		amount := delegateMap[address].Amount
		amInt := core.MustRealString2LedgerIntNoMin(amount)
		return amInt.ToDec(), nil
	}
	return delDec, nil
}

// ChatPoundage ，（coin）
func (k Keeper) ChatPoundage(ctx sdk.Context, accFromAddress, nodeAddress sdk.AccAddress, distCoin sdk.Coin) (sdk.Dec, error) {

	//Dec
	AllCoinDec := distCoin.Amount.ToDec()

	//
	chatParams := k.GetParams(ctx)

	//
	destroyDec := AllCoinDec.Mul(chatParams.AttDestroyPercent)
	//
	gatewayDec := AllCoinDec.Mul(chatParams.AttGatewayPercent)
	//dpos
	DPosDec := AllCoinDec.Mul(chatParams.AttDposPercent)
	//
	feeAllInt := destroyDec.TruncateInt().Add(gatewayDec.TruncateInt()).Add(DPosDec.TruncateInt())
	feeAllDec := feeAllInt.ToDec()

	feeCoin := sdk.NewCoin(distCoin.Denom, feeAllInt)

	// 
	err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, accFromAddress, types.ModuleName, sdk.NewCoins(feeCoin))
	if err != nil {
		return sdk.Dec{}, types.ErrPledgeFeeTransfer
	}

	//
	err = k.BankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(distCoin.Denom, destroyDec.TruncateInt())))
	if err != nil {
		return sdk.Dec{}, types.ErrPledgeFeeBurn
	}
	//
	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, nodeAddress, sdk.NewCoins(sdk.NewCoin(distCoin.Denom, gatewayDec.TruncateInt())))
	if err != nil {
		return sdk.Dec{}, types.ErrPledgeFeeGateway
	}
	//dpos
	err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, sdk.NewCoins(sdk.NewCoin(distCoin.Denom, DPosDec.TruncateInt())))
	if err != nil {
		return sdk.Dec{}, types.ErrPledgeFeeDpos
	}

	return feeAllDec, nil
}
