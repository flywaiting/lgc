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

// var envMap map[string]string
var envMap = sync.Map{}
var task *TaskHub
var handerMux sync.Mutex

func init() {
	hub = &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		// sync:       make(chan *SyncData),
	}
	task := &TaskHub{
		Counter: 0,
		item:    make(chan *TaskItem),
		ctx:     context.Background(),
	}
	config = util.InitConfig()
	list, err := util.BranchList(config.Git.TmpRepository)
	if err != nil {
		panic("仓库分支抓取出错")
	}
	envMap.Store("branch", list)
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
	sync := &SyncData{
		Task:    task,
		Team:    make(map[string]string),
		Pattern: config.Product.Pattern,
	}
	// team := make(map[string]string)
	for _, val := range config.Product.Teams {
		v, ok := envMap.Load(val)
		if !ok {
			sync.Team[val] = ""
		} else {
			sync.Team[val] = v.(string)
		}
	}
	if list, ok := envMap.Load("branch"); ok {
		sync.List = list.([]string)
	}
	c.ResponseInfo(sync)
	return
}

// 更新分支列表
func upBranchList(c *Client, msg []byte) (res bool) {
	if string(msg) != "branch" {
		return
	}

	res = true
	list, err := util.BranchList(config.Git.TmpRepository)
	if err != nil {
		c.ResponseMsg(Err, err.Error())
		return
	}

	hub.response(&SyncData{
		List: list,
	})
	return
}

// 配置分支
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

// 获取日志
func getLog(c *Client, id int) (res bool) {
	if id <= 0 {
		return
	}

	res = true
	log, err := util.GetLog(&config.Server, id)
	if err != nil {
		c.ResponseMsg(Err, err.Error())
		return
	}
	c.ResponseMsg(Log, string(log))
	return
}

func handler(sync *SyncData) {
	handerMux.Lock()
	defer handerMux.Unlock()

	if sync == nil {
		return
	}

	upTeam(sync.Team)

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

// 服务器分支配置
func upTeam(team map[string]string) { // c *Client,
	if team == nil {
		return
	}

	for t, b := range team {
		envMap.Store(t, b)
	}

	hub.response(&SyncData{
		Team: team,
	})
}

type SyncData struct {
	Team map[string]string `json:"team"` // 测试机分配

	// ->server
	Item   *TaskItem `json:"item"`   // 打包任务
	Branch *Branches `json:"branch"` // 配置分支
	Log    int       `json:"log"`    // 日子ID

	// ->client
	Task    *TaskHub `json:"task"`    // 任务相关信息
	List    []string `json:"list"`    // 分支列表
	Pattern []string `json:"pattern"` // 打包模式
}

// ->server 配置分支
type Branches struct {
	Branch string `json:"branch"`
	// 分支别名
	Alias string `json:"alias"`
}
