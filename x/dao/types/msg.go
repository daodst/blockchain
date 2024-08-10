package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgColonyRate{}
)

const (
	TypeMsgColonyRate               = "dao/ColonyRate"
	TypeMsgCreateCluser             = "dao/CreateCluster"
	TypeMsgClusterAddMembers        = "dao/ClusterAddMembers"
	TypeMsgBurnToPower              = "dao/BurnToPower"
	TypeMsgClusterChangeDeviceRatio = "dao/ClusterChangeDeviceRatio"
	TypeMsgClusterChangeSalaryRatio = "dao/ClusterChangeSalaryRatio"
	TypeMsgClusterChangeid          = "dao/ClusterChangeid"
	TypeMsgWithdrawOwnerReward      = "dao/WithdrawOwnerReward"
	TypeMsgWithdrawBurnReward       = "dao/WithdrawBurnReward"
	TypeMsgWithdrawDeviceReward     = "dao/WithdrawDeviceReward"
	TypeMsgDeleteMembers            = "dao/DeleteMembers"
	TypeMsgThawFrozenPower          = "dao/ThawFrozenPower"
	TypeMsgClusterMemberExit        = "dao/ClusterMemberExit"
	TypeMsgClusterChangeName        = "dao/ClusterChangeName"
	TypeMsgUpdateAdmin              = "dao/UpdateAdmin"
	TypeMsgClusterPowerApprove      = "dao/ClusterPowerApprove"
)

// NewMsgClusterChangeName 
func NewMsgUpdateAdmin(
	from, clusterId string, clusterAdminList []string,
) *MsgUpdateAdmin {
	return &MsgUpdateAdmin{
		FromAddress:      from,
		ClusterId:        clusterId,
		ClusterAdminList: clusterAdminList,
	}
}

// Route Implements Msg.
func (m MsgUpdateAdmin) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgUpdateAdmin) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgUpdateAdmin) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgUpdateAdmin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgUpdateAdmin) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

// NewMsgClusterChangeName 
func NewMsgClusterChangeName(
	from, clusterId, clusterName string,
) *MsgClusterChangeName {
	return &MsgClusterChangeName{
		FromAddress: from,
		ClusterId:   clusterId,
		ClusterName: clusterName,
	}
}

// Route Implements Msg.
func (m MsgClusterChangeName) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgClusterChangeName) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeName) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterChangeName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterChangeName) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

// NewMsgClusterMemberExit 
func NewMsgClusterMemberExit(fromAddress, clusterId string) *MsgClusterMemberExit {
	return &MsgClusterMemberExit{
		FromAddress: fromAddress,
		ClusterId:   clusterId,
	}
}

// Route Implements Msg.
func (m MsgClusterMemberExit) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgClusterMemberExit) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterMemberExit) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterMemberExit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterMemberExit) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

// NewMsgThawFrozenPower 
func NewMsgThawFrozenPower(fromAddress, clusterId, gatewayAddr, chatAddr string, thawAmount sdk.Dec) *MsgThawFrozenPower {
	return &MsgThawFrozenPower{
		FromAddress:    fromAddress,
		ClusterId:      clusterId,
		ThawAmount:     thawAmount,
		GatewayAddress: gatewayAddr,
		ChatAddress:    chatAddr,
	}
}

// Route Implements Msg.
func (m MsgThawFrozenPower) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgThawFrozenPower) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgThawFrozenPower) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgThawFrozenPower) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgThawFrozenPower) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

// NewMsgDeleteMembersd 
func NewMsgDeleteMembersd(fromAddress, clusterId string, members []string) *MsgDeleteMembers {
	return &MsgDeleteMembers{
		FromAddress: fromAddress,
		ClusterId:   clusterId,
		Members:     members,
	}
}

// Route Implements Msg.
func (m MsgDeleteMembers) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgDeleteMembers) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgDeleteMembers) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgDeleteMembers) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgDeleteMembers) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

// NewMsgWithdrawOwnerReward 
func NewMsgWithdrawOwnerReward(address, clusterId string) *MsgWithdrawOwnerReward {
	return &MsgWithdrawOwnerReward{
		Address:   address,
		ClusterId: clusterId,
	}
}

// Route Implements Msg.
func (m MsgWithdrawOwnerReward) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgWithdrawOwnerReward) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgWithdrawOwnerReward) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgWithdrawOwnerReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgWithdrawOwnerReward) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

// NewMsgWithdrawBurnReward 
func NewMsgWithdrawBurnReward(memberAddress, clusterId string) *MsgWithdrawBurnReward {
	return &MsgWithdrawBurnReward{
		MemberAddress: memberAddress,
		ClusterId:     clusterId,
	}
}

// Route Implements Msg.
func (m MsgWithdrawBurnReward) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgWithdrawBurnReward) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgWithdrawBurnReward) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.MemberAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgWithdrawBurnReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgWithdrawBurnReward) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.MemberAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}


