package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the staking Querier
const (
	
	QueryBurnLevels = "burn_levels"

	
	QueryBurnLevel = "burn_level"

	
	QueryGatewayClusters = "query_gateway_clusters"
	
	QueryPersonClusterInfo = "person_cluster_info"
	//id
	QueryClusterInfoById = "cluster_info_by_id"
	//gas
	QueryClusterGasReward = "cluster_gas_reward"
	
	QueryClusterOwnerReward = "cluster_owner_reward"
	
	QueryClusterDeviceReward = "cluster_device_reward"
	
	QueryDeductionFee = "cluster_deduction_fee"
	
	QueryInClusters = "in_clusters"
	
	QueryClusterInfo = "cluster_info"
	
	QueryCluster = "cluster"
	
	QueryClusterPersonInfo = "query_cluster_person_info"
	//dvm
	QueryDvmList = "query_dvm_list"
	//dao
	QueryDaoParams = "query_dao_params"
	
	QueryClusterPersonals = "query_cluster_personals"
	
	QueryClusterPersonalInfo = "query_cluster_personal_info"
	
	QueryClusterProposalVoter = "query_cluster_personal_voter"
	
	QueryClusterProposalVoters = "query_cluster_personal_voters"
	
	QueryGroupMembers = "query_group_members"
	
	QueryGroupInfo = "query_group_info"
	
	QueryClusterApproveInfo = "query_cluster_approve_info"
)

type QueryClusterParams struct {
	ClusterId string `json:"cluster_id"`
}

type QueryClusterRewardParams struct {
	Member    string `json:"member"`
	ClusterId string `json:"cluster_id"`
}

func NewRewardParams(member, clusterId string) QueryClusterRewardParams {
	return QueryClusterRewardParams{member, clusterId}
}

// defines the params for the following queries:
// - 'custom/staking/delegatorDelegations'
// - 'custom/staking/delegatorUnbondingDelegations'
// - 'custom/staking/delegatorValidators'
type QueryDelegatorParams struct {
	DelegatorAddr sdk.AccAddress
}

