package rest

import (
	"errors"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/comm/types"
	contractCli "freemasonry.cc/blockchain/x/contract/client/cli"
	contractTypes "freemasonry.cc/blockchain/x/contract/types"
	"freemasonry.cc/blockchain/x/pledge/client/cli"
	pledgeTypes "freemasonry.cc/blockchain/x/pledge/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func MsgHandlerFun(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp struct {
			Data []byte `json:"data"`
		}
		var paramsByte []byte
		if r.Body != nil {
			paramsByte, _ = ioutil.ReadAll(r.Body)
		}
		msgToByte := MsgToByte{}
		err := util.Json.Unmarshal(paramsByte, &msgToByte)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		msgByte := []byte(msgToByte.Msg)
		switch msgToByte.MsgType {
		case types.TypeMsgCreateSmartValidator:
			msg := types.MsgCreateSmartValidator{}
			err = util.Json.Unmarshal(msgByte, &msg)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			err = msg.ValidateBasic()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msgByte, err = msg.Marshal()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			resp.Data = msgByte
			SendReponse(w, clientCtx, resp)
			return
		case types.TypeMsgGatewayRegister:
			msg := types.MsgGatewayRegister{}
			err = util.Json.Unmarshal(msgByte, &msg)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			err = msg.ValidateBasic()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msgByte, err = msg.Marshal()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			resp.Data = msgByte
			SendReponse(w, clientCtx, resp)
			return
		case types.TypeMsgGatewayIndexNum:
			msg := types.MsgGatewayIndexNum{}
			err = util.Json.Unmarshal(msgByte, &msg)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			err = msg.ValidateBasic()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msgByte, err = msg.Marshal()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			resp.Data = msgByte
			SendReponse(w, clientCtx, resp)
			return
		case types.TypeMsgGatewayUndelegation:
			msg := types.MsgGatewayUndelegate{}
			err = util.Json.Unmarshal(msgByte, &msg)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			err = msg.ValidateBasic()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msgByte, err = msg.Marshal()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			resp.Data = msgByte
			SendReponse(w, clientCtx, resp)
			return
		case "cosmos-sdk/MsgEditValidator":
			msg := stakingTypes.MsgEditValidator{}
			err = util.Json.Unmarshal(msgByte, &msg)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			err = msg.ValidateBasic()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msgByte, err = msg.Marshal()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			resp.Data = msgByte
			SendReponse(w, clientCtx, resp)
		case "cosmos-sdk_MsgUnjail":
			msg := slashingTypes.MsgUnjail{}
			err = util.Json.Unmarshal(msgByte, &msg)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			err = msg.ValidateBasic()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msgByte, err = msg.Marshal()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			resp.Data = msgByte
			SendReponse(w, clientCtx, resp)
		case types.TypeMsgGatewayBeginRedelegate:
			msg := types.MsgGatewayBeginRedelegate{}
			err = util.Json.Unmarshal(msgByte, &msg)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			err = msg.ValidateBasic()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msgByte, err = msg.Marshal()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			resp.Data = msgByte
			SendReponse(w, clientCtx, resp)
		case "proposal_delegate":
			//	proposal := cli.DelegateProposalJSON{}
			proposal := struct {
				Proposer       string    `json:"proposer"`
				InitialDeposit sdk.Coins `json:"initial_deposit"`
				Content        struct {
					Type  string `json:"type"`
					Value cli.DelegateProposalJSON
				}
			}{}
			err = util.Json.Unmarshal(msgByte, &proposal)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			content := pledgeTypes.NewPledgeDelegateProposal(proposal.Content.Value.Title, proposal.Content.Value.Description, proposal.Content.Value.Delegate.ToParamChanges())
			deposit := proposal.InitialDeposit
			from, err := sdk.AccAddressFromBech32(msgToByte.Address)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			err = msg.ValidateBasic()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msgByte, err = msg.Marshal()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			resp.Data = msgByte
			SendReponse(w, clientCtx, resp)
		case "proposal_params":
			//proposal := paramsUtils.ParamChangeProposalJSON{}
			proposal := struct {
				Proposer       string    `json:"proposer"`
				InitialDeposit sdk.Coins `json:"initial_deposit"`
				Content        struct {
					Type  string `json:"type"`
					Value struct {
						Title       string           `json:"title" yaml:"title"`
						Description string           `json:"description" yaml:"description"`
						Changes     ParamChangesJSON `json:"changes" yaml:"changes"`
						Deposit     string           `json:"deposit" yaml:"deposit"`
					}
				}
			}{}
			err = util.Json.Unmarshal(msgByte, &proposal)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			content := paramproposal.NewParameterChangeProposal(proposal.Content.Value.Title, proposal.Content.Value.Description, proposal.Content.Value.Changes.ToParamChanges())
			deposit := proposal.InitialDeposit
			from, err := sdk.AccAddressFromBech32(msgToByte.Address)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			err = msg.ValidateBasic()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msgByte, err = msg.Marshal()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			resp.Data = msgByte
			SendReponse(w, clientCtx, resp)
		case "proposal_community":
			//proposal := paramsUtils.ParamChangeProposalJSON{}
			proposal := struct {
				Proposer       string    `json:"proposer"`
				InitialDeposit sdk.Coins `json:"initial_deposit"`
				Content        struct {
					Type  string `json:"type"`
					Value CommunityPoolSpendProposalWithDeposit
				}
			}{}
			err = util.Json.Unmarshal(msgByte, &proposal)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			recpAddr, err := sdk.AccAddressFromBech32(proposal.Content.Value.Recipient)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			content := distributionTypes.NewCommunityPoolSpendProposal(proposal.Content.Value.Title, proposal.Content.Value.Description, recpAddr, proposal.Content.Value.Amount)
			deposit := proposal.InitialDeposit
			from, err := sdk.AccAddressFromBech32(msgToByte.Address)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			err = msg.ValidateBasic()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msgByte, err = msg.Marshal()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			resp.Data = msgByte
			SendReponse(w, clientCtx, resp)
		case "proposal_upgrade":
			proposal := struct {
				Proposer       string    `json:"proposer"`
				InitialDeposit sdk.Coins `json:"initial_deposit"`
				Content        struct {
					Type  string `json:"type"`
					Value SoftwareUpgradeProposal
				}
			}{}
			err = util.Json.Unmarshal(msgByte, &proposal)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			height, err := strconv.ParseInt(proposal.Content.Value.Plan.Height, 10, 64)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			plan := upgradeTypes.Plan{
				Name:   proposal.Content.Value.Plan.Name,
				Height: height,
				Info:   proposal.Content.Value.Plan.Info,
			}

			err = UpgradeJsonValidateBasic(plan)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			content := upgradeTypes.NewSoftwareUpgradeProposal(proposal.Content.Value.Title, proposal.Content.Value.Description, plan)
			deposit := proposal.InitialDeposit
			from, err := sdk.AccAddressFromBech32(msgToByte.Address)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			err = msg.ValidateBasic()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msgByte, err = msg.Marshal()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			resp.Data = msgByte
			SendReponse(w, clientCtx, resp)
		case "proposal_contract":
			proposal := struct {
				Proposer       string    `json:"proposer"`
				InitialDeposit sdk.Coins `json:"initial_deposit"`
				Content        struct {
					Type  string `json:"type"`
					Value contractCli.ContractProposalJSON
				}
			}{}
			err = util.Json.Unmarshal(msgByte, &proposal)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			content := contractTypes.NewContractProposal(proposal.Content.Value.Title, proposal.Content.Value.Description, proposal.Content.Value.Contract.ToParamChanges())
			deposit := proposal.InitialDeposit
			from, err := sdk.AccAddressFromBech32(msgToByte.Address)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if rest.CheckInternalServerError(w, err) {
				return
			}
			err = msg.ValidateBasic()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			msgByte, err = msg.Marshal()
			if rest.CheckInternalServerError(w, err) {
				return
			}
			resp.Data = msgByte
			SendReponse(w, clientCtx, resp)
		}
	}
}

