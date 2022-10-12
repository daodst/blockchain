package types

import (
	"freemasonry.cc/blockchain/cmd/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyMinMortgageCoin        = []byte("MinMortgageCoin")
	KeyMaxPhoneNumber         = []byte("MaxPhoneNumber")
	KeyDestroyPhoneNumberCoin = []byte("DestroyPhoneNumberCoin")
	KeyPreAttCoin             = []byte("PreAttCoin")
	KeyPreAttAccount          = []byte("PreAttAccount")
	KeyAttDestroyPercent      = []byte("AttDestroyPercent")
	KeyAttGatewayPercent      = []byte("AttGatewayPercent")
	KeyAttDposPercent         = []byte("AttDposPercent")
	KeyChatFee                = []byte("ChatFee")
)

// NewParams creates a new Params object
func NewParams(
	minMortgageCoin sdk.Coin,
	maxPhoneNumber uint64,
	destroyPhoneNumberCoin sdk.Coin,
	preAttCoin sdk.Coin,
	preAttAccount string,
	attDestroyPercent sdk.Dec,
	attGatewayPercent sdk.Dec,
	attDposPercent sdk.Dec,
	chatFee sdk.Coin,
) Params {
	return Params{
		MinMortgageCoin:        minMortgageCoin,
		MaxPhoneNumber:         maxPhoneNumber,
		DestroyPhoneNumberCoin: destroyPhoneNumberCoin,
		PreAttCoin:             preAttCoin,
		PreAttAccount:          preAttAccount,
		AttDestroyPercent:      attDestroyPercent,
		AttGatewayPercent:      attGatewayPercent,
		AttDposPercent:         attDposPercent,
		ChatFee:                chatFee,
	}
}

func DefaultParams() Params {

	preCoinAmount, ok := sdk.NewIntFromString("1000000000000000000000000")
	if !ok {
		panic("invalid preCoinAmount")
	}

	//todo -
	return Params{
		//CommunityAddress:  "dex1vmx0e3r7v93axstfkqpxjpvgearcmxls2x287f", //myself flower unusual veteran squeeze rally still can layer wear taste major lesson feed junk crystal left cart evidence egg arena legal dentist name
		//EcologicalAddress: "dex1wrx8tdyv5j3l5lst9n080aar3v0f5zywh303cy", //sail catalog midnight force pole combine charge audit tiny swim puppy muffin sister ginger pretty custom knife pig street race exact goat high crucial
		MinMortgageCoin:        sdk.NewCoin(config.BaseDenom, sdk.NewInt(1000000000000000000)),
		MaxPhoneNumber:         10,
		DestroyPhoneNumberCoin: sdk.NewCoin(config.BaseDenom, sdk.NewInt(1000000000000000000)),
		PreAttCoin:             sdk.NewCoin(config.BaseDenom, preCoinAmount), //  +  180
		PreAttAccount:          "dex1vmx0e3r7v93axstfkqpxjpvgearcmxls2x287f", //todo -
		AttDestroyPercent:      sdk.MustNewDecFromStr("0.05"),
		AttGatewayPercent:      sdk.MustNewDecFromStr("0.05"),
		AttDposPercent:         sdk.MustNewDecFromStr("0.05"),
		ChatFee:                sdk.NewCoin(config.BaseDenom, sdk.NewInt(1000000000000000000)),
	}
}

func (p Params) Validate() error {

	if err := validateCoin(p.MinMortgageCoin); err != nil {
		return err
	}

	if err := validateMaxPhoneNumber(p.MaxPhoneNumber); err != nil {
		return err
	}

	if err := validateCoin(p.DestroyPhoneNumberCoin); err != nil {
		return err
	}

	if err := validateCoin(p.PreAttCoin); err != nil {
		return err
	}

	if err := validateAddrString(p.PreAttAccount); err != nil {
		return err
	}

	if err := validateRatio(p.AttDestroyPercent); err != nil {
		return err
	}

	if err := validateRatio(p.AttGatewayPercent); err != nil {
		return err
	}

	if err := validateRatio(p.AttDposPercent); err != nil {
		return err
	}

	if err := validateCoin(p.ChatFee); err != nil {
		return err
	}
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMinMortgageCoin, &p.MinMortgageCoin, validateCoin),
		paramtypes.NewParamSetPair(KeyMaxPhoneNumber, &p.MaxPhoneNumber, validateMaxPhoneNumber),
		paramtypes.NewParamSetPair(KeyDestroyPhoneNumberCoin, &p.DestroyPhoneNumberCoin, validateCoin),
		paramtypes.NewParamSetPair(KeyPreAttCoin, &p.PreAttCoin, validateCoin),
		paramtypes.NewParamSetPair(KeyPreAttAccount, &p.PreAttAccount, validateAddrString),
		paramtypes.NewParamSetPair(KeyAttDestroyPercent, &p.AttDestroyPercent, validateRatio),
		paramtypes.NewParamSetPair(KeyAttGatewayPercent, &p.AttGatewayPercent, validateRatio),
		paramtypes.NewParamSetPair(KeyAttDposPercent, &p.AttDposPercent, validateRatio),
		paramtypes.NewParamSetPair(KeyChatFee, &p.ChatFee, validateCoin),
	}
}

func validateAddrString(i interface{}) error {

	return nil
}

//func validateChatRewardLog(i interface{}) error {
//    v, ok := i.([]ChatReward)
//    if !ok {
//        return fmt.Errorf("invalid parameter type: %T", i)
//    }
//
//    if len(v) == 0 {
//        return fmt.Errorf("invalid parameter len")
//    }
//
//    if len(v) > 1 {
//        for i := 0; i < len(v)-1; i++ {
//            if v[i].Height > v[i+1].Height {
//                return fmt.Errorf("invalid parameter height")
//            }
//
//            curValue := v[i].Value
//            curValueDec, err := sdk.NewDecFromStr(curValue)
//            if err != nil {
//                return fmt.Errorf("invalid value format")
//            }
//            if curValueDec.GT(sdk.NewDec(1)) {
//                return fmt.Errorf("invalid value (GT1)")
//            }
//            if curValueDec.LT(sdk.MustNewDecFromStr("0.0001")) {
//                return fmt.Errorf("invalid value (LT0.0001)")
//            }
//        }
//    }
//
//    return nil
//}

func validateMaxPhoneNumber(i interface{}) error {

	return nil
}

func validateCoin(i interface{}) error {

	return nil
}

func validateRatio(i interface{}) error {

	return nil
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(KeyMinMortgageCoin, DefaultParams().MinMortgageCoin, validateCoin),
		paramtypes.NewParamSetPair(KeyMaxPhoneNumber, DefaultParams().MaxPhoneNumber, validateMaxPhoneNumber),
		paramtypes.NewParamSetPair(KeyDestroyPhoneNumberCoin, DefaultParams().DestroyPhoneNumberCoin, validateCoin),
		paramtypes.NewParamSetPair(KeyPreAttCoin, DefaultParams().PreAttCoin, validateCoin),
		paramtypes.NewParamSetPair(KeyPreAttAccount, DefaultParams().PreAttAccount, validateAddrString),
		paramtypes.NewParamSetPair(KeyAttDestroyPercent, DefaultParams().AttDestroyPercent, validateRatio),
		paramtypes.NewParamSetPair(KeyAttGatewayPercent, DefaultParams().AttGatewayPercent, validateRatio),
		paramtypes.NewParamSetPair(KeyAttDposPercent, DefaultParams().AttDposPercent, validateRatio),
		paramtypes.NewParamSetPair(KeyChatFee, DefaultParams().ChatFee, validateCoin),
	)
}
