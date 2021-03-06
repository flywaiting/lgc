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

type taskInfo struct {
	Cur  *Task   `json:"cur"`
	Todo []*Task `json:"todo"`
	Done []*Task `json:"done"`
}

var (
	mgr      *taskMgr
	mgrMutex sync.RWMutex
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

	mgrMutex.Lock()
	defer mgrMutex.Unlock()

	t := &Task{}
	err = json.Unmarshal(data, t)
	if err != nil {
		return
	}

	t.Ip = ip
	insertTask(t)
	if t.ID == 0 {
		err = fmt.Errorf("任务添加失败")
		return
	}

	return
}

func Run() {
	mgrMutex.Lock()
	defer mgrMutex.Unlock()

	if mgr.task != nil {
		return
	}

	t := getTask()
	if t == nil {
		return
	}

	mgr.task = t
	t.ctx, t.cancel = context.WithCancel(mgr.ctx)
	go t.run()
}

func StopTask(state int, id int) {
	defer Run()

	mgrMutex.Lock()
	defer mgrMutex.Unlock()

	updateTask(id, state)

	t := mgr.task
	if t != nil && t.ID == id {
		mgr.task = nil
		t.State = state
		t.end()
	}
}

func TaskInfo() (data []byte, err error) {
	info := &taskInfo{}
	info.Cur = mgr.task
	info.Todo = todoList()
	info.Done = doneList()

	data, err = json.Marshal(info)
	return
}

func IsRunning(id int) bool {
	t := mgr.task
	return t != nil && t.ID == id
}
