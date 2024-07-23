package types

import (
	"context"
	chatTypes "freemasonry.cc/blockchain/x/chat/types"
	gatewaytypes "freemasonry.cc/blockchain/x/gateway/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected interface needed to retrieve account info.
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetSequence(sdk.Context, sdk.AccAddress) (uint64, error)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	IsSendEnabledCoin(ctx sdk.Context, coin sdk.Coin) bool
	BlockedAddr(addr sdk.AccAddress) bool
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata)
	HasSupply(ctx sdk.Context, denom string) bool
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

type DistributionKeeper interface {
	GetFeePool(ctx sdk.Context) (feePool distributiontypes.FeePool)
	SetFeePool(ctx sdk.Context, feePool distributiontypes.FeePool)
	AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins)
}

type GatewayKeeper interface {
	GetGatewayInfo(ctx sdk.Context, gatewayAddress string) (*gatewaytypes.Gateway, error)
	UpdateGatewayInfo(ctx sdk.Context, gateway gatewaytypes.Gateway) error
	GetGatewayInfoByNum(ctx sdk.Context, gatewayNum string) (*gatewaytypes.Gateway, error)
}

type ChatKeeper interface {
	GetRegisterInfo(ctx sdk.Context, fromAddress string) (userinfo chatTypes.UserInfo, err error)
	Register(ctx sdk.Context, regAddr, regChatAddr, gateawyAddr string, getDid []string) error
}

type GroupKeeper interface {
	CreateGroup(goCtx context.Context, req *group.MsgCreateGroup) (*group.MsgCreateGroupResponse, error)
	CreateGroupPolicy(goCtx context.Context, req *group.MsgCreateGroupPolicy) (*group.MsgCreateGroupPolicyResponse, error)
	GetGroupInfo(goCtx context.Context, groupID uint64) (group.GroupInfo, error)
	UpdateGroupMembers(goCtx context.Context, req *group.MsgUpdateGroupMembers) (*group.MsgUpdateGroupMembersResponse, error)
	GetGroupMember(ctx sdk.Context, member *group.GroupMember) (*group.GroupMember, error)
	GroupMembers(goCtx context.Context, request *group.QueryGroupMembersRequest) (*group.QueryGroupMembersResponse, error)
	LeaveGroup(goCtx context.Context, req *group.MsgLeaveGroup) (*group.MsgLeaveGroupResponse, error)
	GroupPolicyInfo(goCtx context.Context, request *group.QueryGroupPolicyInfoRequest) (*group.QueryGroupPolicyInfoResponse, error)
	ProposalsByGroupPolicy(goCtx context.Context, request *group.QueryProposalsByGroupPolicyRequest) (*group.QueryProposalsByGroupPolicyResponse, error)
	Proposal(goCtx context.Context, request *group.QueryProposalRequest) (*group.QueryProposalResponse, error)
	VotesByProposal(goCtx context.Context, request *group.QueryVotesByProposalRequest) (*group.QueryVotesByProposalResponse, error)
	VoteByProposalVoter(goCtx context.Context, request *group.QueryVoteByProposalVoterRequest) (*group.QueryVoteByProposalVoterResponse, error)
	GroupInfo(goCtx context.Context, request *group.QueryGroupInfoRequest) (*group.QueryGroupInfoResponse, error)
}
