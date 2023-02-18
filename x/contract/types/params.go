package types

import (
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	//NFT 100
	DefaultDays = int64(8640000)
)

var (
	KeyDays = []byte("Days")
)

// NewParams creates a new Params object
func NewParams(
	days int64,
) Params {
	return Params{
		Days: days,
	}
}

func DefaultParams() Params {

	return Params{
		Days: DefaultDays,
	}
}

func (p Params) Validate() error {
	if err := validateDays(p.Days); err != nil {
		return err
	}
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDays, &p.Days, validateDays),
	}
}

func validateDays(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("Days must be positive: %d", v)
	}
	return nil
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(KeyDays, DefaultParams().Days, validateDays),
	)
}
