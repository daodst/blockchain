package types

import (
	"fmt"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"math"
	"strconv"
)

var (
	KeyRate                         = []byte("Rate") //5%
	KeyMinDeviceRewardRatio         = []byte("MinDeviceRewardRatio")
	KeyMaxDeviceRewardRatio         = []byte("MaxDeviceRewardRatio")
	KeyMinSalaryRewardRatio         = []byte("MinSalaryRewardRatio")
	KeyMaxSalaryRewardRatio         = []byte("MaxSalaryRewardRatio")
	KeyBurnGetPowerRatio            = []byte("BurnGetPowerRatio")
	KeyClusterLevels                = []byte("ClusterLevels")
	KeyMaxClusterMembers            = []byte("MaxClusterMembers")
	KeyMinCreateClusterPledgeAmount = []byte("MinCreateClusterPledgeAmount")
	KeyDaoRewardPercent             = []byte("DaoRewardPercent")
	KeyDposRewardPercent            = []byte("DposRewardPercent")
	KeyBurnCurrentGateRatio         = []byte("BurnCurrentGateRatio")
	KeyBurnRegisterGateRatio        = []byte("BurnRegisterGateRatio")
	KeyDayMintAmount                = []byte("DayMintAmount")
	KeyBurnLevels                   = []byte("BurnLevels")
	KeyPowerGasRatio                = []byte("PowerGasRatio")
)

// NewParams creates a new Params object
func NewParams(
	rate,
	minDeviceRewardRatio,
	maxDeviceRewardRatio,
	minSalaryRewardRatio,
	maxSalaryRewardRatio,
	daoRewardPercent,
	dposRewardPercent,
	burnCurrentGateRatio,
	burnRegisterGateRati,
	burnGetPowerRatio sdk.Dec,
	maxClusterMembers int64,
	minCreateClusterPledgeAmount sdk.Int,
	clusterLevels []ClusterLevel,
	dayMintAmount sdk.Dec,
	burnLevels []BurnLevel,
	powerGasRatio sdk.Dec,
) Params {
	return Params{
		Rate:                         rate,
		MinDeviceRewardRatio:         minDeviceRewardRatio,
		MaxDeviceRewardRatio:         maxDeviceRewardRatio,
		MinSalaryRewardRatio:         minSalaryRewardRatio,
		MaxSalaryRewardRatio:         maxSalaryRewardRatio,
		BurnGetPowerRatio:            burnGetPowerRatio,
		ClusterLevels:                clusterLevels,
		MaxClusterMembers:            maxClusterMembers,
		MinCreateClusterPledgeAmount: minCreateClusterPledgeAmount,
		DaoRewardPercent:             daoRewardPercent,
		DposRewardPercent:            dposRewardPercent,
		BurnCurrentGateRatio:         burnCurrentGateRatio,
		BurnRegisterGateRatio:        burnRegisterGateRati,
		DayMintAmount:                dayMintAmount,
		BurnLevels:                   burnLevels,
		PowerGasRatio:                powerGasRatio,
	}
}

func DefaultParams() Params {

	return Params{
		Rate:                         sdk.NewDecWithPrec(5, 2),   //0.05
		MinDeviceRewardRatio:         sdk.NewDecWithPrec(1, 3),   
		MaxDeviceRewardRatio:         sdk.NewDecWithPrec(999, 3), //  999/1000
		MinSalaryRewardRatio:         sdk.ZeroDec(),              
		MaxSalaryRewardRatio:         sdk.NewDecWithPrec(999, 3), //  999/1000
		BurnGetPowerRatio:            sdk.NewDec(100),            
		ClusterLevels:                DefaultClusterLevelsInfo(),
		MaxClusterMembers:            666,                        
		MinCreateClusterPledgeAmount: sdk.NewInt(5000),           
		DaoRewardPercent:             sdk.NewDecWithPrec(1, 1),   //  10%    10%
		DposRewardPercent:            sdk.NewDecWithPrec(1, 1),   //  10%    dst
		BurnCurrentGateRatio:         sdk.NewDecWithPrec(575, 4), //  5.75%  dst
		BurnRegisterGateRatio:        sdk.NewDecWithPrec(375, 4), //  3.75%  dst
		DayMintAmount:                sdk.MustNewDecFromStr("360000000000000000000000"),
		BurnLevels:                   DefaultBurnLevels(),
		PowerGasRatio:                sdk.NewDec(100),
	}
}

