package contract

import (
	"freemasonry.cc/blockchain/x/contract/keeper"
	"freemasonry.cc/blockchain/x/contract/types"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewContractProposalHandler(k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.ContractProposal:
			return handleContractProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized param proposal content type: %T", c)
		}
	}
}

func handleContractProposal(ctx sdk.Context, k *keeper.Keeper, p *types.ContractProposal) error {
	for _, val := range p.Contract {
		address := common.HexToAddress(val.ContractAddress)
		isExist := k.QueryContractIsExist(ctx, address)
		if !isExist {
			return types.ErrContractNotFound
		}
	}
	return nil
}
