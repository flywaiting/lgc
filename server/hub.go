package server

import "encoding/json"

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	// request    chan *Client
	sync      chan *SyncData
	broadcast chan []byte
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
		case sync := <-h.sync:
			handler(sync)
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

func (h *Hub) response(msg *SyncData) {
	rsp, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.broadcast <- rsp
}
