package server

import (
	"context"
	"fmt"
	"lgc/com"
	"os"
	"os/exec"
)

type Task struct {
	ID      int
	State   int
	Ip      string
	Pattern string `json:"pattern"`
	Team    string `json:"team"`
	Branch  string `json:"branch"`
	ctx     context.Context
	cancel  context.CancelFunc
}

func (t *Task) run() {
	if t.State == 1 || t.cancel == nil {
		return
	}

	t.State = 1

	// log
	logPath := com.GetLogPath(t.ID)
	log, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		StopTask()
		return
	}
	defer log.Close()

	tasks := []string{fmt.Sprintf("mo t %s --%s@%s", t.Pattern, t.Team, t.Branch)}
	tasks = append(tasks, com.GetSufTask()...)
	for _, v := range tasks {
		cmd := exec.CommandContext(t.ctx, v)
		cmd.Stderr = log
		cmd.Stdout = log
		cmd.Dir = com.WkDir

		err = cmd.Run()
		if err != nil {
			StopTask()
			return
		}

		code := 0
		if cmd.ProcessState != nil {
			code = cmd.ProcessState.ExitCode()
		}
		if code != 0 {
			StopTask()
		}
	}
}
