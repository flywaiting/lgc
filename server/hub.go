package server

import (
	"encoding/json"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	// sync       chan *SyncData
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
		// case client := <-h.request:
		// case sync := <-h.sync:
		// 	handler(sync)
		// 	println("after sync")
		case msg := <-h.broadcast:
			// println("hub broadcast", string(msg))
			// println("current size", len(h.clients))
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) response(msg *SyncData) {
	rsp, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.broadcast <- rsp
}
