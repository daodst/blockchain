package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// DeviceCluster 
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
	ClusterVoteId         uint64                         `json:"cluster_vote_id"`          //id
	ClusterVotePolicy     string                         `json:"cluster_vote_policy"`
}

// ClusterDeviceMember 
type ClusterDeviceMember struct {
	Address     string  `json:"address"`
	ActivePower sdk.Dec `json:"active_power"` 
}

// ClusterPowerMember 
type ClusterPowerMember struct {
	Address     string  `json:"address"`
	ActivePower sdk.Dec `json:"active_power"` 
	BurnAmount  sdk.Dec `json:"burn_amount"`  
}

// PowerSupply 
type PowerSupply struct {
	ActivePower sdk.Dec `json:"active_power"` 
}

// PersonalClusterInfo 
type PersonalClusterInfo struct {
	Address string `json:"address"`
	
	Device map[string]struct{} `json:"device"`
	
	Owner map[string]struct{} `json:"owner"`
	
	BePower map[string]struct{} `json:"be_power"`
	
	AllBurn sdk.Dec `json:"all_burn"`
	// （）
	ActivePower sdk.Dec `json:"active_power"`
	
	FreezePower sdk.Dec `json:"freeze_power"`
	
	FirstPowerCluster string `json:"first_power_cluster"`
}

// FreezeUsers 
type FreezeUsers map[string]sdk.Dec


type ApprovePower struct {
	ClusterId string `json:"cluster_id"` //id
	EndBlock  int64  `json:"end_block"`
}

type ClusterCurApprove struct {
	ApproveAddress string `json:"approve_address"` //id
	EndBlock       int64  `json:"end_block"`
}