func NewQueryDelegatorParams(delegatorAddr sdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddr: delegatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/staking/validator'
// - 'custom/staking/validatorDelegations'
// - 'custom/staking/validatorUnbondingDelegations'
type QueryValidatorParams struct {
	ValidatorAddr sdk.ValAddress
	Page, Limit   int
}

func NewQueryValidatorParams(validatorAddr sdk.ValAddress, page, limit int) QueryValidatorParams {
	return QueryValidatorParams{
		ValidatorAddr: validatorAddr,
		Page:          page,
		Limit:         limit,
	}
}

// defines the params for the following queries:
// - 'custom/staking/redelegation'
type QueryRedelegationParams struct {
	DelegatorAddr    sdk.AccAddress
	SrcValidatorAddr sdk.ValAddress
	DstValidatorAddr sdk.ValAddress
}

func NewQueryRedelegationParams(delegatorAddr sdk.AccAddress, srcValidatorAddr, dstValidatorAddr sdk.ValAddress) QueryRedelegationParams {
	return QueryRedelegationParams{
		DelegatorAddr:    delegatorAddr,
		SrcValidatorAddr: srcValidatorAddr,
		DstValidatorAddr: dstValidatorAddr,
	}
}

// QueryValidatorsParams defines the params for the following queries:
// - 'custom/staking/validators'
type QueryValidatorsParams struct {
	Page, Limit int
	Status      string
}

func NewQueryValidatorsParams(page, limit int, status string) QueryValidatorsParams {
	return QueryValidatorsParams{page, limit, status}
}

// QueryValidatorByConsAddrParams 
type QueryValidatorByConsAddrParams struct {
	ValidatorConsAddress sdk.ConsAddress
}

// ValidatorInfor 
type ValidatorInfor struct {
	ValidatorConsAddr string `json:"validator_consaddr"` 
	ValidatorStatus   string `json:"validator_status"`   //0  Unbonded 1 Unbonding 2 Bonded 3  4 
	ValidatorPubAddr  string `json:"validator_pubaddr"`  
	ValidatorOperAddr string `json:"validator_operaddr"` 
	AccAddress        string `json:"acc_address"`        
	ValidatorPubKey   string `json:"validator_pubkey"`   
}

type QueryBurnLevelsParams struct {
	Addresses []string
}

// QueryPersonClusterInfoRequest 
type QueryPersonClusterInfoRequest struct {
	From string `json:"from"`
}

// QueryGatewayClustersParams 
type QueryGatewayClustersParams struct {
	GatewayAddress string `json:"gateway_address"`
}

// InClusters 
type InClusters struct {
	ClusterId string
	IsOwner   bool
}

/*

type DeviceCluster struct {
	ClusterId             string                         `json:"cluster_id"`               //id
	ClusterChatId         string                         `json:"cluster_chat_id"`          //id
	ClusterName           string                         `json:"cluster_name"`             
	ClusterOwner          string                         `json:"cluster_owner"`            
	ClusterGateway        string                         `json:"cluster_gateway"`          
	ClusterLeader         string                         `json:"cluster_leader"`           
	ClusterDeviceMembers  map[string]ClusterDeviceMember `json:"cluster_device_members"`   
	ClusterPowerMembers   map[string]ClusterPowerMember  `json:"cluster_poweer_members"`   
	ClusterPower          sdk.Dec                        `json:"cluster_power"`            
	ClusterLevel          int64                          `json:"cluster_level"`            
	ClusterBurnAmount     sdk.Dec                        `json:"cluster_burn_amount"`      
	ClusterActiveDevice   int64                          `json:"cluster_active_device"`    
	ClusterDaoPool        string                         `json:"cluster_dao_pool"`         //dao
	ClusterDeviceRatio    sdk.Dec                        `json:"cluster_device_ratio"`     
	ClusterSalaryRatio    sdk.Dec                        `json:"cluster_salary_ratio"`     
	OnlineRatio           sdk.Dec                        `json:"online_ratio"`             
	OnlineRatioUpdateTime int64                          `json:"online_ratio_update_time"` 
	ClusterAdminList      map[string]struct{}            `json:"cluster_admin_list"`       
	ClusterVoteId         int64                          `json:"cluster_vote_id"`          //id
}
*/

// ClusterInfo 
type ClusterInfo struct {
	//id
	ClusterId string `json:"cluster_id"`
	//id
	ClusterChatId string `json:"cluster_chat_id"`
	
	ClusterOwner string `json:"cluster_owner"`
	
	ClusterName string `json:"cluster_name"`
	
	ClusterAllBurn sdk.Dec `json:"cluster_all_burn"`
	
	ClusterAllPower sdk.Dec `json:"cluster_all_power"`
	
	OnlineRatio sdk.Dec `json:"online_ratio"` 
	
	ClusterActiveDevice int64 `json:"cluster_active_device"` 
	
	ClusterDeviceAmount int64 `json:"cluster_device_amount"` 
	
	DeviceConnectivityRate sdk.Dec `json:"device_connectivity_rate"` 
	
	ClusterDeviceRatio sdk.Dec `json:"cluster_device_ratio"` 
	
	ClusterSalaryRatio sdk.Dec `json:"cluster_salary_ratio"` 
	//gas
	ClusterDayFreeGas sdk.Dec `json:"cluster_day_free_gas"`
	//DAO
	ClusterDaoPoolPower sdk.Dec `json:"cluster_dao_pool_power"`
	//DAOGAS
	DaoPoolDayFreeGas sdk.Dec `json:"dao_pool_day_free_gas"`
	
	DaoPoolAvailableAmount sdk.Dec `json:"dao_pool_available_amount"`
	//DAO
	DaoLicensingContract string `json:"dao_licensing_contract"`
	
	DaoLicensingHeight int64 `json:"dao_licensing_height"`
	
	LevelInfo ClusterLevelInfo `json:"level_info"`
}

// ClusterLevelInfo 
type ClusterLevelInfo struct {
	
	Level int64 `json:"level"`
	
	BurnAmountNextLevel math.Int `json:"burn_amount_next_level"`
	
	ActiveAmountNextLevel int64 `json:"active_amount_next_level"`
}

// ClusterPersonalInfo 
type ClusterPersonalInfo struct {
	PowerAmount  sdk.Dec `json:"power_amount"`  
	GasDay       sdk.Dec `json:"gas_day"`       //gas
	BurnAmount   sdk.Dec `json:"burn_amount"`   
	IsDevice     bool    `json:"is_device"`     
	IsAdmin      bool    `json:"is_amind"`      
	IsOwner      bool    `json:"is_owner"`      
	BurnRatio    sdk.Dec `json:"burn_ratio"`    
	PowerReward  sdk.Dec `json:"power_reward"`  
	DeviceReward sdk.Dec `json:"device_reward"` 
	OwnerReward  sdk.Dec `json:"owner_reward"`  
	AuthContract string  `json:"auth_contract"` 
	AuthHeight   int64   `json:"auth_height"`   
	ClusterOwner string  `json:"cluster_owner"` 
	ClusterName  string  `json:"cluster_name"`  
}

type QueryClusterPersonalInfoParams struct {
	ClusterId   string `json:"cluster_id"`
	FromAddress string `json:"from_address"`
}

type DvmInfo struct {
	//id
	ClusterChatId string `json:"cluster_chat_id"`
	//id
	ClusterId string `json:"cluster_id"`
	
	PowerReward sdk.Dec `json:"power_reward"`
	//dvm
	PowerDvm sdk.Dec `json:"power_dvm"`
	//gas
	GasDayDvm sdk.Dec `json:"gas_day_dvm"`
	
	AuthContract string `json:"auth_contract"`
	
	AuthHeight int64 `json:"auth_height"`
	
	ClusterName string `json:"cluster_name"`
}

// PersonClusterStatisticsInfo 
type PersonClusterStatisticsInfo struct {
	Address string `json:"address"`
	
	Owner []string `json:"owner"`
	
	BePower []string `json:"be_power"`
	
	AllBurn sdk.Dec `json:"all_burn"`
	// （）
	ActivePower sdk.Dec `json:"active_power"`
	
	FreezePower sdk.Dec `json:"freeze_power"`
	
	DeviceInfo []DeviceInfo `json:"device"`
}

type DeviceInfo struct {
	ClusterChatId string `json:"cluster_chat_id"`
	ClusterName   string `json:"cluster_name"`
	ClusterLevel  int64  `json:"cluster_level"`
	ClusterOwner  string `json:"cluster_owner"`
}

// dao
type DaoParams struct {
	//(1dst)
	BurnGetPowerRatio sdk.Dec `json:"burn_get_power_ratio"`

	
	SalaryRange Range `json:"salary_range"`

	
	DeviceRange Range `json:"device_range"`

	
	CreateClusterMinBurn math.Int `json:"create_cluster_min_burn"`

	//dst
	BurnAddress string `json:"burn_address"`

	
	DayBurnReward sdk.Dec `json:"day_burn_reward"`
}

type Range struct {
	Max sdk.Dec `json:"max"`
	Min sdk.Dec `json:"min"`
}
type QueryClusterProposalVoterParams struct {
	ProposalId string `json:"proposal_id"`
	Voter      string `json:"voter"`
}
