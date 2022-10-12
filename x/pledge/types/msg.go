package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgPledge{}
	_ sdk.Msg = &MsgUnpledge{}
	_ sdk.Msg = &MsgWithdrawDelegatorReward{}
)

const (
	TypeMsgPledge                  = "pledge"
	TypeMsgUnpledge                = "unpledge"
	TypeMsgWithdrawDelegatorReward = "withdraw_delegator_reward"
)

// NewMsgPledge ****************************/
func NewMsgPledge(fromAddress, delegatorAddress, validatorAddress string, amount sdk.Coin) *MsgPledge {
	return &MsgPledge{
		FromAddress:      fromAddress,
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
	}
}

func (msg MsgPledge) Route() string { return RouterKey }
func (msg MsgPledge) Type() string  { return TypeMsgPledge }
func (msg MsgPledge) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgPledge) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgPledge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}

	return nil
}
func (m MsgPledge) XXX_MessageName() string {
	return TypeMsgPledge
}

// NewMsgUnpledge ****************************/
func NewMsgUnpledge(delegatorAddress, validatorAddress string, amount sdk.Coin) *MsgUnpledge {
	return &MsgUnpledge{
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
	}
}

func (msg MsgUnpledge) Route() string { return RouterKey }
func (msg MsgUnpledge) Type() string  { return TypeMsgUnpledge }
func (msg MsgUnpledge) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgUnpledge) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgUnpledge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}

	return nil
}
func (m MsgUnpledge) XXX_MessageName() string {
	return TypeMsgUnpledge
}

// NewMsgPledgeRecieve ****************************/
func NewMsgMsgWithdrawDelegatorReward(delegatorAddress, validatorAddress string) *MsgWithdrawDelegatorReward {
	return &MsgWithdrawDelegatorReward{
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: validatorAddress,
	}
}

func (msg MsgWithdrawDelegatorReward) Route() string { return RouterKey }
func (msg MsgWithdrawDelegatorReward) Type() string  { return TypeMsgWithdrawDelegatorReward }
func (msg MsgWithdrawDelegatorReward) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgWithdrawDelegatorReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgWithdrawDelegatorReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}

	return nil
}
func (m MsgWithdrawDelegatorReward) XXX_MessageName() string {
	return TypeMsgWithdrawDelegatorReward
}
