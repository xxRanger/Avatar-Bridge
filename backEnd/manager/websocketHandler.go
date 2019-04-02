package manager

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/xxRanger/blockchainUtil/contract/bridgeToken"
	"github.com/xxRanger/blockchainUtil/contract/gameToken"
	"github.com/xxRanger/bridge/backEnd/client"
	"log"
	"math/big"
)

func (m *Manager) LoginHandler(c *client.Client, gcuid uint64, data []byte) {
	var req LoginRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}
	c.Login(common.HexToAddress(req.Address))
	err = m.clientLoginHandler(c)
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}
	m.wrapperAndSend(c, gcuid, &LoginResponse{
		SuccessResponse: NewSuccessResponse(gcuid),
	})
}

func (m *Manager) LogoutHandler(c *client.Client, gcuid uint64, data []byte) {
	c.Logout()
	m.clientLogoutHandler(c)

	m.wrapperAndSend(c, gcuid, &LogoutResponse{
		SuccessResponse: NewSuccessResponse(gcuid),
	})
}

func (m *Manager) GetBalanceHandler(c *client.Client, gcuid uint64, data []byte) {
	var req GetBalanceRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	log.Println(req.Address, "get balance");
	account := common.HexToAddress(req.Address)
	balance := big.NewInt(0)
	switch req.Source {
	case SOURCE_PUBLIC_CHAIN:
		bridgeTokenContract := m.chainHandler.PbcHandler.BridgeTokenContract
		balance, err = bridgeTokenContract.BalanceOf(account)
	case SOURCE_PRIVATE_CHAIN:
		bridgeTokenContract := m.chainHandler.PvcHandler.BridgeTokenContract
		balance, err = bridgeTokenContract.BalanceOf(account)
	default:
		err = errors.New("unknown source")
	}
	if err != nil {
		m.errorHandler(c, gcuid, err)
		return
	}
	m.wrapperAndSend(c, gcuid, &GetBalanceResponse{
		Source:          req.Source,
		Amount:          balance,
		SuccessResponse: NewSuccessResponse(gcuid),
	})
}

func (m *Manager) GetAvatarStateHandler(c *client.Client, gcuid uint64, data []byte) {
	var req GetAvatarStateRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	log.Println(req.Address, "get avatar state");
	account := common.HexToAddress(req.Address)

	//get token Id
	tokenId := new(big.Int)
	bridgeTokenContract := new(bridgeToken.BridgeToken)
	switch req.Source {
	case SOURCE_PUBLIC_CHAIN:
		bridgeTokenContract = m.chainHandler.PbcHandler.BridgeTokenContract
		tokenId, err = bridgeTokenContract.OwnedAvatar(account)
	case SOURCE_PRIVATE_CHAIN:
		bridgeTokenContract = m.chainHandler.PvcHandler.BridgeTokenContract
		tokenId, err = bridgeTokenContract.OwnedAvatar(account)
	default:
		err = errors.New("unknown source")
	}
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	// get avatar state
	avatarState := &gameToken.AvatarState{}
	switch req.Source {
	case SOURCE_PUBLIC_CHAIN:
		avatarState, err = bridgeTokenContract.AvatarState(tokenId)
	case SOURCE_PRIVATE_CHAIN:
		avatarState, err = bridgeTokenContract.AvatarState(tokenId)
	}
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}
	log.Println("avatar state:", gcuid, "tokenId:", tokenId)
	m.wrapperAndSend(c, gcuid, &GetAvatarStateResponse{
		TokenId:         tokenId.String(),
		AvatarState:     avatarState,
		Source: req.Source,
		SuccessResponse: NewSuccessResponse(gcuid),
	})
}

func (m *Manager) ConsumeHandler(c *client.Client, gcuid uint64, data []byte) {
	var req ConsumeRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	//log.Println("to:",tx.To().Hex())
	//log.Println("gas price:",tx.GasPrice())
	//log.Println("gas limit:",tx.Gas())
	//r,s,v:=tx.RawSignatureValues()
	//log.Println("r,s,v",r,s,v)
	//log.Println("tx data",hex.EncodeToString(tx.Data()))
	//log.Println("nonce:",tx.Nonce())
	//log.Println("txHash",tx.Hash().Hex())

	if req.Type != CONSUME_TYPE_ERC20TOKEN && req.Type != CONSUME_TYPE_WEAPON && req.Type != CONSUME_TYPE_ARMOR {
		err := errors.New("unkown consume type")
		log.Println(CONSUME_TYPE_ERC20TOKEN)
		m.errorHandler(c, gcuid, err)
		return
	}

	manageAccount := m.chainHandler.PvcHandler.ManagerAccount
	bridgeTokenContract := m.chainHandler.PvcHandler.BridgeTokenContract
	txErr := manageAccount.SendFunction(
		bridgeTokenContract,
		nil,
		bridgeToken.FuncConsume,
		common.HexToAddress(req.Address),
		req.Amount,
	)
	err = <-txErr
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	consumeResponse:= &ConsumeResponse{
		Amount:          req.Amount,
		Type:            req.Type,
		SuccessResponse: NewSuccessResponse(gcuid),
	}

	if req.Type != CONSUME_TYPE_ERC20TOKEN {
		log.Println("buy weapon");
		tokenId,_ := new(big.Int).SetString(req.TokenId,10)
		var txErr chan error
		switch req.Type {
		case CONSUME_TYPE_WEAPON:
			txErr = manageAccount.SendFunction(
				bridgeTokenContract,
				nil,
				bridgeToken.FuncEquipWeapon,
				tokenId,
				common.HexToAddress(req.Address),
			)
		case CONSUME_TYPE_ARMOR:
			txErr = manageAccount.SendFunction(
				bridgeTokenContract,
				nil,
				bridgeToken.FuncEquipArmor,
				tokenId,
				common.HexToAddress(req.Address),
			)
		}
		err = <-txErr
		if err != nil {
			log.Println(err.Error())
			m.errorHandler(c, gcuid, err)
			return
		}
		consumeResponse.TokenId = req.TokenId
	}
	m.wrapperAndSend(c,gcuid,consumeResponse)
}

