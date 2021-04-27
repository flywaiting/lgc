package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
)

type taskMgr struct {
	task   *Task
	ctx    context.Context
	cancel context.CancelFunc
}

type addInfo struct {
	Pattern string `json:"pattern"`
	Team    string `json:"team"`
	Branch  string `json:"branch"`
}

var (
	mgr *taskMgr
	tm  sync.RWMutex
)

func init() {
	mgr = &taskMgr{}
	mgr.ctx, mgr.cancel = context.WithCancel(context.Background())
}

func AddTask(reader io.Reader, ip string) (err error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	info := &addInfo{}
	err = json.Unmarshal(data, info)
	if err != nil {
		return
	}

	t := &Task{
		Pattern: info.Pattern,
		Team:    info.Team,
		Branch:  info.Branch,
		Ip:      ip,
	}
	t.ID = InsertTask(t)
	if t.ID == 0 {
		err = fmt.Errorf("任务添加失败")
		return
	}

	return
}

func Run() {
	// dm.Lock()
	// defer dm.Unlock()

	// if mgr.task == nil {
	// 	if len(mgr.TodoList) == 0 {
	// 		return
	// 	}

	// 	mgr.task = mgr.TodoList[0]
	// 	mgr.TodoList = mgr.TodoList[1:]
	// }

	// t := mgr.task
	// if t.cancel != nil {
	// 	return
	// }

	// t.ctx, t.cancel = context.WithCancel(mgr.ctx)
	// go t.run()
}

func StopTask(status int) {
	t := mgr.task
	if t == nil || t.cancel == nil {
		return
	}

}
