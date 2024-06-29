package types

// distribution module event types
const (
	
	AttributeSendeer           = "sender"
	AttributeSenderBalances    = "sender_balances"
	AttributeKeyClusterChatId  = "cluster_chat_id"
	AttributeKeyClusterOwner   = "cluster_owner"
	AttributeKeyClusterName    = "cluster_name"
	AttributeKeyFeeSingle      = "fee_payer"
	AttributeKeyMemo           = "memo"
	AttributeKeySenderBalances = "sender_balances"

	
	EventTypeCreateCluster = "create_cluster"

	
	EventTypeClusterUpgrade = "cluster_upgrade"
	AttributeKeyOldLevel    = "cluster_old_level"
	AttributeKeyNewLevel    = "cluster_new_level"

	
	EventTypeDeleteMembers = "cluster_delete_members"

	
	EventTypeChangeName = "cluster_change_name"

	
	EventTypeClusterExit = "cluster_exit"

	
	EventTypeWithdrawClusterRewards     = "withdraw_cluster_rewards"
	EventTypeWithdrawDeviceRewards      = "withdraw_device_rewards"
	EventTypeIncrementHistoricalRewards = "increment_historical_rewards"

	AttributeKeyCluster    = "cluster"
	AttributeKeyMember     = "member"
	AttributeKeyTime       = "time"
	AttributeValueCategory = ModuleName

	
	EventTypeAddMembers = "add_members"
	AttributeClusterId  = "cluster_id"

	EventTypeDeductionFee = "deduction"
	AttributeDeductionFee = "deduction_fee"
)
