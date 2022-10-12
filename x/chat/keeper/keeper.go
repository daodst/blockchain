package keeper

import (
	"fmt"
	"freemasonry.cc/blockchain/cmd/config"
	"freemasonry.cc/blockchain/core"
	commkeeper "freemasonry.cc/blockchain/x/comm/keeper"
	types2 "freemasonry.cc/blockchain/x/comm/types"
	pledgekeeper "freemasonry.cc/blockchain/x/pledge/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	"strconv"
	"strings"

	"freemasonry.cc/blockchain/x/chat/types"
)

// Keeper of this module maintains collections of erc20.
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	commKeeper    commkeeper.Keeper
	pledgeKeeper  pledgekeeper.Keeper
}

// NewKeeper creates new instances of the erc20 Keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	cm commkeeper.Keeper,
	pk pledgekeeper.Keeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramstore:    ps,
		accountKeeper: ak,
		bankKeeper:    bk,
		commKeeper:    cm,
		pledgeKeeper:  pk,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) KVHelper(ctx sdk.Context) storeHelper {
	store := ctx.KVStore(k.storeKey)
	return storeHelper{
		store,
	}
}

func (k Keeper) GetRegisterInfo(ctx sdk.Context, fromAddress string) (userinfo types.UserInfo, err error) {

	store := k.KVHelper(ctx)
	key := types.KeyPrefixRegisterInfo + fromAddress
	if store.Has(key) {

		err := store.GetUnmarshal(key, &userinfo)
		if err != nil {
			return types.UserInfo{}, err
		}

		return userinfo, nil

	} else {

		return types.UserInfo{}, types.ErrUserNotFound

	}
}

func (k Keeper) SetRegisterInfo(ctx sdk.Context, userInfo types.UserInfo) error {

	store := k.KVHelper(ctx)
	key := types.KeyPrefixRegisterInfo + userInfo.FromAddress

	err := store.Set(key, userInfo)
	if err != nil {
		return err
	}

	return nil
}

// ChatPoundage ，（coin）
func (k Keeper) ChatPoundage(ctx sdk.Context, accFromAddress, nodeAddress sdk.AccAddress, distCoin sdk.Coin) (sdk.Dec, error) {

	//Dec
	AllCoinDec := distCoin.Amount.ToDec()

	//
	chatParams := k.pledgeKeeper.GetParams(ctx)

	//
	destroyDec := AllCoinDec.Mul(chatParams.AttDestroyPercent)
	//
	gatewayDec := AllCoinDec.Mul(chatParams.AttGatewayPercent)
	//dpos
	DPosDec := AllCoinDec.Mul(chatParams.AttDposPercent)
	//
	feeAllInt := destroyDec.TruncateInt().Add(gatewayDec.TruncateInt()).Add(DPosDec.TruncateInt())
	feeAllDec := feeAllInt.ToDec()

	feeCoin := sdk.NewCoin(distCoin.Denom, feeAllInt)

	// 
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, accFromAddress, types.ModuleName, sdk.NewCoins(feeCoin))
	if err != nil {
		return sdk.Dec{}, types.ErrPledgeFeeTransfer
	}

	//
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(distCoin.Denom, destroyDec.TruncateInt())))
	if err != nil {
		return sdk.Dec{}, types.ErrPledgeFeeBurn
	}
	//
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, nodeAddress, sdk.NewCoins(sdk.NewCoin(distCoin.Denom, gatewayDec.TruncateInt())))
	if err != nil {
		return sdk.Dec{}, types.ErrPledgeFeeGateway
	}
	//dpos
	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, sdk.NewCoins(sdk.NewCoin(distCoin.Denom, DPosDec.TruncateInt())))
	if err != nil {
		return sdk.Dec{}, types.ErrPledgeFeeDpos
	}

	return feeAllDec, nil
}