func DefaultBurnLevels() []BurnLevel {
	burnLevels := make([]BurnLevel, 0)
	defaultAddPercemt := DefaultAddPercent()
	defaultRoomAmount := DefaultRoomAmount()
	for i := int64(1); i < int64(34); i++ {

		ifloat := float64(i)
		amountBaseFloat64 := math.Pow(2, ifloat-1) * 100
		amountFloat64 := amountBaseFloat64 * float64(core.RealToLedgerRateInt64)
		amountString := strconv.FormatFloat(amountFloat64, 'f', 0, 64)
		amountInt, err := sdk.NewDecFromStr(amountString)

		if err != nil {
			panic("Err Default Dao Params BurnLevel")
		}

		burnLevel := BurnLevel{
			Level:      i,
			BurnAmount: amountInt,
			AddPercent: defaultAddPercemt[i],
			RoomAmount: defaultRoomAmount[i],
		}

		burnLevels = append(burnLevels, burnLevel)
	}

	return burnLevels
}

func DefaultRoomAmount() map[int64]sdk.Int {
	//res := make(map[int64]sdk.Int, 33)
	return map[int64]sdk.Int{
		1:  sdk.NewInt(0),
		2:  sdk.NewInt(0),
		3:  sdk.NewInt(0),
		4:  sdk.NewInt(0),
		5:  sdk.NewInt(1),
		6:  sdk.NewInt(2),
		7:  sdk.NewInt(3),
		8:  sdk.NewInt(4),
		9:  sdk.NewInt(5),
		10: sdk.NewInt(6),
		11: sdk.NewInt(7),
		12: sdk.NewInt(8),
		13: sdk.NewInt(9),
		14: sdk.NewInt(10),
		15: sdk.NewInt(11),
		16: sdk.NewInt(12),
		17: sdk.NewInt(13),
		18: sdk.NewInt(14),
		19: sdk.NewInt(15),
		20: sdk.NewInt(16),
		21: sdk.NewInt(17),
		22: sdk.NewInt(18),
		23: sdk.NewInt(19),
		24: sdk.NewInt(20),
		25: sdk.NewInt(21),
		26: sdk.NewInt(22),
		27: sdk.NewInt(23),
		28: sdk.NewInt(24),
		29: sdk.NewInt(25),
		30: sdk.NewInt(26),
		31: sdk.NewInt(27),
		32: sdk.NewInt(28),
		33: sdk.NewInt(29),
	}
}

func DefaultAddPercent() map[int64]sdk.Int {
	//res := make(map[int64]sdk.Int, 33)
	return map[int64]sdk.Int{
		1:  sdk.NewInt(1),
		2:  sdk.NewInt(2),
		3:  sdk.NewInt(3),
		4:  sdk.NewInt(4),
		5:  sdk.NewInt(5),
		6:  sdk.NewInt(6),
		7:  sdk.NewInt(7),
		8:  sdk.NewInt(8),
		9:  sdk.NewInt(9),
		10: sdk.NewInt(12),
		11: sdk.NewInt(15),
		12: sdk.NewInt(18),
		13: sdk.NewInt(21),
		14: sdk.NewInt(24),
		15: sdk.NewInt(27),
		16: sdk.NewInt(30),
		17: sdk.NewInt(33),
		18: sdk.NewInt(36),
		19: sdk.NewInt(39),
		20: sdk.NewInt(43),
		21: sdk.NewInt(47),
		22: sdk.NewInt(51),
		23: sdk.NewInt(55),
		24: sdk.NewInt(59),
		25: sdk.NewInt(63),
		26: sdk.NewInt(67),
		27: sdk.NewInt(71),
		28: sdk.NewInt(75),
		29: sdk.NewInt(79),
		30: sdk.NewInt(84),
		31: sdk.NewInt(89),
		32: sdk.NewInt(94),
		33: sdk.NewInt(100),
	}
}

