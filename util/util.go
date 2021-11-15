package util

import (
	"os/exec"
	"strings"
	"time"
)

// 获取当前时间戳[sec]
func Tsp() int64 {
	return time.Now().Unix()
}

// 获取指定 repo 的所有分支
func Branchs(dir string) ([]string, error) {
	// 更新仓库
	upCmd := exec.Command("bash", "-c", "git pull")
	upCmd.Dir = dir
	err := upCmd.Run()
	if err != nil {
		return nil, err
	}
	// 获取分支
	cmd := exec.Command("bash", "-c", "git branch -r")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	arr := []string{}
	for _, s := range strings.Split(string(out), "\n") {
		if len(s) > 0 && !strings.Contains(s, "HEAD") {
			arr = append(arr, strings.TrimPrefix(strings.TrimSpace(s), "origin/"))
		}
	}
	return arr, nil
}
