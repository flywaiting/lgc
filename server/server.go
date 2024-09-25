package server

import (
	"context"
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
}

func CloseClient(c *Client) {
	c.conn.Close()
	hub.unregister <- c
}

// websocket 连接初始化
func initConnect(c *Client, msg []byte) (res bool) {
	if string(msg) != "init" {
		return
	}

	res = true
	team := make(map[string]string)
	for _, val := range config.Product.Teams {
		v, ok := teamMap.Load(val)
		if !ok {
			team[val] = ""
		} else {
			team[val] = v.(string)
		}
	}
	c.ResponseInfo(&SyncData{
		Team: team,
		Task: task,
	})
	return
}

// 获取分支列表
func getBranchList(c *Client, msg []byte) (res bool) {
	if string(msg) != "branch" {
		return
	}

	res = true
	list, err := util.BranchList(config.Git.TmpRepository)
	if err != nil {
		c.ResponseMsg(Err, err.Error())
		return
	}
	c.ResponseInfo(&SyncData{
		List: list,
	})
	return
}

func upEnvConfig(c *Client, b *Branches) (res bool) {
	if b == nil {
		return
	}

	res = true
	matched, err := regexp.MatchString(`[a-zA-Z]`, b.Branch)
	if err != nil {
		c.ResponseMsg(Err, "分支正则匹配出错了, 再试试")
		return
	}
	if !matched {
		b.Alias = b.Branch
	}
	if b.Alias == "" {
		c.ResponseMsg(Err, "分支需要给别名")
		return
	}
	if util.CheckKeyExist(&config.Product, fmt.Sprintf(`"%s"`, b.Branch)) {
		return
	}
	item := fmt.Sprintf(`"%s": {verIndex:1, branch:{common:"%s", client:"%s", server:"%s", art:"%s"}, },`, b.Alias, b.Branch, b.Branch, b.Branch, b.Branch)
	err = util.UpFile(&config.Product, item)
	if err != nil {
		c.ResponseMsg(Err, err.Error())
	}
	return
}

func handler(sync *SyncData) {
	// if c.msg == nil {
	// 	return
	// }

	// var sync SyncData
	// if err := json.Unmarshal(c.msg, &sync); err != nil {
	// 	c.ResponseMsg(Err, err.Error())
	// 	return
	// }

	// switch sync.Action {
	// case actTeam:
	// 	upTeam(sync.Team)
	// case actBranch:
	// 	upBranch(c, sync.Branch)
	// case actTask:
	// 	task.handler(c, sync.Item)
	// }
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
	Team map[string]string `json:"team"` // 测试机分配

	// -> server
	Item   *TaskItem `json:"item"`   // 打包任务
	Branch *Branches `json:"branch"` // 配置分支

	// -> client
	Task *TaskHub `json:"task"` // ->client, 任务相关信息
	List []string `json:"list"` // ->client, 分支列表
}
type Branches struct {
	// // server -> 当前所有分支
	// List []string `json:"list"`
	// client -> 新增分支
	Branch string `json:"branch"`
	// 新分支别名
	Alias string `json:"alias"`
}

// const (
// 	Init = iota
// 	actTeam
// 	actBranch
// 	actTask
// )