func (p Params) Validate() error {
	if err := validateRate(p.Rate); err != nil {
		return err
	}
	if err := validateMinDeviceRewardRatio(p.MinDeviceRewardRatio); err != nil {
		return err
	}
	if err := validateMaxDeviceRewardRatio(p.MaxDeviceRewardRatio); err != nil {
		return err
	}
	if err := validateMinSalaryRewardRatio(p.MinSalaryRewardRatio); err != nil {
		return err
	}
	if err := validateMaxSalaryRewardRatio(p.MaxSalaryRewardRatio); err != nil {
		return err
	}
	if err := validateBurnGetPowerRatio(p.BurnGetPowerRatio); err != nil {
		return err
	}
	if err := validateClusterLevels(p.ClusterLevels); err != nil {
		return err
	}
	if err := validateMaxClusterMembers(p.MaxClusterMembers); err != nil {
		return err
	}
	if err := validateMinCreateClusterPledgeAmount(p.MinCreateClusterPledgeAmount); err != nil {
		return err
	}
	if err := validateDaoRewardPercent(p.DaoRewardPercent); err != nil {
		return err
	}
	if err := validateDposRewardPercent(p.DposRewardPercent); err != nil {
		return err
	}
	if err := validateBurnCurrentGateRatio(p.BurnCurrentGateRatio); err != nil {
		return err
	}
	if err := validateBurnRegisterGateRatio(p.BurnRegisterGateRatio); err != nil {
		return err
	}
	if err := validateDayMintAmount(p.DayMintAmount); err != nil {
		return err
	}
	if err := validateBurnLevels(p.BurnLevels); err != nil {
		return err
	}
	if err := validatePowerGasRatio(p.PowerGasRatio); err != nil {
		return err
	}
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyRate, &p.Rate, validateRate),
		paramtypes.NewParamSetPair(KeyMinDeviceRewardRatio, &p.MinDeviceRewardRatio, validateMinDeviceRewardRatio),
		paramtypes.NewParamSetPair(KeyMaxDeviceRewardRatio, &p.MaxDeviceRewardRatio, validateMaxDeviceRewardRatio),
		paramtypes.NewParamSetPair(KeyMinSalaryRewardRatio, &p.MinSalaryRewardRatio, validateMinSalaryRewardRatio),
		paramtypes.NewParamSetPair(KeyMaxSalaryRewardRatio, &p.MaxSalaryRewardRatio, validateMaxSalaryRewardRatio),
		paramtypes.NewParamSetPair(KeyBurnGetPowerRatio, &p.BurnGetPowerRatio, validateBurnGetPowerRatio),
		paramtypes.NewParamSetPair(KeyClusterLevels, &p.ClusterLevels, validateClusterLevels),
		paramtypes.NewParamSetPair(KeyMaxClusterMembers, &p.MaxClusterMembers, validateMaxClusterMembers),
		paramtypes.NewParamSetPair(KeyMinCreateClusterPledgeAmount, &p.MinCreateClusterPledgeAmount, validateMinCreateClusterPledgeAmount),
		paramtypes.NewParamSetPair(KeyDaoRewardPercent, &p.DaoRewardPercent, validateDaoRewardPercent),
		paramtypes.NewParamSetPair(KeyDposRewardPercent, &p.DposRewardPercent, validateDposRewardPercent),
		paramtypes.NewParamSetPair(KeyBurnCurrentGateRatio, &p.BurnCurrentGateRatio, validateBurnCurrentGateRatio),
		paramtypes.NewParamSetPair(KeyBurnRegisterGateRatio, &p.BurnRegisterGateRatio, validateBurnRegisterGateRatio),
		paramtypes.NewParamSetPair(KeyDayMintAmount, &p.DayMintAmount, validateDayMintAmount),
		paramtypes.NewParamSetPair(KeyBurnLevels, &p.BurnLevels, validateBurnLevels),
		paramtypes.NewParamSetPair(KeyPowerGasRatio, &p.PowerGasRatio, validatePowerGasRatio),
	}
}