type (
	MsgToByte struct {
		Address string `json:"address"`
		MsgType string `json:"msg_type"`
		Msg     string `json:"msg"`
	}

	ParamChangesJSON []ParamChangeJSON

	ParamChangeJSON struct {
		Subspace string `json:"subspace" yaml:"subspace"`
		Key      string `json:"key" yaml:"key"`
		Value    string `json:"value" yaml:"value"`
	}

	CommunityPoolSpendProposalWithDeposit struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Recipient   string    `json:"recipient"`
		Amount      sdk.Coins `json:"amount"`
	}

	SoftwareUpgradeProposal struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Plan        Plan   `json:"plan"`
	}
	Plan struct {
		Name   string    `json:"name"`
		Time   time.Time `json:"time"`
		Height string    `json:"height"`
		Info   string    `json:"info"`
	}
)

func (pcj ParamChangesJSON) ToParamChanges() []paramproposal.ParamChange {
	res := make([]paramproposal.ParamChange, len(pcj))
	for i, pc := range pcj {
		res[i] = pc.ToParamChange()
	}
	return res
}

func (pcj ParamChangeJSON) ToParamChange() paramproposal.ParamChange {
	return paramproposal.NewParamChange(pcj.Subspace, pcj.Key, pcj.Value)
}

func UpgradeJsonValidateBasic(plan upgradeTypes.Plan) error {
	info, err := plan.UpgradeInfo()
	if err != nil {
		return err
	}
	if info.Gateway == nil && info.App == nil && info.Blockchain == nil {
		return errors.New("The json content is illegal")
	}
	return nil
}
