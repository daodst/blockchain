package types

import (
	"freemasonry.cc/blockchain/cmd/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"time"
)

const DefaultHistoricalEntries uint32 = 10000

var (
	KeyMintDenom           = []byte("MintDenom")
	KeyInflationRateChange = []byte("InflationRateChange")
	KeyInflationMax        = []byte("InflationMax")
	KeyInflationMin        = []byte("InflationMin")
	KeyGoalBonded          = []byte("GoalBonded")
	KeyBlocksPerYear       = []byte("BlocksPerYear")
	KeyUnbondingTime       = []byte("UnbondingTime")
	KeyBondDenom           = []byte("BondDenom")

	KeyMaxValidators     = []byte("MaxValidators")
	KeyMaxEntries        = []byte("MaxEntries")
	KeyHistoricalEntries = []byte("HistoricalEntries")
	KeyPowerReduction    = []byte("PowerReduction")
	KeyMinMortgageCoin   = []byte("MinMortgageCoin")

	KeyAttDestroyPercent = []byte("AttDestroyPercent")
	KeyAttGatewayPercent = []byte("AttGatewayPercent")
	KeyAttDposPercent    = []byte("AttDposPercent")
	KeyPreAttCoin        = []byte("PreAttCoin")
	KeyPreAttAccount     = []byte("PreAttAccount")
)

var (
	ParamStoreKeyCommunityTax        = []byte("communitytax")
	ParamStoreKeyBaseProposerReward  = []byte("baseproposerreward")
	ParamStoreKeyBonusProposerReward = []byte("bonusproposerreward")
	ParamStoreKeyWithdrawAddrEnabled = []byte("withdrawaddrenabled")
	ParamStoreKeyMaxValidators       = []byte("maxvalidators")
	ParamStoreKeyMinMortgageCoin     = []byte("minmortgagecoin")

	ParamStoreKeyAttDestroyPercent = []byte("attdestroypercent")
	ParamStoreKeyAttGatewayPercent = []byte("attgatewaypercent")
	ParamStoreKeyAttDposPercent    = []byte("attdpospercent")
	ParamStoreKeyPreAttCoin        = []byte("preattcoin")
	ParamStoreKeyPreAttAccount     = []byte("preattaccount")
)

// NewParams creates a new Params object
func NewParams(
	mintDenom string,
	inflationRateChange sdk.Dec,
	inflationMax sdk.Dec,
	inflationMin sdk.Dec,
	goalBonded sdk.Dec,
	blocksPerYear uint64,
	unbondingTime time.Duration,
	bondDenom string,
	historicalEntries uint32,
	minMortgageCoin sdk.Coin,

	attDestroyPercent sdk.Dec,
	attGatewayPercent sdk.Dec,
	attDposPercent sdk.Dec,
	preAttCoin sdk.Coin,
	preAttAccount string,
) Params {
	return Params{
		MintDenom:           mintDenom,
		InflationRateChange: inflationRateChange,
		InflationMin:        inflationMin,
		InflationMax:        inflationMax,
		GoalBonded:          goalBonded,
		BlocksPerYear:       blocksPerYear,
		UnbondingTime:       unbondingTime,
		BondDenom:           bondDenom,
		HistoricalEntries:   historicalEntries,
		MinMortgageCoin:     minMortgageCoin,

		AttDestroyPercent: attDestroyPercent,
		AttGatewayPercent: attGatewayPercent,
		AttDposPercent:    attDposPercent,
		PreAttCoin:        preAttCoin,
		PreAttAccount:     preAttAccount,
	}
}

func DefaultParams() Params {
	minmortgageCoinInt, ok := sdk.NewIntFromString("100000000000000000000")
	if !ok {
		panic("Err Default Pledge Params MinMortgageCoin")
	}

	preCoinAmount, ok := sdk.NewIntFromString("1000000000000000000000000")
	if !ok {
		panic("invalid preCoinAmount")
	}

	return Params{
		MintDenom:           config.BaseDenom,
		InflationRateChange: sdk.NewDecWithPrec(13, 2),
		InflationMax:        sdk.MustNewDecFromStr("800"),
		InflationMin:        sdk.MustNewDecFromStr("100"),
		GoalBonded:          sdk.NewDecWithPrec(67, 2),
		BlocksPerYear:       uint64(60 * 60 * 8766 / 5),
		UnbondingTime:       time.Hour * 24 * 7 * 3,
		BondDenom:           config.BaseDenom,
		HistoricalEntries:   DefaultHistoricalEntries,
		MaxValidators:       100,
		MinMortgageCoin:     sdk.NewCoin(config.BaseDenom, minmortgageCoinInt),
		AttDestroyPercent:   sdk.MustNewDecFromStr("0.05"),
		AttGatewayPercent:   sdk.MustNewDecFromStr("0.05"),
		AttDposPercent:      sdk.MustNewDecFromStr("0.05"),

		PreAttCoin:    sdk.NewCoin(config.BaseDenom, preCoinAmount), //  +  180
		PreAttAccount: "dex1vmx0e3r7v93axstfkqpxjpvgearcmxls2x287f", //todo -
	}
}

func DefaultMinter() Minter {
	return Minter{
		Inflation:        sdk.Dec{},
		AnnualProvisions: sdk.Dec{},
	}
}

