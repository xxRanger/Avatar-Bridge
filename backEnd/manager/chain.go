package manager

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xxRanger/blockchainUtil/chain"
	"github.com/xxRanger/blockchainUtil/contract/bridgeToken"
	"github.com/xxRanger/blockchainUtil/sender"
	"io/ioutil"
)

const (
	DEFAULT_TRANSFER_VALUE = "100000000000000" // 0.0001 ether
)


type PublicChainHandler struct {
	ManagerAccount      *sender.User
	Client              *chain.EthClient
	BridgeTokenContract *bridgeToken.BridgeToken
}

func (h *PublicChainHandler) Subscribe()  {
	go h.BridgeTokenContract.Subscribe()
}

type PrivateChainHandler struct {
	ManagerAccount      *sender.User
	Client              *chain.EthClient
	BridgeTokenContract *bridgeToken.BridgeToken
}

func (h *PrivateChainHandler) Subscribe()  {
	go h.BridgeTokenContract.Subscribe()
}

type ChainHandler struct {
	PbcHandler *PublicChainHandler
	PvcHandler *PrivateChainHandler
}

type AccountConfig struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

type PrivateChainConfig struct {
	Account            AccountConfig `json:"account"`
	Port               string        `json:"port"`
	BridgeTokenAddress string        `json:"bridgeTokenAddress"`
}

type PublicChainConfig struct {
	Account            AccountConfig `json:"account"`
	Port               string        `json:"port"`
	BridgeTokenAddress string        `json:"bridgeTokenAddress"`
}

func loadFile(file string, v interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	return err
}

func NewPublicChainHandler(file string) (*PublicChainHandler, error) {
	config := &PublicChainConfig{}
	err := loadFile(file, config)
	if err!=nil {
		panic(err)
	}

	// eth client
	pbcClient, err := chain.NewEthClient(config.Port)
	if err != nil {
		panic(err)
	}
	// manager account
	pk, err := crypto.HexToECDSA(config.Account.PrivateKey)
	if err != nil {
		panic(err)
	}
	address := common.HexToAddress(config.Account.Address)
	managerAccount := sender.NewUser(address, pk)
	managerAccount.BindEthClient(pbcClient,sender.CHAIN_KIND_PUBLIC)

	// eth contract
	bridgeTokenContract := bridgeToken.NewBridgeToken(common.HexToAddress(config.BridgeTokenAddress))
	bridgeTokenContract.BindClient(pbcClient)

	handler := &PublicChainHandler{
		ManagerAccount:      managerAccount,
		BridgeTokenContract: bridgeTokenContract,
		Client:              pbcClient,
	}
	return handler, nil
}

func NewPrivateChainHandler(file string) (*PrivateChainHandler, error) {
	config := &PrivateChainConfig{}
	err := loadFile(file, config)
	if err != nil {
		panic(err)
	}
	// eth client
	pvcClient, err := chain.NewEthClient(config.Port)
	if err != nil {
		panic(err)
	}
	// manager account
	pk, err := HexToPrivateKey(config.Account.PrivateKey)
	if err != nil {
		panic(err)
	}
	address := common.HexToAddress(config.Account.Address)
	managerAccount := sender.NewUser(address, pk)
	managerAccount.BindEthClient(pvcClient,sender.CHAIN_KIND_PRIVATE)

	bridgeTokenContract := bridgeToken.NewBridgeToken(common.HexToAddress(config.BridgeTokenAddress))
	bridgeTokenContract.BindClient(pvcClient)

	handler := &PrivateChainHandler{
		ManagerAccount:      managerAccount,
		BridgeTokenContract: bridgeTokenContract,
		Client:              pvcClient,
	}
	return handler, nil
}

func NewChainHandler(pvcHandler *PrivateChainHandler, pbcHandler *PublicChainHandler) *ChainHandler {
	return &ChainHandler{
		PvcHandler: pvcHandler,
		PbcHandler: pbcHandler,
	}
}
