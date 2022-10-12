package pledge

import (
	"errors"
	"freemasonry.cc/blockchain/cmd/config"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/pledge/keeper"
	"freemasonry.cc/blockchain/x/pledge/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewPledgeDelegateProposalHandler(k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.PledgeDelegateProposal:
			return handlePledgeDelegateProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized param proposal content type: %T", c)
		}
	}
}

func handlePledgeDelegateProposal(ctx sdk.Context, k *keeper.Keeper, p *types.PledgeDelegateProposal) error {
	store := k.KVHelper(ctx)
	//,
	if store.Has(types.PledgeDelegateKey) {
		return types.ErrProposalDelegate
	}
	result := make(map[string]types.PledgeDelegate)
	for _, val := range p.Delegate {
		amInt := core.MustRealString2LedgerIntNoMin(val.Amount)
		coin := sdk.NewCoin(config.BaseDenom, amInt)
		err := k.MintCoins(ctx, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
		addr, err := sdk.AccAddressFromBech32(val.Address)
		if err != nil {
			return err
		}
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
		msg := types.NewMsgPledge(val.Address, val.Address, val.GatewayAddress, coin)

		accFromAddress, err := sdk.AccAddressFromBech32(msg.FromAddress)
		if err != nil {
			return err
		}

		valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddress)
		if valErr != nil {
			return valErr
		}

		validator, found := k.GetValidator(ctx, valAddr)
		if !found {
			return errors.New("validator not found")
		}

		_, err = k.Delegate(ctx, accFromAddress, accFromAddress, msg.Amount.Amount, validator)
		if err != nil {
			return err
		}
		result[val.Address] = val
	}
	err := k.SetPledgeDelegate(ctx, result)
	if err != nil {
		return err
	}
	return nil
}
