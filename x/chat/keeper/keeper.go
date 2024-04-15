package keeper

import (
	"encoding/json"
	"fmt"
	"freemasonry.cc/blockchain/core"
	types3 "freemasonry.cc/blockchain/x/contract/types"
	gatewaykeeper "freemasonry.cc/blockchain/x/gateway/keeper"
	types2 "freemasonry.cc/blockchain/x/gateway/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	"strconv"
	"strings"

	"freemasonry.cc/blockchain/x/chat/types"
)

// Keeper of this module maintains collections of erc20.
type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	commKeeper    gatewaykeeper.Keeper
}

// NewKeeper creates new instances of the erc20 Keeper
func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	cm gatewaykeeper.Keeper,
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


func (k Keeper) Register(ctx sdk.Context, regAddr, regChatAddr, gateawyAddr string, getDid []string) error {

	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainMsgServer)

	if _, err := sdk.AccAddressFromBech32(regAddr); err != nil {
		return core.ErrAddressFormat
	}

	
	gatewayInfo, err := k.commKeeper.GetGatewayInfo(ctx, gateawyAddr)
	if err != nil {
		logs.WithError(err).Error("chat.Register GetGatewayInfo err")
		return core.ErrGetGatewauInfo
	}

	
	if gatewayInfo.Status == 1 {
		logs.Error("chat.Register Validator Status Err")
		return core.ErrValidatorStatusError
	}

	
	userInfo, err := k.GetRegisterInfo(ctx, regAddr)
	if err != core.ErrUserNotFound {
		logs.WithError(err).Error("Register GetRegisterInfo already exist err")
		return err
	}

	
	userInfo.NodeAddress = gateawyAddr
	userInfo.RegisterNodeAddress = gateawyAddr
	userInfo.FromAddress = regAddr

	if len(getDid) != 0 {
		//todo did

		//userInfo.Mobile = append(userInfo.Mobile, getDid)
	}

	err = k.SetRegisterInfo(ctx, userInfo)
	if err != nil {
		logs.WithError(err).Error("SetRegisterInfo err")
		return core.ErrRegister
	}

	
	err = k.SetChatAddr(ctx, regChatAddr, regAddr)
	if err != nil {
		return err
	}

	
	ctx.EventManager().EmitEvents(
		[]sdk.Event{
			
			sdk.NewEvent(
				types.EventTypeRegister,
				
				sdk.NewAttribute(types.EventTypeRegAddress, regAddr),
				
				sdk.NewAttribute(types.EventTypeGatewayAddress, gateawyAddr),
				
				sdk.NewAttribute(types.EventPrefixMobile, gatewayInfo.GatewayNum[0].NumberIndex),
				
				sdk.NewAttribute(types.EventTypeChatAddress, regChatAddr),
			),
		},
	)

	return nil
}

func (k Keeper) GetRegisterInfo(ctx sdk.Context, fromAddress string) (userinfo types.UserInfo, err error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainChatKeeper)

	store := k.KVHelper(ctx)
	key := types.KeyPrefixRegisterInfo + fromAddress
	if store.Has(key) {
		err := store.GetUnmarshal(key, &userinfo)
		if err != nil {
			logs.WithError(err).Error("GetUnmarshal")
			return types.UserInfo{}, err
		}

		return userinfo, nil

	} else {
		logs.WithError(core.ErrUserNotFound).WithField("from:", fromAddress).Error("GetRegisterInfo")
		return types.UserInfo{}, core.ErrUserNotFound

	}
}

func (k Keeper) SetRegisterInfo(ctx sdk.Context, userInfo types.UserInfo) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainChatKeeper)

	if userInfo.FromAddress == "" {
		return core.ErrUserUpdate
	}

	store := k.KVHelper(ctx)

	key := types.KeyPrefixRegisterInfo + userInfo.FromAddress

	err := store.Set(key, userInfo)
	if err != nil {
		logs.WithError(err).WithField("from:", userInfo.FromAddress).Error("SetRegisterInfo Err")
		return core.ErrUserUpdate
	}

	return nil
}


