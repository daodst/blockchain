package keeper

import (
    "encoding/json"
    "fmt"
    "freemasonry.cc/blockchain/contracts"
    "freemasonry.cc/blockchain/util"
    "freemasonry.cc/blockchain/x/contract/types"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
    stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/hexutil"
    ethtypes "github.com/ethereum/go-ethereum/core/types"
    "github.com/evmos/ethermint/server/config"
    evmtypes "github.com/evmos/ethermint/x/evm/types"
    "github.com/tendermint/tendermint/libs/log"
    "math/big"
    "strings"
)

// Keeper of this module maintains collections of erc20.
type Keeper struct {
    storeKey   sdk.StoreKey
    cdc        codec.BinaryCodec
    paramstore paramtypes.Subspace

    stakingKeeper *stakingKeeper.Keeper
    accountKeeper types.AccountKeeper
    bankKeeper    types.BankKeeper
    pledgeKeeper  types.PledgeKeeper
    evmKeeper     types.EVMKeeper
}

// NewKeeper creates new instances of the erc20 Keeper
func NewKeeper(
    storeKey sdk.StoreKey,
    cdc codec.BinaryCodec,
    ps paramtypes.Subspace,
    ak types.AccountKeeper,
    bk types.BankKeeper,
    stakingKeeper *stakingKeeper.Keeper,
    pledgeKeeper types.PledgeKeeper,
    evmKeeper types.EVMKeeper,
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
        stakingKeeper: stakingKeeper,
        pledgeKeeper:  pledgeKeeper,
        evmKeeper:     evmKeeper,
    }
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
    return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) KVHelper(ctx sdk.Context) StoreHelper {
    store := ctx.KVStore(k.storeKey)
    return StoreHelper{
        store,
    }
}


func (k Keeper) QueryContractIsExist(ctx sdk.Context, address common.Address) bool {
    acct := k.evmKeeper.GetAccountWithoutBalance(ctx, address)
    var code []byte
    if acct != nil && acct.IsContract() {
        code = k.evmKeeper.GetCode(ctx, common.BytesToHash(acct.CodeHash))
    }
    if len(code) == 0 {
        return false
    }
    return true
}

//NFT
func (k Keeper) SetNftContractAddress(ctx sdk.Context, contract string) error {
    store := k.KVHelper(ctx)
    err := store.Set(types.KeyNftContractAddress, contract)
    if err != nil {
        return err
    }
    return nil
}

func (k Keeper) GetNftContractAddress(ctx sdk.Context) string {
    store := k.KVHelper(ctx)
    return string(store.Get(types.KeyNftContractAddress))
}

//NFT
func (k Keeper) GetNftInfo(ctx sdk.Context, address, contract common.Address) (nftInfo []types.NftInfo, err error) {
    freeMasonryMedal := contracts.FreeMasonryMedalJSONContract.ABI
    resp, err := k.CallEVM(ctx, freeMasonryMedal, address, contract, false, "getUserNfts", address)
    if err != nil {
        return
    }
    data := make(map[string]interface{})
    if err = freeMasonryMedal.UnpackIntoMap(data, "getUserNfts", resp.Ret); err != nil {
        return
    }
    nftByte, err := util.Json.Marshal(data["tokens"])
    if err != nil {
        return
    }
    err = util.Json.Unmarshal(nftByte, &nftInfo)
    if err != nil {
        return
    }
    return
}

// CallEVM performs a smart contract method call using given args
func (k Keeper) CallEVM(
    ctx sdk.Context,
    abi abi.ABI,
    from, contract common.Address,
    commit bool,
    method string,
    args ...interface{},
) (*evmtypes.MsgEthereumTxResponse, error) {
    data, err := abi.Pack(method, args...)
    if err != nil {
        return nil, sdkerrors.Wrap(
            types.ErrABIPack,
            sdkerrors.Wrap(err, "failed to create transaction data").Error(),
        )
    }

    resp, err := k.CallEVMWithData(ctx, from, &contract, data, commit)
    if err != nil {
        return nil, sdkerrors.Wrapf(err, "contract call failed: method '%s', contract '%s'", method, contract)
    }
    return resp, nil
}

// CallEVMWithData performs a smart contract method call using contract data
func (k Keeper) CallEVMWithData(
    ctx sdk.Context,
    from common.Address,
    contract *common.Address,
    data []byte,
    commit bool,
) (*evmtypes.MsgEthereumTxResponse, error) {
    nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
    if err != nil && !strings.Contains(err.Error(), "does not exist") {
        return nil, err
    }
    gasCap := config.DefaultGasCap
    if commit {
        args, err := json.Marshal(evmtypes.TransactionArgs{
            From: &from,
            To:   contract,
            Data: (*hexutil.Bytes)(&data),
        })
        if err != nil {
            return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "failed to marshal tx args: %s", err.Error())
        }

        gasRes, err := k.evmKeeper.EstimateGas(sdk.WrapSDKContext(ctx), &evmtypes.EthCallRequest{
            Args:   args,
            GasCap: config.DefaultGasCap,
        })
        if err != nil {
            return nil, err
        }
        gasCap = gasRes.Gas
    }
    msg := ethtypes.NewMessage(
        from,
        contract,
        nonce,
        big.NewInt(0), // amount
        gasCap,        // gasLimit
        big.NewInt(0), // gasFeeCap
        big.NewInt(0), // gasTipCap
        big.NewInt(0), // gasPrice
        data,
        ethtypes.AccessList{}, // AccessList
        !commit,               // isFake
    )
    res, err := k.evmKeeper.ApplyMessage(ctx, msg, evmtypes.NewNoOpTracer(), commit)
    if err != nil {
        return nil, err
    }

    if res.Failed() {
        return nil, sdkerrors.Wrap(evmtypes.ErrVMExecution, res.VmError)
    }

    return res, nil
}
