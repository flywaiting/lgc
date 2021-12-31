package server

import (
	"container/list"
	"time"
)

// 队列
var que *list.List

// 队列默认超时时间
var qto time.Duration

// 队列 channel
var qch chan *Task

func init() {
	qto = 3 * time.Second
	que = list.New()
	go gQue()
}

// 独立 goroutine, 处理任务的增减
func gQue() {
	for t := range qch {
		// todo	添加, 移除|取消, 完成
		if t.Id == 0 {
			add(t)
		} else {
			rm(t)
		}

		if que.Len() > 0 {
			taskNext((que.Front().Value).(*Task))
		}
	}
}

// 对外统一API
func queNext(t *Task) bool {
	select {
	case qch <- t:
		return true
	case <-time.After(qto):
		return false
	}
}

func add(t *Task) {
	// todo	数量限制
	que.PushBack(t)
}
func rm(t *Task) {
	v := que.Front()
	// todo	任务完成处理

	for v != nil {
		if v.Value == t {
			que.Remove(v)
			break
		} else {
			v = v.Next()
		}
	}
}