func validateRate(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateBurnLevels(i interface{}) error {
	burnLevels, ok := i.([]BurnLevel)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	burnCount := len(burnLevels)

	if burnCount == 0 {
		return core.ErrParamsBurnLevels
	}

	
	if burnLevels[0].Level != 1 {
		return core.ErrParamsBurnLevels
	}

	
	for n := 1; n < len(burnLevels)+1; n++ {

		if burnCount > 1 && n != burnCount {

			
			if burnLevels[n-1].Level+1 != burnLevels[n].Level {
				return core.ErrParamsLevel
			}

			
			if burnLevels[n-1].BurnAmount.GTE(burnLevels[n].BurnAmount) {
				return core.ErrParamsPeldgeLevel
			}
		}
	}
	return nil
}

func validatePowerGasRatio(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMinDeviceRewardRatio(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
func validateMaxDeviceRewardRatio(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
func validateMinSalaryRewardRatio(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
func validateMaxSalaryRewardRatio(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
func validateBurnGetPowerRatio(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
func validateClusterLevels(i interface{}) error {

	levels, ok := i.([]ClusterLevel)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(levels) == 0 {
		return core.ErrClusterLevelsParams
	}

	levelsCount := len(levels)

	
	if levels[0].Level != 1 {
		return core.ErrClusterLevelsParams
	}

	
	for n := 1; n < len(levels)+1; n++ {

		if levelsCount > 1 && n != levelsCount {

			
			if levels[n-1].Level+1 != levels[n].Level {
				return core.ErrClusterLevelsParams
			}

			
			if levels[n-1].BurnAmount.GT(levels[n].BurnAmount) {
				return core.ErrClusterLevelsParams
			}

			
			
			if levels[n-1].MemberAmount > levels[n].MemberAmount {
				return core.ErrClusterLevelsParams
			}
		}
	}

	return nil
}
func validateMaxClusterMembers(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
func validateMinCreateClusterPledgeAmount(i interface{}) error {
	_, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
func validateDaoRewardPercent(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
func validateDposRewardPercent(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
func validateBurnCurrentGateRatio(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
func validateBurnRegisterGateRatio(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
func validateDayMintAmount(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(KeyRate, DefaultParams().Rate, validateRate),
		paramtypes.NewParamSetPair(KeyMinDeviceRewardRatio, DefaultParams().MinDeviceRewardRatio, validateMinDeviceRewardRatio),
		paramtypes.NewParamSetPair(KeyMaxDeviceRewardRatio, DefaultParams().MaxDeviceRewardRatio, validateMaxDeviceRewardRatio),
		paramtypes.NewParamSetPair(KeyMinSalaryRewardRatio, DefaultParams().MinSalaryRewardRatio, validateMinSalaryRewardRatio),
		paramtypes.NewParamSetPair(KeyMaxSalaryRewardRatio, DefaultParams().MaxSalaryRewardRatio, validateMaxSalaryRewardRatio),
		paramtypes.NewParamSetPair(KeyBurnGetPowerRatio, DefaultParams().BurnGetPowerRatio, validateBurnGetPowerRatio),
		paramtypes.NewParamSetPair(KeyClusterLevels, DefaultParams().ClusterLevels, validateClusterLevels),
		paramtypes.NewParamSetPair(KeyMaxClusterMembers, DefaultParams().MaxClusterMembers, validateMaxClusterMembers),
		paramtypes.NewParamSetPair(KeyMinCreateClusterPledgeAmount, DefaultParams().MinCreateClusterPledgeAmount, validateMinCreateClusterPledgeAmount),
		paramtypes.NewParamSetPair(KeyDaoRewardPercent, DefaultParams().DaoRewardPercent, validateDaoRewardPercent),
		paramtypes.NewParamSetPair(KeyDposRewardPercent, DefaultParams().DposRewardPercent, validateDposRewardPercent),
		paramtypes.NewParamSetPair(KeyBurnCurrentGateRatio, DefaultParams().BurnCurrentGateRatio, validateBurnCurrentGateRatio),
		paramtypes.NewParamSetPair(KeyBurnRegisterGateRatio, DefaultParams().BurnRegisterGateRatio, validateBurnRegisterGateRatio),
		paramtypes.NewParamSetPair(KeyDayMintAmount, DefaultParams().DayMintAmount, validateDayMintAmount),
		paramtypes.NewParamSetPair(KeyBurnLevels, DefaultParams().BurnLevels, validateBurnLevels),
		paramtypes.NewParamSetPair(KeyPowerGasRatio, DefaultParams().PowerGasRatio, validatePowerGasRatio),
	)
}

func DefaultClusterLevelsInfo() []ClusterLevel {
	clusterLevels := make([]ClusterLevel, 0)
	defaultBurnAmount := DefaultClusterMinPledge()
	defaultMembersAmount := DefaultClusterMinMembers()
	for i := int64(1); i < int64(34); i++ {
		clusterLevel := ClusterLevel{
			Level:        i,
			BurnAmount:   defaultBurnAmount[i],
			MemberAmount: defaultMembersAmount[i],
		}

		clusterLevels = append(clusterLevels, clusterLevel)
	}

	return clusterLevels
}

func DefaultClusterMinPledge() map[int64]sdk.Int {
	return map[int64]sdk.Int{
		1:  sdk.NewInt(5000),
		2:  sdk.NewInt(10000),
		3:  sdk.NewInt(15000),
		4:  sdk.NewInt(20000),
		5:  sdk.NewInt(25000),
		6:  sdk.NewInt(30000),
		7:  sdk.NewInt(40000),
		8:  sdk.NewInt(50000),
		9:  sdk.NewInt(65000),
		10: sdk.NewInt(80000),
		11: sdk.NewInt(100000),
		12: sdk.NewInt(120000),
		13: sdk.NewInt(145000),
		14: sdk.NewInt(175000),
		15: sdk.NewInt(210000),
		16: sdk.NewInt(250000),
		17: sdk.NewInt(310000),
		18: sdk.NewInt(375000),
		19: sdk.NewInt(455000),
		20: sdk.NewInt(550000),
		21: sdk.NewInt(660000),
		22: sdk.NewInt(795000),
		23: sdk.NewInt(955000),
		24: sdk.NewInt(1150000),
		25: sdk.NewInt(1385000),
		26: sdk.NewInt(1665000),
		27: sdk.NewInt(2000000),
		28: sdk.NewInt(2405000),
		29: sdk.NewInt(2890000),
		30: sdk.NewInt(3470000),
		31: sdk.NewInt(4165000),
		32: sdk.NewInt(5000000),
		33: sdk.NewInt(6000000),
	}
}

func DefaultClusterMinMembers() map[int64]int64 {
	//res := make(map[int64]sdk.Int, 33)
	return map[int64]int64{
		1:  1,
		2:  2,
		3:  3,
		4:  4,
		5:  5,
		6:  6,
		7:  8,
		8:  10,
		9:  12,
		10: 15,
		11: 18,
		12: 21,
		13: 25,
		14: 30,
		15: 35,
		16: 41,
		17: 48,
		18: 57,
		19: 67,
		20: 79,
		21: 93,
		22: 109,
		23: 128,
		24: 150,
		25: 175,
		26: 204,
		27: 238,
		28: 278,
		29: 325,
		30: 379,
		31: 442,
		32: 515,
		33: 600,
	}
}
