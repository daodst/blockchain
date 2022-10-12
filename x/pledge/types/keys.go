package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

// constants
const (
	// module name
	ModuleName = "pledge"
	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName

	QuerierRoute = ModuleName

	FeeCollectorName = "pledge_fee_collector"

	PledgeDelegateKey = "pledge_delegate_key"

	PledgeDelegateSumKey = "pledge_delegate_sum_key_"
)

// ModuleAddress is the native module address for EVM
var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

// prefix bytes for the chat persistent store
const ()

//mint types

var (
	MinterKey                            = []byte{0x00}
	ProposerKey                          = []byte{0x01} // key for the proposer operator address
	FeePoolKey                           = []byte{0x02} // key for global distribution state
	DelegatorWithdrawAddrPrefix          = []byte{0x03} // key for delegator withdraw address
	DelegatorStartingInfoPrefix          = []byte{0x04} // key for delegator starting info
	ValidatorHistoricalRewardsPrefix     = []byte{0x05} // key for historical validators rewards / stake
	ValidatorCurrentRewardsPrefix        = []byte{0x06} // key for current validator rewards
	ValidatorAccumulatedCommissionPrefix = []byte{0x07} // key for accumulated validator commission
	ValidatorSlashEventPrefix            = []byte{0x08} // key for validator slash fraction
	ValidatorOutstandingRewardsPrefix    = []byte{0x09} // key for outstanding rewards
	AddrPubkeyRelationKeyPrefix          = []byte{0x10} // Prefix for address-pubkey relation
)

// GetDelegatorWithdrawAddrKey creates the key for a delegator's withdraw addr.
func GetDelegatorWithdrawAddrKey(delAddr sdk.AccAddress) []byte {
	return append(DelegatorWithdrawAddrPrefix, address.MustLengthPrefix(delAddr.Bytes())...)
}

// GetDelegatorStartingInfoKey creates the key for a delegator's starting info.
func GetDelegatorStartingInfoKey(v sdk.ValAddress, d sdk.AccAddress) []byte {
	return append(append(DelegatorStartingInfoPrefix, address.MustLengthPrefix(v.Bytes())...), address.MustLengthPrefix(d.Bytes())...)
}

// GetValidatorHistoricalRewardsKey creates the key for a validator's historical rewards.
func GetValidatorHistoricalRewardsKey(v sdk.ValAddress, k uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, k)
	return append(append(ValidatorHistoricalRewardsPrefix, address.MustLengthPrefix(v.Bytes())...), b...)
}

// GetValidatorHistoricalRewardsPrefix creates the prefix key for a validator's historical rewards.
func GetValidatorHistoricalRewardsPrefix(v sdk.ValAddress) []byte {
	return append(ValidatorHistoricalRewardsPrefix, address.MustLengthPrefix(v.Bytes())...)
}

// GetValidatorCurrentRewardsKey creates the key for a validator's current rewards.
func GetValidatorCurrentRewardsKey(v sdk.ValAddress) []byte {
	return append(ValidatorCurrentRewardsPrefix, address.MustLengthPrefix(v.Bytes())...)
}

// GetValidatorCurrentRewardsAddress creates the address from a validator's current rewards key.
func GetValidatorCurrentRewardsAddress(key []byte) (valAddr sdk.ValAddress) {
	// key is in the format:
	// 0x06<valAddrLen (1 Byte)><valAddr_Bytes>: ValidatorCurrentRewards

	// Remove prefix and address length.
	addr := key[2:]
	if len(addr) != int(key[1]) {
		panic("unexpected key length")
	}

	return sdk.ValAddress(addr)
}

// gets the key for a validator's current commission
func GetValidatorAccumulatedCommissionKey(v sdk.ValAddress) []byte {
	return append(ValidatorAccumulatedCommissionPrefix, v.Bytes()...)
}

// GetValidatorAccumulatedCommissionAddress creates the address from a validator's accumulated commission key.
func GetValidatorAccumulatedCommissionAddress(key []byte) (valAddr sdk.ValAddress) {
	// key is in the format:
	// 0x07<valAddrLen (1 Byte)><valAddr_Bytes>: ValidatorCurrentRewards

	// Remove prefix and address length.
	addr := key[2:]
	if len(addr) != int(key[1]) {
		panic("unexpected key length")
	}

	return sdk.ValAddress(addr)
}

// gets the outstanding rewards key for a validator
func GetValidatorOutstandingRewardsKey(valAddr sdk.ValAddress) []byte {
	return append(ValidatorOutstandingRewardsPrefix, valAddr.Bytes()...)
}

// GetValidatorOutstandingRewardsAddress creates an address from a validator's outstanding rewards key.
func GetValidatorOutstandingRewardsAddress(key []byte) (valAddr sdk.ValAddress) {
	// key is in the format:
	// 0x02<valAddrLen (1 Byte)><valAddr_Bytes>

	// Remove prefix and address length.
	addr := key[2:]
	if len(addr) != int(key[1]) {
		panic("unexpected key length")
	}

	return sdk.ValAddress(addr)
}

// gets the prefix key for a validator's slash fraction (ValidatorSlashEventPrefix + height)
func GetValidatorSlashEventKeyPrefix(v sdk.ValAddress, height uint64) []byte {
	heightBz := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBz, height)
	return append(
		ValidatorSlashEventPrefix,
		append(v.Bytes(), heightBz...)...,
	)
}

// gets the key for a validator's slash fraction
func GetValidatorSlashEventKey(v sdk.ValAddress, height, period uint64) []byte {
	periodBz := make([]byte, 8)
	binary.BigEndian.PutUint64(periodBz, period)
	prefix := GetValidatorSlashEventKeyPrefix(v, height)
	return append(prefix, periodBz...)
}

// GetValidatorSlashEventAddressHeight creates the height from a validator's slash event key.
func GetValidatorSlashEventAddressHeight(key []byte) (valAddr sdk.ValAddress, height uint64) {
	// key is in the format:
	// 0x08<valAddrLen (1 Byte)><valAddr_Bytes><height>: ValidatorSlashEvent
	valAddrLen := int(key[1])
	valAddr = key[2 : 2+valAddrLen]
	startB := 2 + valAddrLen
	b := key[startB : startB+8] // the next 8 bytes represent the height
	height = binary.BigEndian.Uint64(b)
	return
}
func AddrPubkeyRelationKey(addr []byte) []byte {
	return append(AddrPubkeyRelationKeyPrefix, address.MustLengthPrefix(addr)...)
}
