package core

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/shopspring/decimal"
	abciTypes "github.com/tendermint/tendermint/abci/types"
)

const (
	MSG_TYPE_CREATE_COPYRIGHT            = "copyright/MsgCreateCopyright"
	MSG_TYPE_REGISTER_COPYRIGHT_PARTY    = "copyright/MsgRegisterCopyrightParty"
	MSG_TYPE_SPACE_MINER                 = "copyright/MsgSpaceMiner"
	MSG_TYPE_DEFLATION_VOTE              = "copyright/MsgDeflationVote"
	MSG_TYPE_DISTRIBUTE_COMMUNITY_REWARD = "copyright/MsgDistributeCommunityReward"
	MSG_TYPE_NFT_TRANSFER                = "copyright/MsgNftTransfer"
	MSG_TYPE_INVITE_CODE                 = "copyright/MsgInviteCode"
	MSG_TYPE_MORTGAGE                    = "copyright/MsgMortgage"
	MSG_TYPE_EDITOR_COPYRIGHT            = "copyright/MsgEditorCopyright"
	MSG_TYPE_DELETE_COPYRIGHT            = "copyright/MsgDeleteCopyright"
	MSG_TYPE_BONUS_COPYRIGHT             = "copyright/MsgCopyrightBonus"
	MSG_TYPE_BONUS_COPYRIGHTV2           = "copyright/MsgCopyrightBonusV2"
	MSG_TYPE_COPYRIGHT_COMPLAIN          = "copyright/MsgCopyrightComplain"
	MSG_TYPE_COMPLAIN_RESPONSE           = "copyright/MsgComplainResponse"
	MSG_TYPE_COMPLAIN_VOTE               = "copyright/MsgComplainVote"
	MSG_TYPE_AUTHORIZE_NODE              = "copyright/MsgAuthorizeNode"
	MSG_TYPE_TRANSFER                    = "copyright/MsgTransfer"
	MSG_TYPE_INVATE_REWARD               = "copyright/MsgInviteReward"
	MSG_TYPE_SPACE_MINER_REWARD          = "copyright/MsgSpaceMinerReward"
	MSG_TYPE_BONUS_COPYRIGHT_REAR        = "copyright/MsgCopyrightBonusRear"
	MSG_TYPE_BONUS_COPYRIGHT_REARV2      = "copyright/MsgCopyrightBonusRearV2"
	MSG_TYPE_COPYRIGHT_VOTE              = "copyright/MsgCopyrightVote"
	MSG_TYPE_CROSSCHAIN_OUT              = "copyright/MsgCrossChainOut"
	MSG_TYPE_CROSSCHAIN_IN               = "copyright/MsgCrossChainIn"
)

//
type RealCoin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

func (this RealCoin) Add(rcoin2 RealCoin) RealCoin {
	dec1 := sdk.MustNewDecFromStr(rcoin2.Amount)
	dec2 := sdk.MustNewDecFromStr(this.Amount)
	if this.Denom != rcoin2.Denom {
		panic(errors.New("added realcoin denom It has to be the same"))
	}
	return RealCoin{Amount: dec1.Add(dec2).String(), Denom: this.Denom}
}
func (this RealCoin) AmountDec() decimal.Decimal {
	if this.Amount == "" {
		return decimal.Decimal{}
	} else {
		ret, err := decimal.NewFromString(this.Amount)
		if err != nil {

		}
		return ret
	}
}

func (this RealCoin) String() string {
	return this.Amount + this.Denom
}

//
func (this RealCoin) FormatAmount() string {
	return RemoveStringLastZero(this.Amount)
}

//string
func (this RealCoin) FormatString() string {
	return this.FormatAmount() + this.Denom
}

//
type RealCoins []RealCoin

func (this RealCoins) String() string {
	tmp := ""
	for _, v := range this {
		tmp += v.String() + " "
	}
	return tmp
}

func (this RealCoins) Get(index int) RealCoin {
	return this[index]
}

type MortgageData struct {
	TradeBase
	TxBase
	Creator           sdk.AccAddress `json:"creator"`             //
	MortgageAccount   sdk.AccAddress `json:"mortgage_account"`    //
	DataHash          string         `json:"data_hash"`           //hash
	CopyrightPrice    RealCoin       `json:"copyright_price"`     //
	CreateTime        int64          `json:"create_time"`         //
	MortgageAmount    RealCoin       `json:"mortgage_amount"`     //
	OfferAccountShare string         `json:"offer_account_share"` //
	DataHashAccount   sdk.AccAddress `json:"data_hash_account"`   //
	BonusType         string         `json:"bonus_type"`          //
}

