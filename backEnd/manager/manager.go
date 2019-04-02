package manager

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/xxRanger/bridge/backEnd/client"
	"log"
	"net/http"
)

const CLIENT_BUFFER = 1024

type handler func(user *client.Client, gcuid uint64, data []byte)

type Manager struct {
	loginClients map[string]*client.Client
	clients      map[*client.Client]bool
	broadcast    chan []byte
	register     chan *client.Client
	unregister   chan *client.Client
	handlers     map[uint64]handler
	chainHandler *ChainHandler
}

func NewManager() *Manager {
	m := &Manager{
		loginClients: make(map[string]*client.Client),
		clients:      make(map[*client.Client]bool),
		register:     make(chan *client.Client, CLIENT_BUFFER),
		unregister:   make(chan *client.Client, CLIENT_BUFFER),
		broadcast:    make(chan []byte, CLIENT_BUFFER),
		handlers:     make(map[uint64]handler),
	}
	return m
}

func (m *Manager) RegisterHandler(gcuid uint64, h handler) {
	m.handlers[gcuid] = h
}

func (m *Manager) UnregisterHandler(gcuid uint64, h handler) {
	delete(m.handlers, gcuid)
}

func (m *Manager) SetChainHandler(chainHandler *ChainHandler) {
	m.chainHandler = chainHandler
}

func (m *Manager) InitHandlers() {
	m.handlers[GCUID_LOGIN] = m.LoginHandler
	m.handlers[GCUID_LOGOUT] = m.LogoutHandler
	m.handlers[GCUID_REWARD] = m.RewardHandler
	m.handlers[GCUID_CONSUME] = m.ConsumeHandler
	m.handlers[GCUID_EXCHANGE] = m.ExchangeHandler
	m.handlers[GCUID_GET_BALANCE] = m.GetBalanceHandler
	m.handlers[GCUID_GET_AVATAR] = m.GetAvatarStateHandler
	m.handlers[GCUID_UPGRADE_AVATAR] = m.UpgradeHanlder
}

func (m *Manager) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("receive a reqeust")
	var upgrader = websocket.Upgrader{}
	upgrader.CheckOrigin = func(rq *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer conn.Close()

	c := client.NewClient()
	c.Conn = conn
	m.connectInHandler(c)

	go c.Sender()
	for {
		_, data, err := conn.ReadMessage()
		var kvs map[string]interface{}
		if err != nil {
			log.Println(err.Error())
			break;
		}
		json.Unmarshal(data, &kvs)
		gid, ok := kvs["gcuid"]
		if !ok {
			log.Println(errors.New("gcuid not exist"))
			continue
		}
		gcuid := uint64(gid.(float64))
		if h, ok := m.handlers[gcuid]; ok {
			go h(c, gcuid, data)
		} else {
			log.Println("unknown message")
		}
	}
	m.closeHandler(c)
}

func (m *Manager) Start() {
	for {
		select {
		case c := <-m.register:
			m.clients[c] = true
			log.Println("a new user connect")
		case c := <-m.unregister:
			delete(m.clients, c)
			log.Println("a user unregister disconnect")
		case message := <-m.broadcast:
			log.Println("broadcast a message:", string(message))
			for c, _ := range m.clients {
				c.Send(message) // an active user may block other user here, fix in the future
			}
		}
	}
}

func (m *Manager) Subscribe() {
	if m.chainHandler.PvcHandler != nil {
		m.registerPrivateChainHandler()
		m.chainHandler.PvcHandler.Subscribe()
	}

	if m.chainHandler.PbcHandler != nil {
		m.registerPublicChainHandler()
		m.chainHandler.PbcHandler.Subscribe()
	}
}
