package types

import (
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/types/address"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

// constants
const (
	// module name
	ModuleName = "dao"
	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	QuerierRoute = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName
)

// ModuleAddress is the native module address for EVM

var (
	ModuleAddress common.Address
	//Cluster
	BurnStartingInfoPrefix          = []byte{0x00} // key for Burn starting info
	ClusterHistoricalRewardsPrefix  = []byte{0x01} // key for historical Cluster rewards / stake
	ClusterCurrentRewardsPrefix     = []byte{0x02} // key for current Cluster rewards
	ClusterOutstandingRewardsPrefix = []byte{0x03} // key for Cluster outstanding rewards
	//device
	DeviceStartingInfoPrefix       = []byte{0x04} // key for Device starting info
	DeviceHistoricalRewardsPrefix  = []byte{0x05} // key for historical Device rewards / stake
	DeviceCurrentRewardsPrefix     = []byte{0x06} // key for current Device rewards
	DeviceOutstandingRewardsPrefix = []byte{0x07} // key for Device outstanding rewards
	//key
	DeviceClusterKey = []byte{0x08}
	//key
	PersonClusterInfoKey = []byte{0x09}
	
	TotalBurnAmount = []byte{0x10}
	
	TotalPowerAmount = []byte{0x11}
	//idid
	ClusterIdKey = []byte{0x12}
	
	ClusterForGateway = []byte{0x13}
	
	ClusterDeductionFee = []byte{0x14}
	
	ClusterEvmDeductionFee = []byte{0x15}

	
	ClusterApprovePowerInfo = []byte{0x16}

	
	ClusterApproveInfo = []byte{0x17}
)

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

func GetClusterEvmDeductionFeeKey(clusterId string) []byte {
	return append(ClusterDeductionFee, address.MustLengthPrefix([]byte(clusterId))...)
}
func GetClusterApprovePowerInfoKey(contractAddress string) []byte {
	return append(ClusterApprovePowerInfo, address.MustLengthPrefix([]byte(contractAddress))...)
}
func GetClusterApproveInfoInfoKey(clusterId string) []byte {
	return append(ClusterApproveInfo, address.MustLengthPrefix([]byte(clusterId))...)
}
func GetClusterDeductionFeeKey(clusterId string) []byte {
	return append(ClusterDeductionFee, address.MustLengthPrefix([]byte(clusterId))...)
}

func GetClusterForGatewayKey(gatewayAddress string) []byte {
	return append(ClusterForGateway, []byte(gatewayAddress)...)
}

func GetClusterIdKey(clusterChatId string) []byte {
	return append(ClusterIdKey, []byte(clusterChatId)...)
}

func GetDeviceClusterKey(clusterId string) []byte {
	return append(DeviceClusterKey, []byte(clusterId)...)
}

func GetPersonClusterInfoKey(address string) []byte {
	return append(PersonClusterInfoKey, []byte(address)...)
}

func GetClusterHistoricalRewardsKey(v string, k uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, k)
	return append(append(ClusterHistoricalRewardsPrefix, address.MustLengthPrefix([]byte(v))...), b...)
}

func GetClusterCurrentRewardsKey(v string) []byte {
	return append(ClusterCurrentRewardsPrefix, address.MustLengthPrefix([]byte(v))...)
}

func GetClusterOutstandingRewardsKey(clusterId string) []byte {
	return append(ClusterOutstandingRewardsPrefix, address.MustLengthPrefix([]byte(clusterId))...)
}

func GetBurnStartingInfoKey(clusterId, memberAddress string) []byte {
	return append(append(BurnStartingInfoPrefix, address.MustLengthPrefix([]byte(clusterId))...), address.MustLengthPrefix([]byte(memberAddress))...)
}

func GetDeviceHistoricalRewardsKey(v string, k uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, k)
	return append(append(DeviceHistoricalRewardsPrefix, address.MustLengthPrefix([]byte(v))...), b...)
}

func GetDeviceCurrentRewardsKey(v string) []byte {
	return append(DeviceCurrentRewardsPrefix, address.MustLengthPrefix([]byte(v))...)
}

func GetDeviceOutstandingRewardsKey(clusterId string) []byte {
	return append(DeviceOutstandingRewardsPrefix, address.MustLengthPrefix([]byte(clusterId))...)
}

func GetDeviceStartingInfoKey(clusterId, memberAddress string) []byte {
	return append(append(DeviceStartingInfoPrefix, address.MustLengthPrefix([]byte(clusterId))...), address.MustLengthPrefix([]byte(memberAddress))...)
}