//
type BlockMessageI interface {
	Type() string //
	Unmarshal([]byte) (interface{}, error)
	Marshal() []byte
}

type SpaceMinerData struct {
	TradeBase
	TxBase
	Creator         sdk.AccAddress `json:"creator"`          //
	DeflationAmount RealCoin       `json:"deflation_amount"` //
	AwardAccount    sdk.AccAddress `json:"award_account"`    //
}

type DeflationVoteData struct {
	TradeBase
	TxBase
	Creator sdk.AccAddress `json:"creator"` //
	Option  string         `json:"option"`  //
}

type NftTransferData struct {
	TradeBase
	TxBase
	From    sdk.AccAddress `json:"from"`     //
	To      sdk.AccAddress `json:"to"`       //
	TokenId string         `json:"token_id"` //id
}

//tx
type TxBase struct {
	BlockTime int64           `json:"block_time"` //
	Height    int64           `json:"height"`
	TxHash    string          `json:"tx_hash"`
	Fee       legacytx.StdFee `json:"fee"`  //
	Memo      string          `json:"memo"` //
}

//
type TradeBase struct {
	TradeType       TranserType `json:"trade_type"`       //
	ContractAddress string      `json:"contract_address"` //
}

//tx
type TxPropValue struct {
	BlockTime int64 //
	TxHash    string
	Height    int64
	Seq       int64
	Fee       legacytx.StdFee
	Memo      string
	Events    []abciTypes.Event
}

//txbase
func (this *TxBase) UpdateTxBase(prop *TxPropValue) {
	this.BlockTime = prop.BlockTime
	this.Height = prop.Height
	this.TxHash = prop.TxHash
	this.Fee = prop.Fee
	this.Memo = prop.Memo
}

//trade base
func (this *TradeBase) UpdateTradeBase(tradeType TranserType, contractAddr string) {
	this.TradeType = tradeType
	this.ContractAddress = contractAddr
}

//event  skipFee 
//func (this *TradeBase) AnalysisEvents(resp *MessageResp, prop *TxPropValue, transerType TranserType, analysisMethod AnalysisMethod) error {
//	for _, event := range prop.Events {
//		if event.Type == "transfer" {
//			AnalysisTransfer(resp, prop, event, transerType, analysisMethod)
//		}
//	}
//	return nil
//}

//
type CopyrightPartyData struct {
	TradeBase
	TxBase
	Id      string         `json:"id"`      //id
	Intro   string         `json:"intro"`   //
	Author  string         `json:"author"`  //
	Creator sdk.AccAddress `json:"creator"` //
}

//
type CopyrightData struct {
	TxBase
	TradeBase
	DataHash       string              `json:"datahash"`        //hash
	Price          RealCoin            `json:"price"`           //
	Creator        sdk.AccAddress      `json:"creator"`         //
	ResourceType   string              `json:"resourcetype"`    //
	PreHash        string              `json:"prehash"`         //ipfs hash
	VideoHash      string              `json:"video_hash"`      //ipfs hash
	Name           string              `json:"name"`            //
	Files          Files               `json:"files"`           //
	Size           int64               `json:"size"`            //
	CreateTime     int                 `json:"create_time"`     //
	Password       string              `json:"password"`        //64
	ChargeRate     string              `json:"charge_rate"`     //
	Ip             string              `json:"ip"`              //ip
	OriginDataHash string              `json:"origin_datahash"` //hash
	Ext            string              `json:"ext"`             //
	ClassifyUid    int64               `json:"classify_uid"`    //
	LinkMap        map[string]Link     `json:"link_map"`        //
	ApproveStatus  int                 `json:"approve_status"`  //
	PicLinkMap     map[string]Link     `json:"pic_link_map"`    //hash
	SecretMap      map[string][][]byte `json:"secret_map"`      //
}

//
type EditorCopyrightData struct {
	TxBase
	TradeBase
	DataHash   string         `json:"data_hash"`   //hash
	Price      RealCoin       `json:"price"`       //
	Creator    sdk.AccAddress `json:"creator"`     //
	Name       string         `json:"name"`        //
	ChargeRate string         `json:"charge_rate"` //
	Ip         string         `json:"ip"`          //ip
}

