package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// create a new ClusterHistoricalRewards
func NewClusterHistoricalRewards(cumulativeRewardRatio sdk.DecCoins, referenceCount uint32) ClusterHistoricalRewards {
	return ClusterHistoricalRewards{
		CumulativeRewardRatio: cumulativeRewardRatio,
		ReferenceCount:        referenceCount,
	}
}
func NewClusterHistoricalRewardsWithHis(cumulativeRewardRatio, HisReward sdk.DecCoins, referenceCount uint32) ClusterHistoricalRewards {
	return ClusterHistoricalRewards{
		CumulativeRewardRatio: cumulativeRewardRatio,
		ReferenceCount:        referenceCount,
		HisReward:             HisReward,
	}
}

// create a new ClusterCurrentRewards
func NewClusterCurrentRewards(rewards sdk.DecCoins, period uint64) ClusterCurrentRewards {
	return ClusterCurrentRewards{
		Rewards: rewards,
		Period:  period,
	}
}

func NewBurnStartingInfo(previousPeriod uint64, stake sdk.Dec, height uint64) BurnStartingInfo {
	return BurnStartingInfo{
		PreviousPeriod: previousPeriod,
		Stake:          stake,
		Height:         height,
	}
}
