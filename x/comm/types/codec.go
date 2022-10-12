package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgCreateSmartValidator{}, MSG_SMART_CREATE_VALIDATOR, nil)
	cdc.RegisterConcrete(&MsgGatewayRegister{}, MSG_GATEWAY_REGISTER, nil)
	cdc.RegisterConcrete(&MsgGatewayIndexNum{}, MSG_GATEWAY_INDEX_NUM, nil)
	cdc.RegisterConcrete(&MsgGatewayUndelegate{}, MSG_GATEWAY_UNDELEGATION, nil)
	cdc.RegisterConcrete(&MsgGatewayBeginRedelegate{}, MSG_GATEWAY_BEGIN_REDELEGATE, nil)
}

// RegisterInterfaces register implementations
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgGatewayRegister{},
		&MsgGatewayIndexNum{},
		&MsgGatewayUndelegate{},
		&MsgGatewayBeginRedelegate{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/smart module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
