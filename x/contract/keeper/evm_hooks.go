package keeper

import (
    "freemasonry.cc/blockchain/contracts"
    "freemasonry.cc/blockchain/x/contract/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/ethereum/go-ethereum/core"
    ethtypes "github.com/ethereum/go-ethereum/core/types"
    evmtypes "github.com/evmos/ethermint/x/evm/types"
    "math/big"
)

// Hooks wrapper struct for slashing keeper
type Hooks struct {
    k Keeper
}

var _ evmtypes.EvmHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
    return Hooks{k}
}

func (h Hooks) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
    params := h.k.GetParams(ctx)
    freeMasonryMedal := contracts.FreeMasonryMedalJSONContract.ABI
    for _, log := range receipt.Logs {
        if len(log.Topics) != 2 {
            continue
        }
        eventID := log.Topics[0]
        event, err := freeMasonryMedal.EventByID(eventID)
        if err != nil {
            continue
        }
        if event.Name != types.NftEventCreate {
            h.k.Logger(ctx).Info("emitted event", "name", event.Name, "signature", event.Sig)
            continue
        }
        nftEvent, err := freeMasonryMedal.Unpack(event.Name, log.Data)
        if err != nil {
            h.k.Logger(ctx).Error("failed to unpack NFT event", "error", err.Error())
            return err
        }

        level, ok := nftEvent[1].(*big.Int)
        if !ok {
            h.k.Logger(ctx).Error("failed to conversion NFT level")
            return types.ErrNftlevel
        }
        addrByte := log.Topics[1].Bytes()
        address := sdk.AccAddress(addrByte[len(addrByte)-20:])
        createTime := nftEvent[2].(*big.Int)
        
        delegateInfo, err := h.k.pledgeKeeper.GetDelegateTime(ctx, address)
        if err != nil {
            return err
        }
        if delegateInfo == nil {
            return types.ErrDelegate
        }
        
        accountLevel, err := h.k.pledgeKeeper.QueryPledgeLevelByAccAddress(ctx, address)
        if err != nil {
            return err
        }

        
        if level.Int64() > accountLevel.Level {
            return types.ErrDelegateLevel
        }

        
        if createTime.Int64()-delegateInfo.Time < params.Days {
            return types.ErrDelegateTime
        }
        
        err = h.k.pledgeKeeper.UpdateDelegateTime(ctx, address)
        if err != nil {
            return err
        }
    }
    return nil
}
