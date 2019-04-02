package manager

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/xxRanger/blockchainUtil/contract/bridgeToken"
	"github.com/xxRanger/blockchainUtil/contract/gameToken"
	"log"
	"math/big"
)

func (m *Manager) registerPrivateChainHandler() {
	m.chainHandler.PvcHandler.BridgeTokenContract.RegisterHandler(bridgeToken.EventExchangeToken, m.eventPrivateChainExchangeTokenHandler)
	m.chainHandler.PvcHandler.BridgeTokenContract.RegisterHandler(bridgeToken.EventExchangeNFT, m.eventPrivateChainExchangeNFTHandler)
}

func (m *Manager) registerPublicChainHandler() {
	m.chainHandler.PbcHandler.BridgeTokenContract.RegisterHandler(bridgeToken.EventExchangeToken, m.eventPublicChainExchangeTokenHandler)
	m.chainHandler.PbcHandler.BridgeTokenContract.RegisterHandler(bridgeToken.EventExchangeNFT, m.eventPublicChainExchangeNFTHandler)
}

// TODO
// error handler for failed exchange
// just a simple dealing solutio now

func (m *Manager) errorHandlerForPrivateChainTokenExchange(user common.Address, amount *big.Int) {
	// pay token back to private chain
	managerAccount := m.chainHandler.PvcHandler.ManagerAccount
	pvcBridgeTokenContract := m.chainHandler.PvcHandler.BridgeTokenContract
	txErr := managerAccount.SendFunction(pvcBridgeTokenContract,
		nil,
		bridgeToken.FuncPayToken,
		user,
		amount)
	err := <-txErr
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (m *Manager) errorHandlerForPublicChainTokenExchange(user common.Address, amount *big.Int) {
	// pay token back to public chain
	managerAccount := m.chainHandler.PbcHandler.ManagerAccount
	pbcBridgeTokenContract := m.chainHandler.PbcHandler.BridgeTokenContract
	txErr := managerAccount.SendFunction(pbcBridgeTokenContract,
		nil,
		bridgeToken.FuncPayToken,
		user,
		amount)
	err := <-txErr
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (m *Manager) errorHandlerForPrivateChainNFTExchange(tokenId *big.Int, owner common.Address, state *gameToken.AvatarState) {
	// pay avatar back to private chain
	managerAccount := m.chainHandler.PvcHandler.ManagerAccount
	pvcBridgeTokenContract := m.chainHandler.PvcHandler.BridgeTokenContract
	txErr := managerAccount.SendFunction(pvcBridgeTokenContract,
		nil,
		bridgeToken.FuncPayNFT,
		owner,
		state.Gene,
		state.AvatarLevel,
		state.Weaponed,
		state.Armored,
	)
	err := <-txErr
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (m *Manager) errorHandlerForPublicChainNFTExchange(tokenId *big.Int, owner common.Address, state *gameToken.AvatarState) {
	// pay avatar back to public chain
	managerAccount := m.chainHandler.PbcHandler.ManagerAccount
	pvcBridgeTokenContract := m.chainHandler.PbcHandler.BridgeTokenContract
	txErr := managerAccount.SendFunction(pvcBridgeTokenContract,
		nil,
		bridgeToken.FuncPayNFT,
		owner,
		state.Gene,
		state.AvatarLevel,
		state.Weaponed,
		state.Armored,
	)
	err := <-txErr
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (m *Manager) eventPrivateChainExchangeTokenHandler(data []byte) {
	log.Println("get a token exchange event from private chain")
	var gcuid uint64 = GCUID_EXCHANGE
	// fetch data and eth client
	exchangeTokenEventPayLoad := &bridgeToken.BridgeTokenEventExchangeToken{}
	pvcBridgeTokenContract := m.chainHandler.PvcHandler.BridgeTokenContract
	err := pvcBridgeTokenContract.Unpack(exchangeTokenEventPayLoad, bridgeToken.EventExchangeToken, data)
	if err != nil {
		panic(err)
	}

	account := exchangeTokenEventPayLoad.User.Hex()
	c, ok := m.loginClients[account]
	if !ok || c == nil { // TODO user can disconnect here
		log.Println("unknown user from exchange private token")
	}

	// pay token in public chain
	managerAccount := m.chainHandler.PbcHandler.ManagerAccount
	pbcBridgeTokenContract := m.chainHandler.PbcHandler.BridgeTokenContract
	txErr := managerAccount.SendFunction(pbcBridgeTokenContract,
		nil,
		bridgeToken.FuncPayToken,
		exchangeTokenEventPayLoad.User,
		exchangeTokenEventPayLoad.Amount)
	err = <-txErr
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		m.errorHandlerForPrivateChainTokenExchange(
			exchangeTokenEventPayLoad.User,
			exchangeTokenEventPayLoad.Amount,
		)
		return
	}
	m.wrapperAndSend(c, gcuid, &ExchangeResponse{
		Source:          SOURCE_PRIVATE_CHAIN,
		Type:            EXCHANGE_TYPE_ERC20TOKEN,
		State:           EXCHANGE_STATE_FINISH,
		Amount:          exchangeTokenEventPayLoad.Amount,
		SuccessResponse: NewSuccessResponse(gcuid),
	})
}

func (m *Manager) eventPublicChainExchangeTokenHandler(data []byte) {
	var gcuid uint64 = GCUID_EXCHANGE
	// fetch data and eth client
	log.Println("get a token exchange event from public chain")
	exchangeTokenEventPayLoad := &bridgeToken.BridgeTokenEventExchangeToken{}
	pbcBridgeTokenContract := m.chainHandler.PbcHandler.BridgeTokenContract
	err := pbcBridgeTokenContract.Unpack(exchangeTokenEventPayLoad, bridgeToken.EventExchangeToken, data)
	if err != nil {
		panic(err)
		return
	}
	account := exchangeTokenEventPayLoad.User.Hex()
	c, ok := m.loginClients[account] // TODO user can disconnect here
	if !ok || c == nil {
		log.Println("unknown user from exchange public token")
	}

	// pay token in private chain
	log.Println("pay token to private chain")
	managerAccount := m.chainHandler.PvcHandler.ManagerAccount
	pvcBridgeTokenContract := m.chainHandler.PvcHandler.BridgeTokenContract
	txErr := managerAccount.SendFunction(pvcBridgeTokenContract,
		nil,
		bridgeToken.FuncPayToken,
		exchangeTokenEventPayLoad.User,
		exchangeTokenEventPayLoad.Amount)
	err = <-txErr
	if err != nil {
		log.Println(err.Error())
		m.errorHandlerForPublicChainTokenExchange(
			exchangeTokenEventPayLoad.User,
			exchangeTokenEventPayLoad.Amount,
		)
		m.errorHandler(c, gcuid, err)
		return
	}
	m.wrapperAndSend(c, gcuid, &ExchangeResponse{
		Source:          SOURCE_PUBLIC_CHAIN,
		Type:            EXCHANGE_TYPE_ERC20TOKEN,
		State:           EXCHANGE_STATE_FINISH,
		Amount:          exchangeTokenEventPayLoad.Amount,
		SuccessResponse: NewSuccessResponse(gcuid),
	})
}

func (m *Manager) eventPrivateChainExchangeNFTHandler(data []byte) {
	var gcuid uint64 = GCUID_EXCHANGE
	// fetch data and eth client
	exchangeNFTEventPayLoad := &bridgeToken.BridgeTokenEventExchangeNFT{}
	pvcBridgeTokenContract := m.chainHandler.PvcHandler.BridgeTokenContract
	err := pvcBridgeTokenContract.Unpack(exchangeNFTEventPayLoad, bridgeToken.EventExchangeNFT, data)
	if err != nil {
		panic(err)
	}
	account := exchangeNFTEventPayLoad.Owner.Hex()
	c, ok := m.loginClients[account]
	if !ok || c == nil { // TODO user can disconnect here
		log.Println("unknown user from exchange private avatar")
	}

	// pay nft in public chain
	managerAccount := m.chainHandler.PbcHandler.ManagerAccount
	pbcBridgeTokenContract := m.chainHandler.PbcHandler.BridgeTokenContract
	txErr := managerAccount.SendFunction(pbcBridgeTokenContract,
		nil,
		bridgeToken.FuncPayNFT,
		exchangeNFTEventPayLoad.TokenID,
		exchangeNFTEventPayLoad.Owner,
		exchangeNFTEventPayLoad.Gene,
		exchangeNFTEventPayLoad.AvatarLevel,
		exchangeNFTEventPayLoad.Weaponed,
		exchangeNFTEventPayLoad.Armored)
	err = <-txErr
	avatarState := &gameToken.AvatarState{
		Gene:        exchangeNFTEventPayLoad.Gene,
		AvatarLevel: exchangeNFTEventPayLoad.AvatarLevel,
		Weaponed:    exchangeNFTEventPayLoad.Weaponed,
		Armored:     exchangeNFTEventPayLoad.Armored,
	}
	if err != nil {
		log.Println(err.Error())
		m.errorHandlerForPrivateChainNFTExchange(
			exchangeNFTEventPayLoad.TokenID,
			exchangeNFTEventPayLoad.Owner,
			avatarState,
		)
		m.errorHandler(c, gcuid, err)
		return
	}
	m.wrapperAndSend(c, gcuid, &ExchangeResponse{
		Source:          SOURCE_PRIVATE_CHAIN,
		Type:            EXCHANGE_TYPE_AVATAR,
		State:           EXCHANGE_STATE_FINISH,
		AvatarState:     avatarState,
		TokenId:         exchangeNFTEventPayLoad.TokenID.String(),
		SuccessResponse: NewSuccessResponse(gcuid),
	})
}

func (m *Manager) eventPublicChainExchangeNFTHandler(data []byte) {
	var gcuid uint64 = GCUID_EXCHANGE
	// fetch data and eth client
	exchangeNFTEventPayLoad := &bridgeToken.BridgeTokenEventExchangeNFT{}
	pbcBridgeTokenContract := m.chainHandler.PbcHandler.BridgeTokenContract
	err := pbcBridgeTokenContract.Unpack(exchangeNFTEventPayLoad, bridgeToken.EventExchangeNFT, data)
	if err != nil {
		panic(err)
	}
	account := exchangeNFTEventPayLoad.Owner.Hex()
	c, ok := m.loginClients[account] // TODO user can disconnect here
	if !ok || c == nil {
		log.Println("unknown user from exchange public")
		return
	}

	// pay nft in private chain
	managerAccount := m.chainHandler.PvcHandler.ManagerAccount
	pvcBridgeTokenContract := m.chainHandler.PvcHandler.BridgeTokenContract
	txErr := managerAccount.SendFunction(pvcBridgeTokenContract,
		nil,
		bridgeToken.FuncPayNFT,
		exchangeNFTEventPayLoad.TokenID,
		exchangeNFTEventPayLoad.Owner,
		exchangeNFTEventPayLoad.Gene,
		exchangeNFTEventPayLoad.AvatarLevel,
		exchangeNFTEventPayLoad.Weaponed,
		exchangeNFTEventPayLoad.Armored)

	avatarState := &gameToken.AvatarState{
		Gene:        exchangeNFTEventPayLoad.Gene,
		AvatarLevel: exchangeNFTEventPayLoad.AvatarLevel,
		Weaponed:    exchangeNFTEventPayLoad.Weaponed,
		Armored:     exchangeNFTEventPayLoad.Armored,
	}
	err = <-txErr
	if err != nil {
		log.Println(err.Error())
		m.errorHandlerForPublicChainNFTExchange(
			exchangeNFTEventPayLoad.TokenID,
			exchangeNFTEventPayLoad.Owner,
			avatarState,
		)
		m.errorHandler(c, gcuid, err)
	}
	m.wrapperAndSend(c, gcuid, &ExchangeResponse{
		Source:          SOURCE_PUBLIC_CHAIN,
		Type:            EXCHANGE_TYPE_AVATAR,
		State:           EXCHANGE_STATE_FINISH,
		AvatarState:     avatarState,
		TokenId:         exchangeNFTEventPayLoad.TokenID.String(),
		SuccessResponse: NewSuccessResponse(gcuid),
	})
}
