package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNftlevel              = sdkerrors.Register(ModuleName, 301, "conversion NFT level error")
	ErrDelegate              = sdkerrors.Register(ModuleName, 302, "not find delegate info")
	ErrDelegateTime          = sdkerrors.Register(ModuleName, 303, "Insufficient pledge time")
	ErrDelegateLevel         = sdkerrors.Register(ModuleName, 304, "Insufficient pledge level")
	ErrABIPack               = sdkerrors.Register(ModuleName, 305, "contract ABI pack failed")
	ErrNftTokenId            = sdkerrors.Register(ModuleName, 306, "conversion NFT tokenId error")
	ErrEmptyProposalContract = sdkerrors.Register(ModuleName, 307, "invalid proposal contract")
	ErrContractNotFound      = sdkerrors.Register(ModuleName, 308, "contract not found")
)
