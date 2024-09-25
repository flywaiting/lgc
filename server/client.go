package server

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
	// msg  []byte
}

type Message struct {
	Type int    `json:"type"` // 消息类型
	Info string `json:"info"`
}

// 单独响应的消息类型
const (
	Msg = iota
	Log
	Err
)

func (c *Client) readPump() {
	defer CloseClient(c)

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		if initConnect(c, msg) || getBranchList(c, msg) {
			continue
		}

		var sync SyncData
		if err := json.Unmarshal(msg, &sync); err != nil {
			c.ResponseMsg(Err, err.Error())
			continue
		}
		if upEnvConfig(c, sync.Branch) {
			continue
		}

		// c.msg = msg
		// hub.request <- c
		hub.sync <- &sync
	}
}

func (c *Client) writePump() {
	defer CloseClient(c)

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

// 文本消息
func (c *Client) ResponseMsg(t int, s string) {
	info := Message{
		Type: t,
		Info: s,
	}
	if msg, err := json.Marshal(info); err == nil {
		c.send <- msg
	}
}

// 服务数据
func (c *Client) ResponseInfo(info interface{}) {
	rsp, err := json.Marshal(info)
	if err != nil {
		return
	}
	c.send <- rsp
}
