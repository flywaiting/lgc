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
}

var tm sync.Mutex

// 新建任务
func NewTask() {
	if len(taskChan) >= 5 {

	}

	t := &Task{}
	t.Ctime = util.Tsp()
	taskChan <- t
}

// 取消任务
func CancelTask(tId int) {
	// if curTask == nil || curTask.Id != tId {
	// 	// todo	无效的取消
	// }

	// done(0)
}

// 任务列表
func TaskList() {

}

// type Task struct{}

var taskChan chan *Task

// var curTask *Task

func TaskInit(ctx context.Context) {
	taskChan = make(chan *Task, 5)
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
			break
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

	status := 1
	if err != nil {
		status = 2
	}
	done(t, status)
}

func done(t *Task, status int) {
	tm.Lock()
	defer tm.Unlock()

	if t == nil || t.Status != 0 {
		return
	}

	t.cancel()
	t.Status = status
	t.Etime = util.Tsp()
	// todo 数据写入
}
