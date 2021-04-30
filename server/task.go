package server

import (
	"context"
	"fmt"
	"lgc/com"
	"os"
	"os/exec"
)

type Task struct {
	ID      int    `json:"id" sql:"id"`
	State   int    `json:"state" sql:"state"`
	Ip      string `json:"ip" sql:"ip"`
	Pattern string `json:"pattern" sql:"pattern"`
	Team    string `json:"team" sql:"team"`
	Branch  string `json:"branch" sql:"branch"`
	CT      string
	ET      string
	ctx     context.Context
	cancel  context.CancelFunc
	// CT      util.JsonTime `json:"create"`
	// ET      util.JsonTime `json:"time"`
}

func (t *Task) run() {
	if t.State == com.Running || t.cancel == nil {
		return
	}

	t.State = com.Running
	updateTask(t.ID, t.State)

	logPath := com.LogPath(t.ID)
	log, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
		StopTask(com.Interrupt, t.ID)
		return
	}
	defer log.Close()

	tasks := []string{fmt.Sprintf("mo t %s ---%s@%s", t.Pattern, t.Team, t.Branch)}
	tasks = append(tasks, com.SufTask()...)
	for _, v := range tasks {
		if t.State != com.Running {
			return
		}

		bash := "bash"
		par := "-c"
		if com.OS == "windows" {
			bash = "cmd"
			par = "/c"
		}

		cmd := exec.CommandContext(t.ctx, bash, par, v)
		cmd.Stderr = log
		cmd.Stdout = log
		cmd.Dir = com.WkDir()

		err = cmd.Run()
		if err != nil {
			fmt.Println("running:", err.Error())
			if t.State < com.Succ {
				StopTask(com.Interrupt, t.ID)
			}
			return
		}

		code := 0
		if cmd.ProcessState != nil {
			code = cmd.ProcessState.ExitCode()
		}
		if code != 0 {
			fmt.Println("coed error: ", code)
			StopTask(com.Interrupt, t.ID)
		}
	}

	StopTask(com.Succ, t.ID)
}

func (t *Task) end() {
	if t.cancel == nil {
		return
	}
	t.cancel()
	t.cancel = nil
}
