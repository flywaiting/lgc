package server

import "github.com/gorilla/websocket"

var teamMap = make(map[string]string)

func upTeamMap(key, value string) {
	teamMap[key] = value
	// todo	更新客户端
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

var wsList []websocket.Conn

func NewWS(ws websocket.Conn) {

}