func (k Keeper) RegisterMobile(ctx sdk.Context, nodeAddress, fromAddress, mobilePrefix string) (mobile string, err error) {

	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainChatKeeper)

	
	GatewayNumInfo, _, err := k.commKeeper.GetGatewayNum(ctx, mobilePrefix)

	if err != nil {
		logs.WithError(err).Info("Err GetGatewayNum: err")
		return mobile, core.ErrNumberOfGateWay
	}

	if GatewayNumInfo == nil {
		logs.WithError(core.ErrNumberOfGateWay).Info("Err GetGatewayNum: nil")
		return mobile, core.ErrGetMobile
	}

	if GatewayNumInfo.GatewayAddress != nodeAddress || GatewayNumInfo.Status != 0 {
		logs.WithError(core.ErrNumberOfGateWay).Info("Err GetGatewayNum: status or address")
		return mobile, core.ErrGateway
	}

	
	mobileSuffixInt := len(GatewayNumInfo.NumberEnd)

	if mobileSuffixInt >= types.MobileSuffixMax {
		logs.WithError(core.ErrGetMobile).Info("Err MobileSuffixMax:", mobileSuffixInt)
		return mobile, core.ErrMobileexhausted
	}

	
	mobileSuffixString := strconv.Itoa(mobileSuffixInt)

	
	mobileSuffixString = strings.Repeat("0", types.MobileSuffixLength-len(mobileSuffixString)) + mobileSuffixString

	mobile = mobilePrefix + mobileSuffixString

	GatewayNumInfo.NumberEnd = append(GatewayNumInfo.NumberEnd, mobile)
	GatewayNumInfos := []types2.GatewayNumIndex{
		*GatewayNumInfo,
	}
	
	err = k.commKeeper.SetGatewayNum(ctx, GatewayNumInfos)
	if err != nil {
		logs.WithError(err).Info("Err SetGatewayNum")
		return mobile, core.ErrMobileSetError
	}
	
	err = k.commKeeper.UpdateGatewayNum(ctx, GatewayNumInfos)
	if err != nil {
		logs.WithError(err).Info("Err UpdateGatewayNum")
		return mobile, err
	}

	
	err = k.SetMobileOwner(ctx, mobile, fromAddress)
	if err != nil {
		logs.WithError(err).Info("Err SetMobileOwner")
		return mobile, err
	}

	return mobile, nil
}

func (k Keeper) GetUserByMobile(ctx sdk.Context, mobile string) (string, error) {

	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainChatKeeper)
	store := k.KVHelper(ctx)
	key := types.KeyPrefixMobileOwner + mobile

	var fromAddress string
	if store.Has(key) {

		addressByte := store.Get(key)
		fromAddress = string(addressByte)
		if fromAddress == "" {
			logs.WithError(core.ErrUserNotFound).Error("GetUserByMobile")
			return "", core.ErrUserNotFound
		}
		return fromAddress, nil

	} else {
		logs.WithError(core.ErrUserNotFound).Error("GetUserByMobile")
		return fromAddress, core.ErrMobileNotFount

	}
}

func (k Keeper) SetMobileOwner(ctx sdk.Context, mobile, address string) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainChatKeeper)

	KeyMobileOwner := types.KeyPrefixMobileOwner + mobile
	store := k.KVHelper(ctx)
	err := store.Set(KeyMobileOwner, address)
	if err != nil {
		logs.WithError(core.ErrSetMobileOwner).Error("SetMobileOwner Err")
		return core.ErrSetMobileOwner
	}
	return nil
}

func (k Keeper) GetUserInfos(ctx sdk.Context, fromAddresses []string) (userInfos []types.AllUserInfo, err error) {
	//logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainChatKeeper)
	
	//store := k.KVHelper(ctx)
	
	//for _, address := range fromAddresses {
	//	accAddress, err := sdk.AccAddressFromBech32(address)
	//	if err != nil {
	//		return nil, err
	//	}
	
	//	key := types.KeyPrefixRegisterInfo + address
	//	if store.Has(key) {
	//		var userInfo types.UserInfo
	//		err := store.GetUnmarshal(key, &userInfo)
	//		if err != nil {
	//			return nil, err
	//		}
	
	//		
	//		var gatewayProfixMobile string
	//		gateWayInfo, err := k.commKeeper.GetGatewayInfo(ctx, userInfo.NodeAddress)
	//		if err != nil {
	//			logs.WithError(err).Error("GetGatewayInfo err")
	//			return nil, err
	//		}
	
	//		if gateWayInfo.Status == 1 {
	//			gatewayProfixMobile = ""
	//		}
	
	//		gatewayProfixMobile = gateWayInfo.GatewayNum[0].NumberIndex
	
	//		userInfos = append(userInfos, types.AllUserInfo{
	//			UserInfo:            userInfo,
	//			IsExist:             0,
	//			//PledgeLevel:         pledgeLevelInfo.Level,
	//			GatewayProfixMobile: gatewayProfixMobile,
	//		})
	//	} else {
	//		userInfos = append(userInfos, types.AllUserInfo{
	//			UserInfo: types.UserInfo{
	//				FromAddress: address,
	//			},
	//			IsExist: 1,
	//		})
	//	}
	//}
	return nil, nil
}

