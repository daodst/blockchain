package client

import (
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

func TestCreateCluster(t *testing.T) {
	txClient := NewTxClient()
	clusterClient := NewClusterClient(&txClient)

	accClient := NewAccountClient(&txClient)

	from, err := accClient.CreateAccountFromSeed("fly crawl forest insect report bargain office shuffle rifle art scan grid")

	fee := core.NewLedgerFee(9)

	_, result, err := clusterClient.CreateCluster(
		from.Address,
		fee,
		"dstvaloper10x9a5ym9hunn9splh0y3pfhakc89hxnvjwan6z",
		"!!szyyUvrLcY6bz7rf:1111111.nxn",
		"FFF",
		"dst1fttgy363wn8z4pxeqfp4ck63gjrvgmxwj579vr",
		sdk.MustNewDecFromStr("0.5"),
		sdk.MustNewDecFromStr("0.5"),
		sdk.MustNewDecFromStr("7000000000000000000000"),
		from.PrivateKey,
	)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("data:", result.Code)
}

func TestGetParams(t *testing.T) {
	txClient := NewTxClient()
	clusterClient := NewClusterClient(&txClient)

	res, err := clusterClient.QueryDaoParams()
	t.Log("BurnGetPowerRatio", res.BurnGetPowerRatio)
	t.Log("CreateClusterMinBurn", res.CreateClusterMinBurn)
	t.Log("BurnAddress", res.BurnAddress)
	t.Log("DayBurnReward", res.DayBurnReward)
	t.Log("DeviceRange", res.DeviceRange)
	t.Log("SalaryRange", res.SalaryRange)

	t.Log(err)
}
