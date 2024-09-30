package server

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
	// msg  []byte

	owner string // 持有者 简单的用户管理 连接阶段初始化
}

type Message struct {
	// 消息类型
	Type int    `json:"type"`
	Info string `json:"info"`
}

// 单独响应的消息类型
const (
	Log = iota + 1
	Msg
	Err
)

func (c *Client) readPump() {
	defer CloseClient(c)

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		if initConnect(c, msg) || upBranchList(c, msg) {
			continue
		}

		var sync SyncData
		if err := json.Unmarshal(msg, &sync); err != nil {
			c.ResponseMsg(Err, err.Error())
			continue
		}
		if upEnvConfig(c, sync.Branch) || getLog(c, sync.Log) {
			continue
		}

		if sync.Item != nil {
			sync.Item.From = c.owner
		}
		handler(&sync)
	}
}

func (c *Client) writePump() {
	defer CloseClient(c)
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
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
