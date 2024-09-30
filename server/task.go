package server

import (
	"context"
	"sync"
	"time"
)

// type Task struct {
// 	ID      int    `json:"id" sql:"id"`
// 	State   int    `json:"state" sql:"state"`
// 	Ip      string `json:"ip" sql:"ip"`
// 	Pattern string `json:"pattern" sql:"pattern"`
// 	Team    string `json:"team" sql:"team"`
// 	Branch  string `json:"branch" sql:"branch"`
// 	CT      string
// 	ET      string
// 	ctx     context.Context
// 	cancel  context.CancelFunc
// 	// CT      util.JsonTime `json:"create"`
// 	// ET      util.JsonTime `json:"time"`
// }

type TaskHub struct {
	TodoList   []TaskItem `json:"todoList"`
	FinishList []TaskItem `json:"finishList"`
	Current    *TaskItem  `json:"current"` // 当前进行中的任务
	counter    int        // 任务计数

	item chan *TaskItem

	ctx context.Context
	sync.RWMutex
}

// 任务状态
const (
	Wait = iota // 队列
	Doing
	Git // git环节
	Finish
	Error
	Interrupt
	Cancel
)

func (t *TaskHub) enqueue(item *TaskItem) {
	t.Lock()
	defer t.Unlock()
	t.TodoList = append(t.TodoList, *item)
}
func (t *TaskHub) dequeue() *TaskItem {
	t.Lock()
	defer t.Unlock()

	if len(t.TodoList) == 0 {
		return nil
	}
	item := t.TodoList[0]
	t.TodoList = t.TodoList[1:]
	return &item
}

//	func (t *TaskHub) size() int {
//		t.RLock()
//		defer t.RUnlock()
//		return len(t.TodoList)
//	}

func (t *TaskHub) del(id int) *TaskItem {
	t.Lock()
	defer t.Unlock()

	if id < 0 || len(t.TodoList) == 0 {
		return nil
	}

	idx := -1
	for k, v := range t.TodoList {
		if v.Id == id {
			idx = k
			break
		}
	}
	if idx < 0 {
		return nil
	}
	item := t.TodoList[idx]
	t.TodoList = append(t.TodoList[:idx], t.TodoList[idx+1:]...)
	return &item
}

func (t *TaskHub) push(item *TaskItem) {
	t.Lock()
	defer t.Unlock()

	cnt := len(t.FinishList)
	if cnt > 10 {
		t.FinishList = t.FinishList[cnt-10:]
	}
	t.FinishList = append(t.FinishList, *item)
}

func (t *TaskHub) run() {
	// for {
	// 	select {
	// 	case item := <-t.item:
	// 		t.Counter++
	// 		item.Id = t.Counter
	// 		item.CreateTime = time.Now().Unix()
	// 		t.enqueue(item)
	// 		t.next()
	// 	case id := <-t.id:
	// 		t.remove(*id)
	// 	}
	// }

	for item := range t.item {
		if item.Id > 0 {
			t.remove(item.Id)
			continue
		}

		if alias, ok := envMap.Load(item.Branch); ok {
			item.alias = alias.(string)
		} else {
			item.alias = item.Branch
		}
		item.Id = t.counter
		t.counter++
		item.CreateTime = time.Now().Unix()
		item.Status = Wait
		t.enqueue(item)
		t.next()
	}
}

// func (t *TaskHub) handler(item *TaskItem) {
// 	if item.Id > 0 {
// 		t.remove(item.Id)
// 		return
// 	}

// 	t.item <- item

// 	// 新任务
// 	if alias, ok := envMap.Load(item.Branch); ok {
// 		item.alias = alias.(string)
// 	}
// }

func (t *TaskHub) next() {
	if t.Current != nil {
		hub.response(&SyncData{
			Task: t,
		})
		return
	}

	item := t.dequeue()
	if item == nil {
		hub.response(&SyncData{
			Task: t,
		})
		return
	}

	item.Status = Doing
	item.ActiveTime = time.Now().Unix()
	t.Current = item
	item.ctx, item.fn = context.WithCancel(t.ctx)
	go item.run()

	hub.response(&SyncData{
		Task: t,
	})
}

func (t *TaskHub) remove(id int) {
	cur := t.Current
	if cur != nil && cur.Id == id {
		if cur.Status >= Git {
			// c.ResponseMsg(Msg, "git阶段不能中断")
			return
		}

		cur.finish(Interrupt)
		t.push(cur)
		t.Current = nil
		t.next()
		return
	}

	item := t.del(id)
	if item == nil {
		return
	}
	item.finish(Cancel)
	// t.FinishList = append(t.FinishList, *item)
	t.push(item)
	t.next()
}

func (t *TaskHub) finish(state int) {
	item := t.Current
	t.Current = nil
	if item != nil {
		item.finish(state)
		t.push(item)
	}
	t.next()
}

// ---------------

// func (t *Task) run() {
// 	if t.State == com.Running || t.cancel == nil {
// 		return
// 	}

// 	t.State = com.Running
// 	updateTask(t.ID, t.State)

// 	logPath := com.LogPath(t.ID)
// 	log, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		StopTask(com.Interrupt, t.ID)
// 		return
// 	}
// 	defer log.Close()

// 	tasks := []string{fmt.Sprintf("mo t %s ---%s@%s", t.Pattern, t.Team, t.Branch)}
// 	tasks = append(tasks, com.SufTask()...)
// 	for _, v := range tasks {
// 		if t.State != com.Running {
// 			return
// 		}

// 		bash := "bash"
// 		par := "-c"
// 		if com.OS == "windows" {
// 			bash = "cmd"
// 			par = "/c"
// 		}

// 		cmd := exec.CommandContext(t.ctx, bash, par, v)
// 		cmd.Stderr = log
// 		cmd.Stdout = log
// 		cmd.Dir = com.WkDir()

// 		err = cmd.Run()
// 		if err != nil {
// 			fmt.Println("running:", err.Error())
// 			if t.State < com.Succ {
// 				StopTask(com.Interrupt, t.ID)
// 			}
// 			return
// 		}

// 		code := 0
// 		if cmd.ProcessState != nil {
// 			code = cmd.ProcessState.ExitCode()
// 		}
// 		if code != 0 {
// 			fmt.Println("coed error: ", code)
// 			StopTask(com.Interrupt, t.ID)
// 		}
// 	}

// 	StopTask(com.Succ, t.ID)
// }

// func (t *Task) end() {
// 	if t.cancel == nil {
// 		return
// 	}
// 	t.cancel()
// 	t.cancel = nil
// }
