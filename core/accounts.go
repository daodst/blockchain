package core

import (
	"freemasonry.cc/blockchain/x/chat/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcTransferTypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

var (
	//
	ContractAddressFee = authtypes.NewModuleAddress(authtypes.FeeCollectorName)

	//
	ContractAddressBank = authtypes.NewModuleAddress(bankTypes.ModuleName)

	//
	ContractAddressDistribution = authtypes.NewModuleAddress(distrtypes.ModuleName)

	//staking 
	ContractAddressStakingBonded = authtypes.NewModuleAddress(stakingtypes.BondedPoolName)

	//staking 
	ContractAddressStakingNotBonded = authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName)

	ContractAddressGov = authtypes.NewModuleAddress(govtypes.ModuleName)

	//IBC
	ContractAddressIbcTransfer = authtypes.NewModuleAddress(ibcTransferTypes.ModuleName)

	//
	ContractGatewayBonus = authtypes.NewModuleAddress(GatewayBonusAddress)

	//
	ContractChatBurn = authtypes.NewModuleAddress(types.ModuleBurnName)

	//
	ContractChat = authtypes.NewModuleAddress(types.ModuleName)
)
