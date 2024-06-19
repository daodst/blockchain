package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2.
	cdc.RegisterConcrete(&MsgColonyRate{}, TypeMsgColonyRate, nil)
	cdc.RegisterConcrete(&MsgBurnToPower{}, TypeMsgBurnToPower, nil)
	cdc.RegisterConcrete(&MsgClusterChangeDeviceRatio{}, TypeMsgClusterChangeDeviceRatio, nil)
	cdc.RegisterConcrete(&MsgClusterChangeSalaryRatio{}, TypeMsgClusterChangeSalaryRatio, nil)
	cdc.RegisterConcrete(&MsgCreateCluster{}, TypeMsgCreateCluser, nil)
	cdc.RegisterConcrete(&MsgClusterAddMembers{}, TypeMsgClusterAddMembers, nil)
	cdc.RegisterConcrete(&MsgClusterChangeId{}, TypeMsgClusterChangeid, nil)
	cdc.RegisterConcrete(&MsgWithdrawOwnerReward{}, TypeMsgWithdrawOwnerReward, nil)
	cdc.RegisterConcrete(&MsgWithdrawBurnReward{}, TypeMsgWithdrawBurnReward, nil)
	cdc.RegisterConcrete(&MsgWithdrawDeviceReward{}, TypeMsgWithdrawDeviceReward, nil)
	cdc.RegisterConcrete(&MsgDeleteMembers{}, TypeMsgDeleteMembers, nil)
	cdc.RegisterConcrete(&MsgThawFrozenPower{}, TypeMsgThawFrozenPower, nil)
	cdc.RegisterConcrete(&MsgClusterMemberExit{}, TypeMsgClusterMemberExit, nil)
	cdc.RegisterConcrete(&MsgClusterChangeName{}, TypeMsgClusterChangeName, nil)
	cdc.RegisterConcrete(&MsgUpdateAdmin{}, TypeMsgUpdateAdmin, nil)
	cdc.RegisterConcrete(&MsgClusterPowerApprove{}, TypeMsgClusterPowerApprove, nil)
}

// RegisterInterfaces register implementations
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgColonyRate{},
		&MsgBurnToPower{},
		&MsgClusterChangeDeviceRatio{},
		&MsgClusterChangeSalaryRatio{},
		&MsgCreateCluster{},
		&MsgClusterAddMembers{},
		&MsgClusterChangeId{},
		&MsgWithdrawOwnerReward{},
		&MsgWithdrawBurnReward{},
		&MsgWithdrawDeviceReward{},
		&MsgDeleteMembers{},
		&MsgThawFrozenPower{},
		&MsgClusterMemberExit{},
		&MsgClusterChangeName{},
		&MsgUpdateAdmin{},
		&MsgClusterPowerApprove{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/smart module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
