package server

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Message struct {
	Msg string `json:"msg"`
}

func (c *Client) readPump() {
	defer CloseClient(c)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		RequestHandle(c, message)
	}
}

func (c *Client) writePump() {
	defer CloseClient(c)

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

func (c *Client) ResponseMsg(s string) {
	info := Message{
		Msg: s,
	}
	if msg, err := json.Marshal(info); err == nil {
		c.send <- msg
	}
}

func (c *Client) ResponseInfo(info interface{}) {
	rsp, err := json.Marshal(info)
	if err != nil {
		return
	}
	c.send <- rsp
}