//
type DeleteCopyrightData struct {
	TxBase
	TradeBase
	DataHash string         `json:"data_hash"` //hash
	Creator  sdk.AccAddress `json:"creator"`   //
}

//
type CopyrightBonusData struct {
	TxBase
	TradeBase
	DataHash          string         `json:"data_hash"` //hash
	Downer            sdk.AccAddress `json:"downer"`    //
	HashAccount       sdk.AccAddress `json:"hash_account"`
	OfferAccountShare string         `json:"offer_account_share"`
	BonusType         string         `json:"bonus_type"`    //
	BonusAddress      string         `json:"bonus_address"` //
}

//
type CopyrightBonusRearData struct {
	TxBase
	TradeBase
	DataHash          string         `json:"data_hash"` //hash
	Downer            sdk.AccAddress `json:"downer"`    //
	OfferAccountShare string         `json:"offer_account_share"`
	BonusAddress      string         `json:"bonus_address"`
}

//
type CopyrightComplainData struct {
	TxBase
	TradeBase
	DataHash        string         `json:"datahash"`         //hash
	Author          string         `json:"author"`           //
	Productor       string         `json:"productor"`        //
	LegalNumber     string         `json:"legal_number"`     //
	LegalTime       string         `json:"legal_time"`       //
	ComplainInfor   string         `json:"complain_infor"`   //
	ComplainAccount sdk.AccAddress `json:"complain_account"` //
	AccuseAccount   sdk.AccAddress `json:"accuse_account"`   //
	ComplainId      string         `json:"complain_id"`      //id
	ComplainTime    int64          `json:"complain_time"`    //
}

//
type ComplainResponseData struct {
	TxBase
	TradeBase
	DataHash      string         `json:"datahash"`       //hash
	RemoteIp      string         `json:"remote_ip"`      //ip
	AccuseInfor   string         `json:"accuse_infor"`   //
	AccuseAccount sdk.AccAddress `json:"accuse_account"` //
	ComplainId    string         `json:"complain_id"`    //id
	Status        string         `json:"status"`         //
	ResponseTime  int64          `json:"complain_time"`  //
}

//
type ComplainVoteData struct {
	TxBase
	TradeBase
	VoteAccount sdk.AccAddress `json:"vote_account"` //
	ComplainId  string         `json:"complain_id"`  //id
	VoteStatus  string         `json:"vote_status"`  //
	VoteShare   sdk.Dec        `json:"vote_share"`   //
}

//
type AuthorizeAccountData struct {
	TxBase
	Account  sdk.AccAddress `json:"account"`   //
	ConsAddr string         `json:"cons_addr"` //id
	Sign     string         `json:"sign"`      //
	Message  string         `json:"message"`   //
}

//
type InviteCodeData struct {
	TxBase
	Address    string `json:"address"`
	InviteCode string `json:"invite_code"`
	InviteTime int64  `json:"invite_time"`
}

//
type Files struct {
	IsDir   bool    `json:"is_dir"`
	Size    int64   `json:"size"`
	Name    string  `json:"name"`
	Content []Files `json:"content"`
}
type FilesSizeStr struct {
	IsDir   bool           `json:"is_dir"`
	Size    int64          `json:"size"`
	Name    string         `json:"name"`
	Content []FilesSizeStr `json:"content"`
}

//
type Link struct {
	Name string `json:"name"` //
	Size uint64 `json:"size"` //
	Cid  string `json:"cid"`  //hash
}

func (fs *Files) Files2FilesSizeStr() (fs_ *FilesSizeStr) {
	fs_ = &FilesSizeStr{
		IsDir:   fs.IsDir,
		Size:    fs.Size,
		Name:    fs.Name,
		Content: []FilesSizeStr{},
	}

	if len(fs.Content) != 0 {
		for _, f := range fs.Content {
			t := f.Files2FilesSizeStr()
			fs_.Content = append(fs_.Content, *t)
		}
	}
	return
}

//
type TransferData struct {
	TxBase
	TradeBase
	FromAddress    string    `json:"from_address"` //
	ToAddress      string    `json:"to_address"`   //
	Coins          RealCoins `json:"coins"`
	FromFsvBalance string    `json:"from_fsv_balance"`
	ToFsvBalance   string    `json:"to_fsv_balance"`
	FromTipBalance string    `json:"from_tip_balance"`
	ToTipBalance   string    `json:"to_tip_balance"`
}

