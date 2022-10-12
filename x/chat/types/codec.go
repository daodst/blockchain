package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

//var ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgRegister{}, MsgTypeRegister, nil)
	cdc.RegisterConcrete(&MsgSetChatFee{}, MsgTypeSetChatFee, nil)
	cdc.RegisterConcrete(&MsgSendGift{}, MsgTypeSendGift, nil)
	cdc.RegisterConcrete(&MsgAddressBookSave{}, MsgTypeAddressBookSave, nil)
	cdc.RegisterConcrete(&MsgMobileTransfer{}, MsgTypeMobileTransfer, nil)
	cdc.RegisterConcrete(&MsgChangeGateway{}, MsgTypeChangeGateway, nil)
	cdc.RegisterConcrete(&MsgBurnGetMobile{}, MsgTypeBurnGetMobile, nil)
	cdc.RegisterConcrete(&MsgSetChatInfo{}, MsgTypeSetChatInfo, nil)
}

// RegisterInterfaces register implementations
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegister{},
		&MsgSetChatFee{},
		&MsgSendGift{},
		&MsgAddressBookSave{},
		&MsgMobileTransfer{},
		&MsgChangeGateway{},
		&MsgBurnGetMobile{},
		&MsgSetChatInfo{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