func NewMsgWithdrawDeviceReward(memberAddress, clusterId string) *MsgWithdrawDeviceReward {
	return &MsgWithdrawDeviceReward{
		MemberAddress: memberAddress,
		ClusterId:     clusterId,
	}
}

// Route Implements Msg.
func (m MsgWithdrawDeviceReward) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgWithdrawDeviceReward) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgWithdrawDeviceReward) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.MemberAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgWithdrawDeviceReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgWithdrawDeviceReward) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.MemberAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

// NewMsgClusterChangeid id
func NewMsgClusterChangeId(
	from, clusterId, newClusterId string,
) *MsgClusterChangeId {
	return &MsgClusterChangeId{
		FromAddress:  from,
		ClusterId:    clusterId,
		NewClusterId: newClusterId,
	}
}

// Route Implements Msg.
func (m MsgClusterChangeId) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgClusterChangeId) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeId) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterChangeId) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterChangeId) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

// NewMsgClusterAddMembers 
func NewMsgClusterAddMembers(
	from, clusterId string, members []Members,
) *MsgClusterAddMembers {
	return &MsgClusterAddMembers{
		FromAddress: from,
		ClusterId:   clusterId,
		Members:     members,
	}
}

// Route Implements Msg.
func (m MsgClusterAddMembers) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgClusterAddMembers) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterAddMembers) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterAddMembers) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterAddMembers) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

// NewMsgCreateCluster 
func NewMsgCreateCluster(
	from, gateAddress, clusterId, chatAddress, clusterName string, deviceRatio, salaryRatio, burnAmount, freezeAmount sdk.Dec,
) *MsgCreateCluster {
	return &MsgCreateCluster{
		FromAddress:  from,
		GateAddress:  gateAddress,
		ClusterId:    clusterId,
		DeviceRatio:  deviceRatio,
		SalaryRatio:  salaryRatio,
		BurnAmount:   burnAmount,
		ChatAddress:  chatAddress,
		ClusterName:  clusterName,
		FreezeAmount: freezeAmount,
	}
}

// Route Implements Msg.
func (m MsgCreateCluster) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgCreateCluster) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgCreateCluster) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgCreateCluster) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgCreateCluster) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

// NewMsgClusterChangeSalaryRatio 
func NewMsgClusterChangeSalaryRatio(from, clusterId string, salaryRatio sdk.Dec) *MsgClusterChangeSalaryRatio {
	return &MsgClusterChangeSalaryRatio{
		FromAddress: from,
		ClusterId:   clusterId,
		SalaryRatio: salaryRatio,
	}
}

// Route Implements Msg.
func (m MsgClusterChangeSalaryRatio) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgClusterChangeSalaryRatio) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeSalaryRatio) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterChangeSalaryRatio) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterChangeSalaryRatio) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

// NewMsgClusterChangeDeviceRatio 
func NewMsgClusterChangeDeviceRatio(
	from, clusterId string, deviceRatio sdk.Dec,
) *MsgClusterChangeDeviceRatio {

	return &MsgClusterChangeDeviceRatio{
		FromAddress: from,
		ClusterId:   clusterId,
		DeviceRatio: deviceRatio,
	}
}

// Route Implements Msg.
func (m MsgClusterChangeDeviceRatio) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgClusterChangeDeviceRatio) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeDeviceRatio) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterChangeDeviceRatio) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterChangeDeviceRatio) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

// NewMsgBurnToPower 
func NewMsgBurnToPower(
	from, to, clusterId, gatewayAddress, chatAddress string, burnAmoumt, useFreezeAmount sdk.Dec,
) *MsgBurnToPower {

	return &MsgBurnToPower{
		FromAddress:     from,
		ToAddress:       to,
		ClusterId:       clusterId,
		BurnAmount:      burnAmoumt,
		UseFreezeAmount: useFreezeAmount,
		GatewayAddress:  gatewayAddress,
		ChatAddress:     chatAddress,
	}
}

// Route Implements Msg.
func (m MsgBurnToPower) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgBurnToPower) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgBurnToPower) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgBurnToPower) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgBurnToPower) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgColonyRate(
	address, gatewayAddress string,
	rate []ColonyRate,
) *MsgColonyRate {

	return &MsgColonyRate{
		Address:        address,
		GatewayAddress: gatewayAddress,
		OnlineRate:     rate,
	}
}

// Route Implements Msg.
func (m MsgColonyRate) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgColonyRate) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgColonyRate) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgColonyRate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgColonyRate) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}



// Route Implements Msg.
func (m MsgClusterPowerApprove) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgClusterPowerApprove) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterPowerApprove) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterPowerApprove) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterPowerApprove) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid from address")
	}
	return nil
}
