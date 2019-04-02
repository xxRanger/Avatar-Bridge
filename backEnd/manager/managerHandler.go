package manager

import (
	"encoding/json"
	"github.com/xxRanger/bridge/backEnd/client"
	"log"
)

func (m *Manager) clientLogoutHandler(c *client.Client) {
	delete(m.loginClients, c.Address.Hex())
}

func (m *Manager) clientLoginHandler(c *client.Client) error {
	//_,loggedInBefore:=m.loginClients[c.Address.Hex()]   // TODO comment for test
	m.loginClients[c.Address.Hex()] = c

	//if !loggedInBefore { // transfer ether to user have not login before
	//	value,_:=new(big.Int).SetString(DEFAULT_TRANSFER_VALUE,10)
	//	txError:=m.chainHandler.PbcHandler.ManagerAccount.Transfer(c.Address,value)
	//	err:= <-txError
	//	if err!=nil {
	//		return err
	//	}
	//}    // TODO comment for test
	return nil
}

func (m *Manager) Register(c *client.Client) {
	m.register <- c
}

func (m *Manager) Unregister(c *client.Client) {
	m.unregister <- c
}

func (m *Manager) closeHandler(c *client.Client) {
	c.Close()
	m.Unregister(c)
}

func (m *Manager) connectInHandler(c *client.Client) {
	m.Register(c)
}

func (m *Manager) wrapperAndSend(c *client.Client, gcuid uint64, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Println(err.Error())
		m.errorHandler(c, gcuid, err)
		return
	}

	c.Send(data)
}