package types

import (
	"fmt"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	//
	DefaultIndexNumHeight = int64(100)
	//()
	DefaultRedeemFeeHeight = int64(432000)
	// 0.1
	DefaultRedeemFee = sdk.NewDec(1).Quo(sdk.NewDec(10))
	// 10000
	DefaultMinDelegate = sdk.NewInt(10000).Mul(sdk.NewInt(core.RealToLedgerRateInt64))
	//()
	DefaultValidity = int64(5256000)
	// ()
	DefaultBonusCycle = int64(14400)
	// 3()
	DefaultBonusHalve = int64(15768000)
	// 10000
	DefaultBonus = sdk.NewInt(10000).Mul(sdk.NewInt(core.RealToLedgerRateInt64))
)

var (
	KeyIndexNumHeight  = []byte("IndexNumHeight")
	KeyRedeemFeeHeight = []byte("RedeemFeeHeight")
	KeyRedeemFee       = []byte("RedeemFee")
	KeyMinDelegate     = []byte("MinDelegate")
	KeyValidity        = []byte("Validity")
	KeyBonusCycle      = []byte("BonusCycle")
	KeyBonusHalve      = []byte("BonusHalve")
	KeyBonus           = []byte("Bonus")
)

// NewParams creates a new Params object
func NewParams(
	IndexNumHeight int64,
	RedeemFeeHeight int64,
	RedeemFee sdk.Dec,
	MinDelegate sdk.Int,
	Validity int64,
	BonusCycle int64,
	BonusHalve int64,
	Bonus sdk.Int,
) Params {
	return Params{
		IndexNumHeight:  IndexNumHeight,
		RedeemFeeHeight: RedeemFeeHeight,
		RedeemFee:       RedeemFee,
		MinDelegate:     MinDelegate,
		Validity:        Validity,
		BonusCycle:      BonusCycle,
		BonusHalve:      BonusHalve,
		Bonus:           Bonus,
	}
}

func DefaultParams() Params {

	return Params{
		IndexNumHeight:  DefaultIndexNumHeight,
		RedeemFeeHeight: DefaultRedeemFeeHeight,
		RedeemFee:       DefaultRedeemFee,
		MinDelegate:     DefaultMinDelegate,
		Validity:        DefaultValidity,
		BonusCycle:      DefaultBonusCycle,
		BonusHalve:      DefaultBonusHalve,
		Bonus:           DefaultBonus,
	}
}

func (p Params) Validate() error {
	if err := validateIndexNumHeight(p.IndexNumHeight); err != nil {
		return err
	}

	if err := validateRedeemFeeHeight(p.RedeemFeeHeight); err != nil {
		return err
	}

	if err := validateRedeemFee(p.RedeemFee); err != nil {
		return err
	}

	if err := validateMinDelegate(p.MinDelegate); err != nil {
		return err
	}
	if err := validateValidity(p.Validity); err != nil {
		return err
	}
	if err := validateBonusCycle(p.BonusCycle); err != nil {
		return err
	}
	if err := validateBonusHalve(p.BonusHalve); err != nil {
		return err
	}
	if err := validateBonus(p.Bonus); err != nil {
		return err
	}
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyIndexNumHeight, &p.IndexNumHeight, validateIndexNumHeight),
		paramtypes.NewParamSetPair(KeyRedeemFeeHeight, &p.RedeemFeeHeight, validateRedeemFeeHeight),
		paramtypes.NewParamSetPair(KeyRedeemFee, &p.RedeemFee, validateRedeemFee),
		paramtypes.NewParamSetPair(KeyMinDelegate, &p.MinDelegate, validateMinDelegate),
		paramtypes.NewParamSetPair(KeyValidity, &p.Validity, validateValidity),
		paramtypes.NewParamSetPair(KeyBonusCycle, &p.BonusCycle, validateBonusCycle),
		paramtypes.NewParamSetPair(KeyBonusHalve, &p.BonusHalve, validateBonusHalve),
		paramtypes.NewParamSetPair(KeyBonus, &p.Bonus, validateBonus),
	}
}

func validateIndexNumHeight(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("IndexNumHeight must be positive: %d", v)
	}
	return nil
}

func validateRedeemFeeHeight(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("RedeemFeeHeight must be positive: %d", v)
	}
	return nil
}

func validateRedeemFee(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("RedeemFee cannot be negative: %s", v)
	}

	return nil
}

func validateMinDelegate(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("MinDelegate cannot be negative: %s", v)
	}
	return nil
}

func validateValidity(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("Validity must be positive: %d", v)
	}
	return nil
}

func validateBonusCycle(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("BonusCycle must be positive: %d", v)
	}
	return nil
}

func validateBonusHalve(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("BonusHalve must be positive: %d", v)
	}
	return nil
}

func validateBonus(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("Bonus cannot be negative: %s", v)
	}
	return nil
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(KeyIndexNumHeight, DefaultParams().IndexNumHeight, validateIndexNumHeight),
		paramtypes.NewParamSetPair(KeyRedeemFeeHeight, DefaultParams().RedeemFeeHeight, validateRedeemFeeHeight),
		paramtypes.NewParamSetPair(KeyRedeemFee, DefaultParams().RedeemFee, validateRedeemFee),
		paramtypes.NewParamSetPair(KeyMinDelegate, DefaultParams().MinDelegate, validateMinDelegate),
		paramtypes.NewParamSetPair(KeyValidity, DefaultParams().Validity, validateValidity),
		paramtypes.NewParamSetPair(KeyBonusCycle, DefaultParams().BonusCycle, validateBonusCycle),
		paramtypes.NewParamSetPair(KeyBonusHalve, DefaultParams().BonusHalve, validateBonusHalve),
		paramtypes.NewParamSetPair(KeyBonus, DefaultParams().Bonus, validateBonus),
	)
}
