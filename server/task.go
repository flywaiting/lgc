package server

import (
	"context"
	"lgc/util"
	"os/exec"
	"time"
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

// 任务 channel
var tch chan *Task

// 任务默认超时
var tto time.Duration

func init() {
	tto = 1 * time.Second
	tch = make(chan *Task)
	go gTask()
}

func gTask() {
	for t := range tch {
		run(t)
	}
}

func taskNext(t *Task) bool {
	select {
	case tch <- t:
		return true
	case <-time.After(tto):
		return false
	}
}

func run(t *Task) {
	// todo	判定任务的有效性
	
	defer queNext(t)

	ctx, cancel := context.WithCancel(context.TODO())
	t.cancel = cancel
	t.Stime = util.Tsp()
	// todo	输出记录
	// todo	任务执行, 走配置
	cmd := exec.CommandContext(ctx, "c")
	err := cmd.Run()
	if err != nil || (cmd.ProcessState != nil && cmd.ProcessState.ExitCode() != 0) {
		// todo	task abort
		return
	}
	// todo finish
}

func done(t *Task, status int) {
	// tm.Lock()
	// defer tm.Unlock()

	// if t == nil || t.Status != util.Ready {
	// 	return
	// }

	// t.cancel()
	// t.Status = status
	// t.Etime = util.Tsp()
	// // todo 数据写入
}
