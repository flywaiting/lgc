package server

import (
	"context"
	"sync"
	"time"
)

var queue *taskNode
var taskChan chan *Task
var taskSize int = 5

var mutex sync.Mutex

func init() {
	taskChan = make(chan *Task, taskSize)
	queue = &taskNode{}
	go handleTask()
}

func handleTask() {
	for t := range taskChan {
		if t.Status == 0 {
			ctx, cancel := context.WithCancel(context.Background())
			t.Fn = cancel
			t.StartTsp = time.Now().Unix()

			go t.Run(ctx)
			<-ctx.Done()
		}
		queue.rm(t.Id)
	}
}

func Add(t *Task) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if len(taskChan) < taskSize {
		queue.add(t)
		t.AddTsp = time.Now().Unix()
		taskChan <- t
		return true
	}
	return false
}

func Rm(tid int) {
	mutex.Lock()
	defer mutex.Unlock()

	t := queue.rm(tid)
	if t != nil {
		t.Done(1)
	}
}