//
func (k Keeper) RegisterMobile(ctx sdk.Context, nodeAddress, fromAddress, mobilePrefix string) (mobile string, err error) {

	//log := core.BuildLog(core.GetFuncName(), core.LmChainChatKeeper)

	//
	GatewayNumInfo, _, err := k.commKeeper.GetGatewayNum(ctx, mobilePrefix)

	if err != nil {
		return mobile, types.ErrGeneratingMobile
	}

	if GatewayNumInfo == nil {
		return mobile, types.ErrGetMobile
	}

	if GatewayNumInfo.GatewayAddress != nodeAddress || GatewayNumInfo.Status != 0 {
		return mobile, types.ErrGateway
	}

	//
	mobileSuffixInt := len(GatewayNumInfo.NumberEnd)

	if mobileSuffixInt >= types.MobileSuffixMax {
		return mobile, types.ErrGetMobile
	}

	//
	mobileSuffixString := strconv.Itoa(mobileSuffixInt)

	//
	mobileSuffixString = strings.Repeat("0", types.MobileSuffixLength-len(mobileSuffixString)) + mobileSuffixString

	mobile = mobilePrefix + mobileSuffixString

	GatewayNumInfo.NumberEnd = append(GatewayNumInfo.NumberEnd, mobile)
	GatewayNumInfos := []types2.GatewayNumIndex{
		*GatewayNumInfo,
	}
	//
	err = k.commKeeper.SetGatewayNum(ctx, GatewayNumInfos)
	if err != nil {
		return mobile, types.ErrMobileSetError
	}
	//
	err = k.commKeeper.UpdateGatewayNum(ctx, GatewayNumInfos)
	if err != nil {
		return mobile, err
	}
	return mobile, nil
}

// CaculateChatReward 
/*
	var newAmountDec = 

	for  {
		for ， ratio
		for ，（） mortgageAmount

		//，
		 =  + （ * ）
		newAmountDec = newAmountDec + mortgageAmount * ratio
	}

	（Dec -> sdk.Coin）
	return newAmountCoin
*/
func (k Keeper) CaculateChatReward(ctx sdk.Context, fromAddress string, chatRewardChangeLog []types.ChatReward, canRedemAmount sdk.Coin) (reward sdk.Coin, err error) {
	logs := core.BuildLog(core.GetFuncName(), core.LmChainChatKeeper)
	//
	mortgageLog, err := k.GetMortgageLog(ctx, fromAddress)
	if err != nil {
		return reward, err
	}

	logs.Info(":", mortgageLog)

	//(comm)
	commParams := k.commKeeper.GetParams(ctx)

	HeightPerDay := commParams.BonusCycle

	logs.Info(":", HeightPerDay)

	//
	LastGetInfo, err := k.GetLastGetHeight(ctx, fromAddress)
	if err != nil {
		return reward, err
	}

	logs.Info(":", LastGetInfo)

	//
	NowHeight := ctx.BlockHeight()

	logs.Info(":", NowHeight)

	newAmountDec := canRedemAmount.Amount.ToDec()
	//LastGetInfo.Height()  NowHeight（） ， HeightPerDay（）
	//
	for SendBonusHeight := GetCommonMultipleGt(LastGetInfo.Height, HeightPerDay); SendBonusHeight < NowHeight; SendBonusHeight += HeightPerDay {

		//，
		var ratio sdk.Dec
		for changeIndex, chatReward := range chatRewardChangeLog {

			//
			if SendBonusHeight >= chatReward.Height {
				if HasIndex(changeIndex+1, chatRewardChangeLog) { //
					if SendBonusHeight >= chatRewardChangeLog[changeIndex+1].Height { //，
						continue
					} else {
						ratio, err = sdk.NewDecFromStr(chatReward.Value) //，
						if err != nil {
							return reward, types.ErrGetBonus
						}
					}
				} else { //,
					ratio, err = sdk.NewDecFromStr(chatReward.Value)
					if err != nil {
						return reward, types.ErrGetBonus
					}
				}
			}

		}

		//， （SendBonusHeight）
		mortgageAmount := sdk.NewCoin(mortgageLog[0].MortgageValue.Denom, sdk.ZeroInt()) //0
		for _, addLog := range mortgageLog {

			//， ，
			if addLog.Height < SendBonusHeight-HeightPerDay {
				mortgageAmount = addLog.MortgageValue
			} else {
				break
			}
		}

		logs.Info("-----------------")
		logs.Info("---:", SendBonusHeight)
		logs.Info("---:", ratio)
		logs.Info("---:", mortgageAmount)
		logs.Info("---:", newAmountDec)

		// =  +  * 
		newAmountDec = newAmountDec.Add(mortgageAmount.Amount.ToDec().Mul(ratio))
		logs.Info("---:", newAmountDec)
		logs.Info("-----------------")
	}

	logs.Info("：", newAmountDec)
	return sdk.NewCoin(LastGetInfo.Value.Denom, newAmountDec.TruncateInt()), nil
}

