package task

import (
	"context"
	"sync"
	"time"
)

type Task struct {
	Id       int
	Status   int
	AddTsp   int
	StartTsp int
	EndTsp   int
	Uid      string
	CmdList  []string
	Ctx      context.Context
	Fn       context.CancelFunc
}

var taskChan chan *Task
var taskSize int = 5

var wait sync.WaitGroup
var mutex sync.Mutex

func init() {
	taskChan = make(chan *Task, taskSize)
	go taskQueue()
}

func Add(t *Task) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if len(taskChan) < taskSize {
		t.AddTsp = int(time.Now().Unix())
		taskChan <- t
		return true
	}
	return false
}

func taskQueue() {
	for t := range taskChan {
		ctx, cancel := context.WithCancel(context.Background())
		t.Ctx = ctx
		t.Fn = cancel
		t.StartTsp = int(time.Now().Unix())

		wait.Add(1)
		go run(t)
		wait.Wait()
	}
}

func run(t *Task) {
	done()
}

func done() {
	wait.Done()
}
