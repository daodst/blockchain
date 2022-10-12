package rest

import (
	"encoding/json"
	"errors"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/pledge/client/cli"
	"freemasonry.cc/blockchain/x/pledge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gorilla/mux"
	"net/http"
)

//
var txHandles *TxHandles

// 
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 2
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)

	// 
	txHandles = newTxHandles(clientCtx)
	/********  ********/

	//
	txHandles.Add(types.TypeMsgPledge, PledgeHandlerFn)
	txHandles.Add(types.TypeMsgUnpledge, UnpledgeHandlerFn)
	txHandles.Add(types.TypeMsgWithdrawDelegatorReward, PledgeReceiveHandlerFn)
}

// 
func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// Query a delegation between a delegator and a validator
	r.HandleFunc(
		"/pledge/delegators/{delegatorAddr}/delegations/{validatorAddr}",
		delegationHandlerFn(clientCtx),
	).Methods("GET")

	r.HandleFunc(
		"/pledge/delegators/{delegatorAddr}/delegations",
		delegatorDelegationsHandlerFn(clientCtx),
	).Methods("GET")

	r.HandleFunc(
		"/distribution/parameters",
		paramsHandlerFn(clientCtx),
	).Methods("GET")

}

// 
func registerTxHandlers(clientCtx client.Context, r *mux.Router) {}

func ProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "delegate",
		Handler:  postProposalHandlerFn(clientCtx),
	}
}

func SendReponse(w http.ResponseWriter, clientCtx client.Context, body interface{}) {
	resBytes, err := json.Marshal(body)
	if err != nil {
		return
	}
	rest.PostProcessResponseBare(w, clientCtx, resBytes)
}

type TxHandles struct {
	ctx     client.Context
	funcMap map[string]func([]byte, *client.Context, legacytx.StdFee, string) error
}

func (this *TxHandles) Add(type1 string, func1 func([]byte, *client.Context, legacytx.StdFee, string) error) {
	this.funcMap[type1] = func1
}

func newTxHandles(ctx client.Context) *TxHandles {
	txHandles := &TxHandles{
		ctx:     ctx,
		funcMap: make(map[string]func([]byte, *client.Context, legacytx.StdFee, string) error),
	}
	return txHandles
}

//
func (this *TxHandles) HaveRegistered(msgType string) bool {
	_, ok := this.funcMap[msgType]
	return ok
}

//
func (this *TxHandles) Handle(msgType string, msgBytes []byte, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithField("msg", msgType)
	if !this.HaveRegistered(msgType) {
		log.Error("No handle registered!")
		return errors.New("msgType:" + msgType + " No handle registered!!")
	}
	log.Info("do") //
	return this.funcMap[msgType](msgBytes, &this.ctx, fee, memo)
}

func postProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req cli.PledgeDelegateProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		content := types.NewPledgeDelegateProposal(req.Title, req.Description, req.Delegate.ToParamChanges())

		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
