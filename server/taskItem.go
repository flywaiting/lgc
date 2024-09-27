package server

import (
	"context"
	"time"
)

type TaskItem struct {
	Id         int      `json:"id"`
	From       string   `json:"from"`   // 发起人
	To         []string `json:"to"`     // 完成通知
	Pattern    string   `json:"patten"` // 任务类型
	Team       string   `json:"team"`
	Branch     string   `json:"branch"`
	Status     int      `json:"status"`
	CreateTime int64    `json:"createTime"`
	ActiveTime int64    `json:"activeTime"` // 开始执行时间戳
	FinishTime int64    `json:"finishTime"`

	Log int `json:"log"` // 日志查询

	ctx context.Context
	fn  context.CancelFunc
}

func (t *TaskItem) run() {
}

// 任务结束
func (t *TaskItem) finish(status int) {
	t.Status = status
	t.FinishTime = time.Now().Unix()
	if t.fn == nil {
		return
	}

	t.fn()
	t.fn = nil
	t.ctx = nil
}
