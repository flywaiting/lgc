package server

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	ActiveTime int64    `json:"activeTime"` // 开始执行时间
	EndTime    int64    `json:"finishTime"`

	alias string // 别名

	ctx context.Context
	fn  context.CancelFunc
}

func (t *TaskItem) run() {
	name := filepath.Join(config.Workspace.Root, config.Workspace.Log, fmt.Sprint(t.Id, ".log"))
	log, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		println(err.Error())
		taskHub.finish(Error)
		return
	}
	defer log.Close()

	tasks := [][]string{strings.Split(fmt.Sprintf("mo t %s ---%s@%s", t.Pattern, t.Team, t.alias), " ")}
	tasks = append(tasks, config.Git.Apply...)
	for idx, task := range tasks {
		if idx > 0 {
			t.Status = Git
		}
		cmd := exec.CommandContext(t.ctx, task[0], task[1:]...)
		cmd.Stdout = log
		cmd.Stderr = log
		cmd.Dir = config.Project.Root

		err = cmd.Run()
		if err != nil {
			println("运行中报错了: ", err.Error())
			taskHub.finish(Error)
			return
		}

		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() != 0 {
			println("exit code: ", cmd.ProcessState.ExitCode())
			taskHub.finish(Error)
			return
		}
	}
	taskHub.finish(Finish)
}

// 任务结束
func (t *TaskItem) finish(status int) {
	t.Status = status
	t.EndTime = time.Now().Unix()
	if t.fn == nil {
		return
	}

	t.fn()
	t.fn = nil
	t.ctx = nil
}
