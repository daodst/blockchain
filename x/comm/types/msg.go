package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"regexp"
)

var (
	_ sdk.Msg = &MsgGatewayRegister{}
	_ sdk.Msg = &MsgGatewayIndexNum{}
	_ sdk.Msg = &MsgGatewayUndelegate{}
)

const (
	TypeMsgCreateSmartValidator   = "create_smart_validator"
	TypeMsgGatewayRegister        = "gateway_register"
	TypeMsgGatewayIndexNum        = "gateway_index_num"
	TypeMsgGatewayUndelegation    = "gateway_undelegation"
	TypeMsgGatewayBeginRedelegate = "gateway_begin_redelegate"
)

func NewMsgCreateSmartValidator(
	valAddr sdk.ValAddress, pubKey string, //nolint:interfacer
	selfDelegation sdk.Coin, description stakingtypes.Description, commission stakingtypes.CommissionRates, minSelfDelegation sdk.Int,
) (*MsgCreateSmartValidator, error) {
	return &MsgCreateSmartValidator{
		Description:       description,
		DelegatorAddress:  sdk.AccAddress(valAddr).String(),
		ValidatorAddress:  valAddr.String(),
		PubKey:            pubKey,
		Value:             selfDelegation,
		Commission:        commission,
		MinSelfDelegation: minSelfDelegation,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgCreateSmartValidator) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCreateSmartValidator) Type() string { return TypeMsgCreateSmartValidator }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgCreateSmartValidator) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	addr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(delAddr.Bytes(), addr.Bytes()) {
		addrs = append(addrs, sdk.AccAddress(addr))
	}

	return addrs
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateSmartValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreateSmartValidator) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures either empty or valid
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return err
	}
	if delAddr.Empty() {
		return stakingtypes.ErrEmptyDelegatorAddr
	}

	if msg.ValidatorAddress == "" {
		return stakingtypes.ErrEmptyValidatorAddr
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return err
	}
	if !sdk.AccAddress(valAddr).Equals(delAddr) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "validator address is invalid")
	}

	if msg.PubKey == "" {
		return stakingtypes.ErrEmptyValidatorPubKey
	}

	if !msg.Value.IsValid() || !msg.Value.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid delegation amount")
	}

	if msg.Description == (stakingtypes.Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.Commission == (stakingtypes.CommissionRates{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty commission")
	}

	if err := msg.Commission.Validate(); err != nil {
		return err
	}

	if !msg.MinSelfDelegation.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"minimum self delegation must be a positive integer",
		)
	}

	if msg.Value.Amount.LT(msg.MinSelfDelegation) {
		return stakingtypes.ErrSelfDelegationBelowMinimum
	}
	return nil
}

//
func NewMsgGatewayRegister(address, gatewayName, gatewayUrl, delegation string, indexNumber []string) *MsgGatewayRegister {
	return &MsgGatewayRegister{
		Address:     address,
		GatewayName: gatewayName,
		GatewayUrl:  gatewayUrl,
		Delegation:  delegation,
		IndexNumber: indexNumber,
	}
}

func (msg MsgGatewayRegister) Route() string { return RouterKey }
func (msg MsgGatewayRegister) Type() string  { return TypeMsgGatewayRegister }
func (msg MsgGatewayRegister) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgGatewayRegister) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgGatewayRegister) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid send address")
	}
	for _, val := range msg.IndexNumber {
		if len(val) != 7 {
			return ErrGatewayNumLength
		}
	}
	matched, err := regexp.MatchString("^(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5]):([0-9]|[1-9]\\d|[1-9]\\d{2}|[1-9]\\d{3}|[1-5]\\d{4}|6[0-4]\\d{3}|65[0-4]\\d{2}|655[0-2]\\d|6553[0-5])$", msg.GatewayUrl)
	if err != nil {
		return sdkerrors.Wrap(err, "gateway url")
	}
	if !matched {
		return sdkerrors.Wrap(err, "gateway url not match")
	}
	return nil
}

//
func NewMsgGatewayIndexNum(address, validatorAddress string, indexNumber []string) *MsgGatewayIndexNum {
	return &MsgGatewayIndexNum{
		DelegatorAddress: address,
		ValidatorAddress: validatorAddress,
		IndexNumber:      indexNumber,
	}
}
func (msg MsgGatewayIndexNum) Route() string { return RouterKey }
func (msg MsgGatewayIndexNum) Type() string  { return TypeMsgGatewayIndexNum }
func (msg MsgGatewayIndexNum) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgGatewayIndexNum) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgGatewayIndexNum) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid send address")
	}
	if len(msg.IndexNumber) == 0 {
		return sdkerrors.Wrap(err, "invalid message")
	}
	for _, val := range msg.IndexNumber {
		if len(val) != 7 {
			return ErrGatewayNumLength
		}
	}
	return nil
}

//
func NewMsgGatewayUndelegation(address, validatorAddress string, amount sdk.Coin, indexNumber []string) *MsgGatewayUndelegate {
	return &MsgGatewayUndelegate{
		DelegatorAddress: address,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
		IndexNumber:      indexNumber,
	}
}

func (msg MsgGatewayUndelegate) Route() string { return RouterKey }
func (msg MsgGatewayUndelegate) Type() string  { return TypeMsgGatewayUndelegation }
func (msg MsgGatewayUndelegate) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgGatewayUndelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgGatewayUndelegate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid send address")
	}
	if len(msg.IndexNumber) == 0 && msg.Amount.IsZero() {
		return sdkerrors.Wrap(err, "invalid message")
	}
	for _, val := range msg.IndexNumber {
		if len(val) != 7 {
			return ErrGatewayNumLength
		}
	}

	return nil
}

func NewMsgGatewayBeginRedelegate(
	delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, amount sdk.Coin, indexNumber []string) *MsgGatewayBeginRedelegate {
	return &MsgGatewayBeginRedelegate{
		DelegatorAddress:    delAddr.String(),
		ValidatorSrcAddress: valSrcAddr.String(),
		ValidatorDstAddress: valDstAddr.String(),
		Amount:              amount,
		IndexNumber:         indexNumber,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgGatewayBeginRedelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg MsgGatewayBeginRedelegate) Type() string { return TypeMsgGatewayBeginRedelegate }

// GetSigners implements the sdk.Msg interface
func (msg MsgGatewayBeginRedelegate) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delAddr}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgGatewayBeginRedelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgGatewayBeginRedelegate) ValidateBasic() error {
	if msg.DelegatorAddress == "" {
		return stakingtypes.ErrEmptyDelegatorAddr
	}

	if msg.ValidatorSrcAddress == "" {
		return stakingtypes.ErrEmptyValidatorAddr
	}

	if msg.ValidatorDstAddress == "" {
		return stakingtypes.ErrEmptyValidatorAddr
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid shares amount",
		)
	}
	for _, val := range msg.IndexNumber {
		if len(val) != 7 {
			return ErrGatewayNumLength
		}
	}
	return nil
}
