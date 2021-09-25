package server

import (
	"context"
	"sync"
	"time"
)

type Task struct {
	Status   byte
	Id       int
	AddTsp   int64
	StartTsp int64
	EndTsp   int64
	Uid      string
	CmdList  []string
	Fn       context.CancelFunc
}

var tm sync.Mutex

func (t *Task) Run(ctx context.Context) {

}

func (t *Task) Done(status byte) {
	tm.Lock()
	defer tm.Unlock()

	if t.Status == 0 {
		t.Status = status
		t.EndTsp = time.Now().Unix()
	}

	if t.Fn != nil {
		t.Fn()
	}
}
