package manager

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"log"
)

func RLPToTransaction(rawTransactionString string)(*types.Transaction,error) {
	rawTx,err:=hex.DecodeString(rawTransactionString)
	if err!=nil {
		log.Println(err.Error())
		return nil,err
	}

	tx:=new(types.Transaction)
	err=rlp.DecodeBytes(rawTx,tx)
	if err!=nil {
		log.Println(err.Error())
	}
	return tx,err
}

func HexToPrivateKey(rawKeyString string ) (*ecdsa.PrivateKey,error) {
	return crypto.HexToECDSA(rawKeyString)
}