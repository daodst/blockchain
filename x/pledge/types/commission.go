package types

import (
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"time"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewCommissionRates returns an initialized validator commission rates.
func NewCommissionRates(rate, maxRate, maxChangeRate sdk.Dec) CommissionRates {
	return CommissionRates{
		Rate:          rate,
		MaxRate:       maxRate,
		MaxChangeRate: maxChangeRate,
	}
}

// NewCommission returns an initialized validator commission.
func NewCommission(rate, maxRate, maxChangeRate sdk.Dec) Commission {
	return Commission{
		CommissionRates: NewCommissionRates(rate, maxRate, maxChangeRate),
		UpdateTime:      time.Unix(0, 0).UTC(),
	}
}

// NewCommissionWithTime returns an initialized validator commission with a specified
// update time which should be the current block BFT time.
func NewCommissionWithTime(rate, maxRate, maxChangeRate sdk.Dec, updatedAt time.Time) Commission {
	return Commission{
		CommissionRates: NewCommissionRates(rate, maxRate, maxChangeRate),
		UpdateTime:      updatedAt,
	}
}

// String implements the Stringer interface for a Commission object.
func (c Commission) String() string {
	out, _ := yaml.Marshal(c)
	return string(out)
}

// String implements the Stringer interface for a CommissionRates object.
func (cr CommissionRates) String() string {
	out, _ := yaml.Marshal(cr)
	return string(out)
}

// Validate performs basic sanity validation checks of initial commission
// parameters. If validation fails, an SDK error is returned.
func (cr CommissionRates) Validate() error {
	switch {
	case cr.MaxRate.IsNegative():
		// max rate cannot be negative
		return stakingTypes.ErrCommissionNegative

	case cr.MaxRate.GT(sdk.OneDec()):
		// max rate cannot be greater than 1
		return stakingTypes.ErrCommissionHuge

	case cr.Rate.IsNegative():
		// rate cannot be negative
		return stakingTypes.ErrCommissionNegative

	case cr.Rate.GT(cr.MaxRate):
		// rate cannot be greater than the max rate
		return stakingTypes.ErrCommissionGTMaxRate

	case cr.MaxChangeRate.IsNegative():
		// change rate cannot be negative
		return stakingTypes.ErrCommissionChangeRateNegative

	case cr.MaxChangeRate.GT(cr.MaxRate):
		// change rate cannot be greater than the max rate
		return stakingTypes.ErrCommissionChangeRateGTMaxRate
	}

	return nil
}

// ValidateNewRate performs basic sanity validation checks of a new commission
// rate. If validation fails, an SDK error is returned.
func (c Commission) ValidateNewRate(newRate sdk.Dec, blockTime time.Time) error {
	switch {
	case blockTime.Sub(c.UpdateTime).Hours() < 24:
		// new rate cannot be changed more than once within 24 hours
		return stakingTypes.ErrCommissionUpdateTime

	case newRate.IsNegative():
		// new rate cannot be negative
		return stakingTypes.ErrCommissionNegative

	case newRate.GT(c.MaxRate):
		// new rate cannot be greater than the max rate
		return stakingTypes.ErrCommissionGTMaxRate

	case newRate.Sub(c.Rate).GT(c.MaxChangeRate):
		// new rate % points change cannot be greater than the max change rate
		return stakingTypes.ErrCommissionGTMaxChangeRate
	}

	return nil
}
