package manager

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"log"
	"math/big"
	"net/http"
)

func (m *Manager) GetNonce(w http.ResponseWriter,r *http.Request) {
	vars:=mux.Vars(r)
	chain:=vars["chain"]
	user:=vars["user"]
	var nonce uint64
	var err error
	if chain == CHAIN_TYPE_PRIVATE {
		nonce,err=m.chainHandler.PvcHandler.Client.GetNonce(common.HexToAddress(user))
	} else {
		nonce,err=m.chainHandler.PbcHandler.Client.GetNonce(common.HexToAddress(user))
	}
	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	nonceWrapper,err:=json.Marshal(nonce)
	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Write(nonceWrapper)
}

func (m *Manager) GetChainId(w http.ResponseWriter,r *http.Request) {
	vars:=mux.Vars(r)
	chain:=vars["chain"]
	var chainId *big.Int
	var err error
	if chain == CHAIN_TYPE_PRIVATE {
		chainId,err = m.chainHandler.PvcHandler.Client.GetChainId()
	} else {
		chainId,err = m.chainHandler.PbcHandler.Client.GetChainId()
	}
	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	chainIdWrapper,err:=json.Marshal(chainId)
	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Write(chainIdWrapper)
}
