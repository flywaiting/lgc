package server

import (
	"context"
	"encoding/json"
	"fmt"
	"lgc/util"
	"regexp"
	"sync"

	"github.com/gorilla/websocket"
)

var hub *Hub
var config *util.Config

// var teamMap map[string]string
var teamMap = sync.Map{}
var task *TaskHub

func init() {
	hub = &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	task := &TaskHub{
		Counter: 0,
		item:    make(chan *TaskItem),
		ctx:     context.Background(),
	}
	config = util.InitConfig()
	// teamMap = make(map[string]string, 0)

	go hub.run()
	go task.run()
}

func ServerWs(ws *websocket.Conn) {
	client := &Client{
		conn: ws,
		send: make(chan []byte),
	}

	hub.register <- client
	go client.readPump()
	go client.writePump()

	// todo 初始化
}

func CloseClient(c *Client) {
	c.conn.Close()
	hub.unregister <- c
}

func handler(c *Client) {
	if c.msg == nil {
		return
	}

	var sync SyncData
	if err := json.Unmarshal(c.msg, &sync); err != nil {
		c.ResponseMsg(Err, err.Error())
		return
	}

	switch sync.Action {
	case actTeam:
		upTeam(sync.Team)
	case actBranch:
		upBranch(c, sync.Branch)
	case actTask:
		task.handler(c, sync.Item)
	}
}

func upBranch(c *Client, b *Branches) {
	if b == nil {
		return
	}
	// 获取分支列表
	if b.Branch == "" {
		list, err := util.BranchList(config.Git.TmpRepository)
		if err != nil {
			c.ResponseMsg(Err, err.Error())
			return
		}
		b.List = list
		c.ResponseInfo(&SyncData{
			Branch: b,
		})
		return
	}

	// 添加分支
	matched, err := regexp.MatchString(`[a-zA-Z]`, b.Branch)
	if err != nil {
		c.ResponseMsg(Err, "分支正则匹配出错了")
		return
	}
	if !matched {
		b.Alias = b.Branch
	}

	if b.Alias == "" {
		c.ResponseMsg(Err, "分支需要个别名")
		return
	}

	// 分支记录模板
	item := fmt.Sprintf(`"%s": {verIndex:1, branch:{common:"%s", client:"%s", server:"%s", art:"%s"}, },`, b.Alias, b.Branch, b.Branch, b.Branch, b.Branch)
	err = util.UpFile(&config.Product, item)
	if err != nil {
		c.ResponseMsg(Err, err.Error())
		return
	}
}

func upTeam(team map[string]string) { // c *Client,
	if team == nil {
		return
	}
	for t, b := range team {
		teamMap.Store(t, b)
	}

	hub.response(&SyncData{
		Team: team,
	})
}

type SyncData struct {
	Action int               `json:"action"`
	Team   map[string]string `json:"team"`
	Branch *Branches         `json:"branch"`
	Task   *TaskHub          `json:"task"` // ->client, 任务相关信息
	Item   *TaskItem         `json:"item"` // ->server, 任务操作
}
type Branches struct {
	// server -> 当前所有分支
	List []string `json:"list"`
	// client -> 新增分支
	Branch string `json:"branch"`
	// 新分支别名
	Alias string `json:"alias"`
}

const (
	Init = iota
	actTeam
	actBranch
	actTask
)
