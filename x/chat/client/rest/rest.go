package rest

import (
	"encoding/json"
	"errors"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/chat/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
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
	txHandles.Add(types.TypeMsgRegister, RegisterHandlerFn)
	txHandles.Add(types.TypeMsgSetChatFee, SetChatFeeHandlerFn)
	txHandles.Add(types.TypeMsgSendGift, SendGiftHandlerFn)
	txHandles.Add(types.TypeMsgAddressBookSave, AddressBookSaveHandlerFn)
	txHandles.Add(types.TypeMsgMobileTransfer, MobileTransferHandlerFn)
	txHandles.Add(types.TypeMsgChangeGateway, ChangeGatewayHandlerFn)
	txHandles.Add(types.TypeMsgBurnGetMobile, BurnGetMobileHandlerFn)
	txHandles.Add(types.TypeMsgSetChatInfo, SetChatInfoHandlerFn)
}

// 
func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	//
	r.HandleFunc("/chat/accountNumberSeq/{address}", AccountNumberSeqHandlerFn(clientCtx)).Methods("GET")
}

// 
func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	//tx，
	r.HandleFunc("/chat/tx/broadcast", BroadcastTxHandlerFn(clientCtx)).Methods("POST")

	r.HandleFunc(
		"/chat/parameters",
		paramsHandlerFn(clientCtx),
	).Methods("GET")

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
