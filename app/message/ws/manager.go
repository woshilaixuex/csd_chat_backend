package ws

import (
	"encoding/json"
	"log"
	"time"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2025-03-02 22:14
 */
var DefaultClientManager = NewWSClientManager()

func (wsm *WSClientManager) Start() {
	go wsm.connectionGC()
	for {
		select {
		case client := <-wsm.register:
			wsm.handleRegister(client)

		case client := <-wsm.unregister:
			wsm.handleUnregister(client)
		case msg := <-wsm.chat:
			wsm.handleChat(msg)
		}
	}
}
func (wsm *WSClientManager) SendChat(msg *Message) {
	wsm.chat <- msg
}
func (wsm *WSClientManager) Close() {

}
func (wsm *WSClientManager) handleRegister(client *WSClient) {
	wsm.Lock()
	defer wsm.Unlock()

	if oldClient, exist := wsm.clients[client.WSID]; exist {
		oldClient.Close()
		delete(wsm.clients, client.WSID)
		wsm.c_num--
	}
	wsm.clients[client.WSID] = client
	wsm.c_num++
	log.Printf("connected %v", client.WSID)
}

func (wsm *WSClientManager) handleUnregister(client *WSClient) {
	wsm.Lock()
	defer wsm.Unlock()

	oldClient, exist := wsm.clients[client.WSID]
	if !exist {
		log.Fatalf("unfound %v", client.WSID)
		return
	}

	oldClient.Close()
	delete(wsm.clients, client.WSID)
	wsm.c_num--
	log.Printf("unconnected %v", client.WSID)
}

func (m *WSClientManager) connectionGC() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.RLock()

		for CID, client := range m.clients {
			if time.Since(client.LastSeen) > 2*maxLinkTime {
				client.Close()
				delete(m.clients, CID)
				log.Printf("clear connect: %v", CID)
			}
		}
		m.Unlock()
	}
}
func (wsm *WSClientManager) handleChat(msg *Message) {
	userID := msg.SendID
	client, exist := wsm.clients[userID]
	if !exist {
		log.Fatalf("unfound client id : %v", userID)
	}
	sendMsg, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
		return
	}
	client.Send <- sendMsg
}