// a  b 
func GetCommonMultipleGt(a, b int64) int64 {
	return (a/b + 1) * b
}

//
func HasIndex(index int, data []types.ChatReward) bool {
	return len(data) > index
}

func (k Keeper) GetLastGetHeight(ctx sdk.Context, fromAddress string) (types.LastReceiveLog, error) {
	store := k.KVHelper(ctx)
	key := types.KeyPrefixLastGetRewardLog + fromAddress

	logs := types.LastReceiveLog{
		Height: 1,
		Value:  sdk.NewCoin(config.BaseDenom, sdk.NewInt(0)),
	}

	if store.Has(key) {

		err := store.GetUnmarshal(key, &logs)
		if err != nil {
			return logs, types.ErrGetLastReveiveHeight
		}
		return logs, nil
	}

	return logs, nil
}

func (k Keeper) SetLastGetHeight(ctx sdk.Context, fromAddress string, height int64, mortgage sdk.Coin) error {

	store := k.KVHelper(ctx)
	key := types.KeyPrefixLastGetRewardLog + fromAddress

	err := store.Set(key, types.LastReceiveLog{
		Height: height,
		Value:  mortgage,
	})

	if err != nil {
		return types.ErrSetLastReveiveHeight
	}

	//log := types.LastReceiveLog{
	//	Height: 1,
	//	Value:  sdk.NewCoin(config.BaseDenom, sdk.NewInt(0)),
	//}
	//store.GetUnmarshal(key, &log)
	//
	//fmt.Println("log:", log)

	return nil
}

func (k Keeper) SetMortgageLog(ctx sdk.Context, fromAddress string, mortgageNew sdk.Coin) error {

	//logs := core.BuildLog(core.GetFuncName(), core.LmChainChatKeeper)

	store := k.KVHelper(ctx)
	key := types.KeyPrefixMortgageAddLog + fromAddress

	logNew := types.MortgageAddLog{
		Height:        ctx.BlockHeight(),
		MortgageValue: mortgageNew,
	}
	logs := make([]types.MortgageAddLog, 0)

	if store.Has(key) {
		err := store.GetUnmarshal(key, &logs)
		if err != nil {
			return types.ErrSetMortgageLog
		}

		//
		mortgageAddLogCheck := k.CheckMortgageAddLog(ctx, logs, logNew)
		if !mortgageAddLogCheck {
			return types.ErrSetMortgageLog
		}

	}

	logs = append(logs, logNew)
	err := store.Set(key, logs)
	if err != nil {
		return types.ErrSetMortgageLog
	}

	return nil
}

func (k Keeper) GetMortgageLog(ctx sdk.Context, fromAddress string) ([]types.MortgageAddLog, error) {
	store := k.KVHelper(ctx)
	key := types.KeyPrefixMortgageAddLog + fromAddress

	logs := make([]types.MortgageAddLog, 0)
	if store.Has(key) {
		err := store.GetUnmarshal(key, &logs)
		if err != nil {
			return logs, types.ErrGetMortgageLog
		}
	} else {
		return logs, types.ErrUserNotFound
	}
	return logs, nil
}

func (k Keeper) CheckMortgageAddLog(ctx sdk.Context, log []types.MortgageAddLog, logNew types.MortgageAddLog) bool {

	lastMortgageInfo := log[len(log)-1]
	//
	if lastMortgageInfo.Height >= logNew.Height {
		return false
	}

	//
	if lastMortgageInfo.MortgageValue.Amount.GT(logNew.MortgageValue.Amount) {
		return false
	}

	return true
}