func (m *Manager) RewardHandler(c *client.Client, gcuid uint64, data []byte) {
	var req RewardRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	user := common.HexToAddress(req.Address)

	rewardResponse := &RewardResponse{
		Type:            req.Type,
		SuccessResponse: NewSuccessResponse(gcuid),
	}
	bridgeTokenContract := m.chainHandler.PvcHandler.BridgeTokenContract
	switch req.Type {
	case REWARD_TYPE_ERC20TOKEN:
		txErr := m.chainHandler.PvcHandler.ManagerAccount.SendFunction(bridgeTokenContract,
			nil,
			bridgeToken.FuncReward,
			user,
			req.Amount,
		)
		err = <-txErr
		if err != nil {
			log.Println(err.Error())
			m.errorHandler(c, gcuid, err)
			return
		}
		rewardResponse.Amount = req.Amount
	case REWARD_TYPE_AVATAR:
		tokenId := big.NewInt(0)
		for tokenId.Cmp(big.NewInt(0)) == 0 {
			b := make([]byte, 256)
			rand.Read(b)
			tokenId = new(big.Int).SetBytes(b)
		}
		txErr := m.chainHandler.PvcHandler.ManagerAccount.SendFunction(bridgeTokenContract,
			nil,
			bridgeToken.FuncMint,
			user,
			tokenId,
		)
		err = <-txErr
		if err != nil {
			log.Println(err.Error())
			m.errorHandler(c, gcuid, err)
			return
		}
		rewardResponse.TokenId = tokenId.String()
		// get avatar info
		avatarState, err := m.chainHandler.PvcHandler.BridgeTokenContract.AvatarState(tokenId)
		if err != nil {
			log.Println(err.Error())
			m.errorHandler(c, gcuid, err)
			return
		}
		rewardResponse.AvatarState = avatarState
	}

	m.wrapperAndSend(c, gcuid, rewardResponse)
}

func (m *Manager) ExchangeHandler(c *client.Client, gcuid uint64, data []byte) {
	var req ExchangeRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	txRLP := req.Transaction
	tx, err := RLPToTransaction(txRLP)
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	sourceUnknown := false
	var txErr chan error
	if req.Source == SOURCE_PRIVATE_CHAIN {
		if req.Type != EXCHANGE_TYPE_ERC20TOKEN && req.Type != EXCHANGE_TYPE_AVATAR {
			sourceUnknown = true
		} else {
			txErr = m.chainHandler.PvcHandler.Client.Send(tx)
		}
	} else if req.Source == SOURCE_PUBLIC_CHAIN {
		if req.Type != EXCHANGE_TYPE_ERC20TOKEN && req.Type != EXCHANGE_TYPE_AVATAR {
			sourceUnknown = true
		} else {
			txErr = m.chainHandler.PbcHandler.Client.Send(tx)
		}
	} else {
		err := errors.New("unknown exchange source chain")
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	if sourceUnknown {
		err := errors.New("unknown exchange type")
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	err = <-txErr
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	m.wrapperAndSend(c, gcuid, &ExchangeResponse{
		Amount:          req.Amount,
		State:           EXCHANGE_STATE_PENDING,
		Source:          req.Source,
		Type:            req.Type,
		SuccessResponse: NewSuccessResponse(gcuid),
	})
}

func (m *Manager) UpgradeHanlder(c *client.Client, gcuid uint64, data []byte) {
	var req UpgradeRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	var txErr chan error
	manageAccount := m.chainHandler.PvcHandler.ManagerAccount
	bridgeTokenContract := m.chainHandler.PvcHandler.BridgeTokenContract

	tokenId,_:= new(big.Int).SetString(req.TokenId,10)
	txErr = manageAccount.SendFunction(
		bridgeTokenContract,
		nil,
		bridgeToken.FuncUpgrade,
		tokenId,
	)

	err = <-txErr
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	m.wrapperAndSend(c, gcuid, &UpgradeResponse{
		SuccessResponse: NewSuccessResponse(gcuid),
	})
}

func (m *Manager) errorHandler(c *client.Client, gcuid uint64, err error) {
	res := NewErrorResponse(gcuid, err.Error())

	resWrapper, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		return
	}

	c.Send(resWrapper)
}