func (p Params) Validate() error {

	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateInflationRateChange(p.InflationRateChange); err != nil {
		return err
	}
	if err := validateInflationMax(p.InflationMax); err != nil {
		return err
	}
	if err := validateInflationMin(p.InflationMin); err != nil {
		return err
	}
	if err := validateGoalBonded(p.GoalBonded); err != nil {
		return err
	}
	if err := validateBlocksPerYear(p.BlocksPerYear); err != nil {
		return err
	}
	if err := validateUnbondingTime(p.UnbondingTime); err != nil {
		return err
	}
	if err := validateBondDenom(p.BondDenom); err != nil {
		return err
	}
	if err := validateHistoricalEntries(p.BondDenom); err != nil {
		return err
	}
	if err := validateMaxValidators(p.MaxValidators); err != nil {
		return err
	}
	if err := validateMinMortgageCoin(p.MinMortgageCoin); err != nil {
		return err
	}
	if err := validateAttDestroyPercent(p.AttDestroyPercent); err != nil {
		return err
	}
	if err := validateAttGatewayPercent(p.AttGatewayPercent); err != nil {
		return err
	}
	if err := validateAttDposPercent(p.AttDposPercent); err != nil {
		return err
	}
	if err := validatePreAttCoin(p.PreAttCoin); err != nil {
		return err
	}
	if err := validatePreAttAccount(p.PreAttAccount); err != nil {
		return err
	}
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyInflationRateChange, &p.InflationRateChange, validateInflationRateChange),
		paramtypes.NewParamSetPair(KeyInflationMax, &p.InflationMax, validateInflationMax),
		paramtypes.NewParamSetPair(KeyInflationMin, &p.InflationMin, validateInflationMin),
		paramtypes.NewParamSetPair(KeyGoalBonded, &p.GoalBonded, validateGoalBonded),
		paramtypes.NewParamSetPair(KeyBlocksPerYear, &p.BlocksPerYear, validateBlocksPerYear),
		paramtypes.NewParamSetPair(KeyUnbondingTime, &p.UnbondingTime, validateUnbondingTime),
		paramtypes.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
		paramtypes.NewParamSetPair(KeyHistoricalEntries, &p.HistoricalEntries, validateHistoricalEntries),
		paramtypes.NewParamSetPair(KeyMaxValidators, &p.MaxValidators, validateMaxValidators),
		paramtypes.NewParamSetPair(KeyMinMortgageCoin, &p.MinMortgageCoin, validateMinMortgageCoin),
		paramtypes.NewParamSetPair(KeyAttDestroyPercent, &p.AttDestroyPercent, validateAttDestroyPercent),
		paramtypes.NewParamSetPair(KeyAttGatewayPercent, &p.AttGatewayPercent, validateAttGatewayPercent),
		paramtypes.NewParamSetPair(KeyAttDposPercent, &p.AttDposPercent, validateAttDposPercent),
		paramtypes.NewParamSetPair(KeyPreAttCoin, &p.PreAttCoin, validatePreAttCoin),
		paramtypes.NewParamSetPair(KeyPreAttAccount, &p.PreAttAccount, validatePreAttAccount),
	}
}

func validateMintDenom(i interface{}) error {
	return nil
}
func validateInflationRateChange(i interface{}) error {
	return nil
}
func validateInflationMax(i interface{}) error {
	return nil
}
func validateInflationMin(i interface{}) error {
	return nil
}
func validateGoalBonded(i interface{}) error {
	return nil
}
func validateBlocksPerYear(i interface{}) error {
	return nil
}
func validateUnbondingTime(i interface{}) error {
	return nil
}
func validateBondDenom(i interface{}) error {
	return nil
}
func validateHistoricalEntries(i interface{}) error {
	return nil
}
func validateMaxValidators(i interface{}) error {
	return nil
}
func validateMinMortgageCoin(i interface{}) error {
	return nil
}
func validateAttDestroyPercent(i interface{}) error {
	return nil
}
func validateAttGatewayPercent(i interface{}) error {
	return nil
}
func validateAttDposPercent(i interface{}) error {
	return nil
}
func validatePreAttCoin(i interface{}) error {
	return nil
}
func validatePreAttAccount(i interface{}) error {
	return nil
}
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(KeyMintDenom, DefaultParams().MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyInflationRateChange, DefaultParams().InflationRateChange, validateInflationRateChange),
		paramtypes.NewParamSetPair(KeyInflationMax, DefaultParams().InflationMax, validateInflationMax),
		paramtypes.NewParamSetPair(KeyInflationMin, DefaultParams().InflationMin, validateInflationMin),
		paramtypes.NewParamSetPair(KeyGoalBonded, DefaultParams().GoalBonded, validateGoalBonded),
		paramtypes.NewParamSetPair(KeyBlocksPerYear, DefaultParams().BlocksPerYear, validateBlocksPerYear),
		paramtypes.NewParamSetPair(KeyUnbondingTime, DefaultParams().UnbondingTime, validateUnbondingTime),
		paramtypes.NewParamSetPair(KeyBondDenom, DefaultParams().BondDenom, validateBondDenom),
		paramtypes.NewParamSetPair(KeyHistoricalEntries, DefaultParams().HistoricalEntries, validateHistoricalEntries),
		paramtypes.NewParamSetPair(KeyMaxValidators, DefaultParams().MaxValidators, validateMaxValidators),
		paramtypes.NewParamSetPair(KeyMinMortgageCoin, DefaultParams().MinMortgageCoin, validateMinMortgageCoin),
		paramtypes.NewParamSetPair(KeyAttDestroyPercent, DefaultParams().AttDestroyPercent, validateAttDestroyPercent),
		paramtypes.NewParamSetPair(KeyAttGatewayPercent, DefaultParams().AttGatewayPercent, validateAttGatewayPercent),
		paramtypes.NewParamSetPair(KeyAttDposPercent, DefaultParams().AttDposPercent, validateAttDposPercent),
		paramtypes.NewParamSetPair(KeyPreAttCoin, DefaultParams().PreAttCoin, validatePreAttCoin),
		paramtypes.NewParamSetPair(KeyPreAttAccount, DefaultParams().PreAttAccount, validatePreAttAccount),
	)
}