//   @dex1zkcz0qal5l60wmxk0fzvffkw32pyj7ze8tgwuk:1888858.fm
func (k Keeper) GetGatewayProfixMobiles(ctx sdk.Context, addresses []string) []types.CustomInfo {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainChatKeeper)

	res := make([]types.CustomInfo, 0)

	for _, address := range addresses {
		
		userMobile := ""
		userInfo, err := k.GetRegisterInfo(ctx, address)
		if err != nil {
			logs.WithError(err).Error("GetRegisterInfo err:" + address)
			res = append(res, types.CustomInfo{
				Address:     address,
				CommAddress: "",
				Mobile:      "",
			})
			continue
		}

		if len(userInfo.Mobile) > 0 {
			userMobile = userInfo.Mobile[0]
		}

		gatewayAddress := userInfo.NodeAddress
		gateWayInfo, err := k.commKeeper.GetGatewayInfo(ctx, gatewayAddress)
		if err != nil {
			logs.WithError(err).Error("GetGatewayInfo err" + gatewayAddress)
			res = append(res, types.CustomInfo{
				Address:     address,
				CommAddress: "",
				Mobile:      userMobile,
			})
			continue
		}

		if gateWayInfo.Status == 1 {
			logs.WithError(core.ErrGetGatewauInfo).Error("GetGateway Status Error" + gatewayAddress)
			res = append(res, types.CustomInfo{
				Address:     address,
				CommAddress: "",
				Mobile:      userMobile,
			})
			continue
		}

		res = append(res, types.CustomInfo{
			Address:     address,
			CommAddress: "@" + address + ":" + gateWayInfo.GatewayNum[0].NumberIndex + "." + core.GovDenom,
			Mobile:      userMobile,
		})

	}
	return res
}


func (k Keeper) SetGatewayIssueToken(ctx sdk.Context, gatewayAddress string, tokenInfo types3.GatewayTokenInfo) error {
	//log := core.BuildLog(core.GetPackageFuncName(), core.LmChainChatKeeper)
	store := ctx.KVStore(k.storeKey)

	key := types.KeyPrefixGatewayIssueToken + gatewayAddress

	data, _ := json.Marshal(tokenInfo)

	store.Set([]byte(key), data)

	return nil
}


func (k Keeper) GetGatewayIssueToken(ctx sdk.Context, gatewayAddress string) (*types3.GatewayTokenInfo, error) {
	//log := core.BuildLog(core.GetPackageFuncName(), core.LmChainChatKeeper)
	store := ctx.KVStore(k.storeKey)

	key := []byte(types.KeyPrefixGatewayIssueToken + gatewayAddress)

	if store.Has(key) {

		resp := types3.GatewayTokenInfo{}

		data := store.Get(key)

		err := json.Unmarshal(data, &resp)

		if err != nil {
			return nil, err
		}

		return &resp, nil

	}

	return nil, core.ErrTokenNotFound
}

func (k Keeper) SetChatAddr(ctx sdk.Context, chatAddress, fromAddress string) error {
	store := ctx.KVStore(k.storeKey)

	key := []byte(types.KeyChatAddress + chatAddress)

	if store.Has(key) {
		return core.ErrChatAddressExist
	}

	v := []byte(fromAddress)

	store.Set(key, v)

	return nil
}

func (k Keeper) GetAddrFromChatAddr(ctx sdk.Context, chatAddress string) (string, error) {
	store := ctx.KVStore(k.storeKey)

	key := []byte(types.KeyChatAddress + chatAddress)

	if !store.Has(key) {
		return "", core.ErrChatAddressNotExist
	}

	vb := store.Get(key)

	return string(vb), nil
}
