package server

import (
	"fmt"
	"lgc/util"
	"regexp"

	"github.com/gorilla/websocket"
)

var hub *Hub
var config *util.Config
var teamMap map[string]string

func init() {
	hub = &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	config = util.InitConfig()
	teamMap = make(map[string]string, 0)

	go hub.run()
}

func NewWS(ws *websocket.Conn) {
	client := &Client{
		conn: ws,
		send: make(chan []byte),
	}

	// todo 初始化

	client.readPump()
	client.writePump()
}

func CloseClient(c *Client) {
	c.conn.Close()
	hub.unregister <- c
}

func RequestHandle(c *Client, message []byte) {

}

// 分支相关
func upBranch(c *Client, b *Branches) {
	// 获取分支列表
	if b.Branch == "" {
		list, err := util.BranchList(config.Git.TmpRepository)
		if err != nil {
			return
		}
		b.List = list
		c.ResponseInfo(&Request{
			Branch: b,
		})
		return
	}

	// 添加分支
	matched, err := regexp.MatchString(`[a-zA-Z]`, b.Branch)
	if err != nil {
		c.ResponseMsg("分支正则匹配出错了")
		return
	}
	if !matched {
		b.Alias = b.Branch
	}

	if b.Alias == "" {
		c.ResponseMsg("分支需要个别名")
		return
	}

	// 分支记录模板
	item := fmt.Sprintf(`"%s": {verIndex:1, branch:{common:"%s", client:"%s", server:"%s", art:"%s"}, },`, b.Alias, b.Branch, b.Branch, b.Branch, b.Branch)
	// todo 将格式化的分支写入文件
	fmt.Println(item)
}

func upTeam(team map[string]string) { // c *Client,
	for t, b := range team {
		teamMap[t] = b
	}

	hub.response(&Request{
		Team: teamMap,
	})
}

type Request struct {
	Action int               `json:"action"`
	Team   map[string]string `json:"team"`
	Branch *Branches         `json:"branch"`
	Task   Task              `json:"task"`
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
	UpTeams
	UpBranches
	UpTasks
)
