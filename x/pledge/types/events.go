package types

// staking module event types
const (
	//staking events

	EventTypeCompleteUnbonding    = "complete_unbonding"
	EventTypeCompleteRedelegation = "complete_redelegation"
	EventTypeCreateValidator      = "create_validator"
	EventTypeEditValidator        = "edit_validator"
	EventTypeDelegate             = "chat_delegate"
	EventTypeUnbond               = "chat_unbond"
	EventTypeRedelegate           = "redelegate"

	AttributeKeyValidator         = "validator"
	AttributeKeyCommissionRate    = "commission_rate"
	AttributeKeyMinSelfDelegation = "min_self_delegation"
	AttributeKeySrcValidator      = "source_validator"
	AttributeKeyDstValidator      = "destination_validator"
	AttributeKeyDelegator         = "delegator"
	AttributeKeyCompletionTime    = "completion_time"
	AttributeKeyNewShares         = "new_shares"
	AttributeValueCategory        = ModuleName

	EventTypeSlashAmount         = "slash_amount"
	AttributeKeyDelegatorAddr    = "delegatorAddress"
	AttributeKeyDelegatorBalance = "delegator_balance"

	//mint module event types

	EventTypeMint = ModuleName

	AttributeKeyBondedRatio      = "bonded_ratio"
	AttributeKeyInflation        = "inflation"
	AttributeKeyAnnualProvisions = "annual_provisions"

	//distribution module event types

	EventTypeRewards            = "rewards"
	EventTypeCommission         = "commission"
	EventTypeWithdrawRewards    = "withdraw_rewards"
	EventTypeWithdrawCommission = "withdraw_commission"
	EventTypeProposerReward     = "proposer_reward"

	AttributeKeyWithdrawAddress = "withdraw_address"

	// 

	EventPledgeFromAddress = "pledge_from_address"
	EventPledgeToAddress   = "pledge_to_address"
	EventPledgeAmount      = "pledge_amount"
	EventPledgeDenom       = "pledge_denom"
	EventPledgeFromBalance = "pledge_from_balance"

	//

	EventUnPledgeFromAddress = "pledge_from_address"
	EventUnPledgeToAddress   = "pledge_to_address"
	EventUnPledgeAmount      = "pledge_amount"
	EventUnPledgeDenom       = "pledge_denom"
	EventUnPledgeFromBalance = "pledge_from_balance"

	// 
	ChatWithDrawAllEventType            = "withdraw_delegator_reward_all"
	ChatWithDrawAllEventTypeFromAddress = "from_address"
	ChatWithDrawAllEventTypeFromBalance = "from_balance"

	// 
	ChatWithDrawEventType              = "withdraw_delegator_reward"
	ChatWithDrawEventTypeFromAddress   = "from_address"
	ChatWithDrawEventTypeFromBalance   = "from_balance"
	ChatWithDrawEventTypeAmount        = "amount"
	ChatWithDrawEventTypeDenom         = "denom"
	ChatWithDrawEventTypeModuleAddress = "module_address"
)
