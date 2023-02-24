package types

import (
    "context"
    "freemasonry.cc/blockchain/x/pledge/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core"
    "github.com/ethereum/go-ethereum/core/vm"
    "github.com/evmos/ethermint/x/evm/statedb"
    evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type PledgeKeeper interface {
    GetDelegateTime(ctx sdk.Context, address sdk.AccAddress) (*types.DelegateTime, error)
    UpdateDelegateTime(ctx sdk.Context, address sdk.AccAddress) error
    QueryPledgeLevelByAccAddress(ctx sdk.Context, accAddress sdk.AccAddress) (types.PledgeLevel, error)
}

type EVMKeeper interface {
    GetParams(ctx sdk.Context) evmtypes.Params
    GetAccountWithoutBalance(ctx sdk.Context, addr common.Address) *statedb.Account
    GetCode(ctx sdk.Context, codeHash common.Hash) []byte
    EstimateGas(c context.Context, req *evmtypes.EthCallRequest) (*evmtypes.EstimateGasResponse, error)
    ApplyMessage(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*evmtypes.MsgEthereumTxResponse, error)
}