//
type InviteRewardData struct {
	TxBase
	TradeBase
	Address string `json:"from_address"` //
}

//
type SpaceMinerRewardData struct {
	TxBase
	TradeBase
	Address string `json:"from_address"` //
}

//
type CrossChainInData struct {
	TxBase
	TradeBase
	SendAddress string `json:"send_address"` //
	Coins       string `json:"coins"`        // eg: 15MIP
	ChainType   string `json:"chain_type"`   // eg: ETH
	Remark      string `json:"remark"`       // 
}

//
type CrossChainOutData struct {
	TxBase
	TradeBase
	SendAddress string `json:"send_address"` //
	ToAddress   string `json:"to_address"`   //
	Coins       string `json:"coins"`        // eg: 15MIP
	ChainType   string `json:"chain_type"`   // eg: ETH
	Remark      string `json:"remark"`       // 
}

//
type CopyrightVoteData struct {
	TxBase
	TradeBase
	Address  string `json:"address"`   //
	DataHash string `json:"data_hash"` //hash
	Power    string `json:"Power"`     //
}

//
type WithdrawDelegatorRewardData struct {
	TxBase
	TradeBase
	Amount           RealCoin `json:"amount"`            //coin
	DelegatorAddress string   `json:"delegator_address"` //
	ValidatorAddress string   `json:"validator_address"` //
}

//
type UndelegationData struct {
	TxBase
	TradeBase
	DelegatorAddress string   `json:"delegator_address"`
	ValidatorAddress string   `json:"validator_address"`
	Amount           RealCoin `json:"amount"` //
	Shares           string   `json:"shares"` //
	CompletionTime   string   `json:"completion_time"`
}

//
type DelegationData struct {
	TxBase
	TradeBase
	DelegatorAddress string   `json:"delegator_address"`
	ValidatorAddress string   `json:"validator_address"`
	Coin             RealCoin `json:"coin"`   //
	Shares           string   `json:"shares"` //
}

//
type DelegationDetail struct {
	DelegatorAddress          sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress          sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
	DelegationAmount          string         `json:"delegator_amount" yaml:"delegator_amount"`                       // 
	UnbindingDelegationAmount string         `json:"unbinding_delegation_amount" yaml:"unbinding_delegation_amount"` // 
	DelegationShareNumber     string         `json:"delegator_share_number" yaml:"delegator_share_number"`           //  
	ValidatorShareNumber      string         `json:"validator_share_number" yaml:"validator_share_number"`           //  
	ValidatorInfor            stakingtypes.Validator
	ValidatorSingInfo         slashingTypes.ValidatorSigningInfo
}

func NewDelegationDetail(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) DelegationDetail {
	return DelegationDetail{
		DelegatorAddress:          delegatorAddr,
		ValidatorAddress:          validatorAddr,
		DelegationAmount:          "0",
		UnbindingDelegationAmount: "0",
		DelegationShareNumber:     "0",
		ValidatorShareNumber:      "0",
	}
}

//
type ValidatorInfo struct {
	OperatorAddress     string `json:"operator_address"`
	ConsAddress         string `json:"cons_address"`
	Jailed              bool   `json:"jailed"`           //
	Status              int    `json:"status"`           // 3
	Tokens              string `json:"tokens"`           //
	DelegatorShares     string `json:"delegator_shares"` //
	Moniker             string `json:"moniker"`          //
	Identity            string `json:"identity"`         //
	Website             string `json:"website"`          //
	SecurityContact     string `json:"security_contact"` //
	Details             string `json:"details"`
	UnbondingHeight     int64  `json:"unbonding_height"`
	UnbondingTime       int64  `json:"unbonding_time"`
	Rate                string `json:"rate"`                  //
	MaxRate             string `json:"max_rate"`              //
	MaxChangeRate       string `json:"max_change_rate"`       //
	MinSelfDelegation   string `json:"min_self_delegation"`   //
	MissedBlocksCounter int64  `json:"missed_blocks_counter"` //
	IndexOffset         int64  `json:"index_offset"`          //
	PunishCount         int64  `json:"punish_count"`          //
}
