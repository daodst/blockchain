package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

//var ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgPledge{}, MsgTypePledge, nil)
	cdc.RegisterConcrete(&MsgWithdrawDelegatorReward{}, MsgTypeWithdrawDelegatorReward, nil)
	cdc.RegisterConcrete(&MsgUnpledge{}, MsgTypeUnpledge, nil)
	cdc.RegisterConcrete(&PledgeDelegateProposal{}, TypePledgeDelegateProposal, nil)
}

// RegisterInterfaces register implementations
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgPledge{},
		&MsgWithdrawDelegatorReward{},
		&MsgUnpledge{},
	)
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&PledgeDelegateProposal{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
