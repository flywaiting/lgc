package server

import (
	"context"
	"lgc/util"
	"os/exec"
	"sync"
)

type Task struct {
	Status int `db:"tstatus"`
	Id     int `db:"tid"`
	Ctime  int64
	Stime  int64
	Etime  int64
	Rname  string
	Cmds   string
	cancel context.CancelFunc

	nxt *Task
}

var tm sync.Mutex

// 新建任务
func NewTask() {
	if que.size >= limit {
		return
	}

	t := &Task{}
	t.Ctime = util.Tsp()
	taskChan <- t
}

// 取消任务
func CancelTask(tId int) {
	t := que.find(tId)
	if t == nil {
		return
	}
	done(t, util.Cancel)
}

// 任务列表
func TaskList() {

}

// type Task struct{}

var taskChan chan *Task
var que *queTask

// var curTask *Task

func TaskInit(ctx context.Context) {
	taskChan = make(chan *Task, limit)
	que = NewQue()
	go shedule(ctx)
}

// 任务排班
func shedule(ctx context.Context) {
	for {
		select {
		case t := <-taskChan:
			run(ctx, t)
		case <-ctx.Done():
			close(taskChan)
			return
		}
	}
}

func run(ctx context.Context, t *Task) {
	c, cancel := context.WithCancel(ctx)

	t.cancel = cancel
	t.Stime = util.Tsp()
	cmd := exec.CommandContext(c, "a")
	// todo	输出设置
	err := cmd.Run()

	var status int
	if err != nil {
		status = util.Wrong
	} else {
		status = util.Done
	}
	done(t, status)
}

func done(t *Task, status int) {
	tm.Lock()
	defer tm.Unlock()

	if t == nil || t.Status != util.Ready {
		return
	}

	t.cancel()
	t.Status = status
	t.Etime = util.Tsp()
	// todo 数据写入
}
