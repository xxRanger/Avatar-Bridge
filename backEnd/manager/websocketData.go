package manager

import (
	"github.com/xxRanger/blockchainUtil/contract/gameToken"
	"math/big"
)

// GCUID
const (
	GCUID_LOGIN = iota
	GCUID_LOGOUT
	GCUID_REWARD
	GCUID_CONSUME
	GCUID_EXCHANGE
	GCUID_GET_BALANCE
	GCUID_GET_AVATAR
	GCUID_UPGRADE_AVATAR
)

//  status

const (
	SUCCESS = 0
	FAIL    = 1
	PENDING = 2
)

// reward type
const (
	REWARD_TYPE_ERC20TOKEN = iota
	REWARD_TYPE_AVATAR
)

// consume type
const (
	CONSUME_TYPE_ERC20TOKEN = iota
	CONSUME_TYPE_WEAPON
	CONSUME_TYPE_ARMOR
)

// exchange type
const (
	EXCHANGE_TYPE_ERC20TOKEN = iota
	EXCHANGE_TYPE_AVATAR
)

// source type
const (
	SOURCE_PRIVATE_CHAIN = iota
	SOURCE_PUBLIC_CHAIN
)

// Exchange State
const (
	EXCHANGE_STATE_PENDING = iota
	EXCHANGE_STATE_FINISH
)

type Response struct {
	Gcuid  uint64 `json:"gcuid"`
	Status int    `json:"status"`
}

func NewResponse(gcuid uint64, status int) *Response {
	return &Response{
		Gcuid:  gcuid,
		Status: status,
	}
}

type ErrorResponse struct {
	*Response
	Reason string `json:"reason"`
}

func NewErrorResponse(gcuid uint64, reason string) *ErrorResponse {
	return &ErrorResponse{
		Response: NewResponse(gcuid, FAIL),
		Reason:   reason,
	}
}

type PendingResponse struct {   // use for time consuming
	*Response
}

func NewPendingResponse(gcuid uint64) *PendingResponse {
	return &PendingResponse {
		Response: NewResponse(gcuid,PENDING),
	}
}

type SuccessResponse struct {
	*Response
}

func NewSuccessResponse(gcuid uint64) *SuccessResponse {
	return &SuccessResponse{
		Response:&Response{
			Gcuid:gcuid,
			Status:SUCCESS,
		},
	}
}

type LoginRequest struct {
	Gcuid   uint64 `json:"gcuid"`
	Address string `json:"address"`
}

type LoginResponse struct {
	*SuccessResponse
}

type LogoutRequest struct {
	Gcuid uint64 `json:"gcuid"`
}

type LogoutResponse struct {
	*SuccessResponse
}

type ConsumeRequest struct {
	Gcuid  uint64   `json:"gcuid"`
	Type int `json:"type"`
	Amount *big.Int `json:"amount"` // used for token
	Address string `json:"address"`  // used for token and avatar
	TokenId string `json:"tokenId"` // used for avatar
}

type ConsumeResponse struct {
	Type int `json:"type"`
	Amount *big.Int `json:"amount"` // used for token
	TokenId string `json:"tokenId"` //used for avatar
	*SuccessResponse
}

type RewardRequest struct {
	Gcuid  uint64   `json:"gcuid"`
	Address string `json:"address"`
	Type int `json:"type"`
	Amount *big.Int `json:"amount"`  // used for reward token
}

type RewardResponse struct {
	Type int `json:"type"`
	Amount *big.Int `json:"amount"`  // used for reward token
	TokenId string `json:"tokenId"` // used for reward avatar
	*gameToken.AvatarState // used for reward avatar
	*SuccessResponse
}

type ExchangeRequest struct {
	Gcuid  uint64 `json:"gcuid"`
	Type   int    `json:"type"`
	Amount *big.Int `json:"amount"`  // used for token
	Source int    `json:"source"` // public chain or private chain
	Transaction string `json:"transaction"`
}

type ExchangeResponse struct {
	State int `json:"state"`
	Type int `json:"type"`
	Source int `json:"source"`
	Amount *big.Int `json:"amount"` // used for exchange Token
	TokenId string `json:"tokenId"`  // used for exchange Avatar
	*gameToken.AvatarState  //used for exchange Avatar
	*SuccessResponse
}

type GetBalanceRequest struct {
	Gcuid uint64 `json:"gcuid"`
	Source int `json:"source"`
	Address string `json:"address"`
}

type GetBalanceResponse struct {
	Amount *big.Int `json:"amount"`
	Source int `json:"source"`
	*SuccessResponse
}

type GetAvatarStateRequest struct {
	Gcuid uint64 `json:"gcuid"`
	Source int `json:"source"`
	Address string `json:"address"`
}

type GetAvatarStateResponse struct {
	*gameToken.AvatarState
	TokenId string `json:"tokenId"`
	Source int `json:"source"`
	*SuccessResponse
}

type UpgradeRequest struct {
	Gcuid uint64 `json:"gcuid"`
	TokenId string `json:"tokenId"`
}

type UpgradeResponse struct {
	*SuccessResponse
}