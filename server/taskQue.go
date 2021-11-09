package server

import "sync"

type queTask struct {
	head, tail *Task
	size       int
}

var mut sync.RWMutex
var limit = 5

func (q *queTask) add(t *Task) bool {
	mut.Lock()
	defer mut.Unlock()

	if q.size >= limit {
		return false
	}

	q.tail.nxt = t
	q.tail = t
	q.size++
	return true
}

// 查找
func (q *queTask) find(tid int) *Task {
	mut.RLock()
	mut.RUnlock()

	t := q.head.nxt
	for t != nil {
		if t.Id == tid {
			return t
		}
		t = t.nxt
	}
	return nil
}

// 移除
func (q *queTask) rm(tid int) *Task {
	mut.Lock()
	defer mut.Unlock()

	t := q.head
	for t != nil && t.nxt != nil {
		if t.nxt.Id == tid {
			v := t.nxt
			t.nxt = v.nxt
			v.nxt = nil
			return v
		}
	}
	return nil
}

func NewQue() *queTask {
	q := &queTask{}
	q.head = &Task{}
	q.tail = q.head
	return q
}
