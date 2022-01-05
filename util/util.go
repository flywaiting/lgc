package util

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

type BranchSection struct {
	Dir       string
	UpdateCmd string
	GetCmd    string
}

var SecBranch BranchSection

// Tsp 获取当前时间戳[sec]
func Tsp() int64 {
	return time.Now().Unix()
}

// Branches 获取指定 repo 的所有分支
func Branches(dir string) ([][]byte, error) {
	// 更新仓库
	upCmd := exec.Command("bash", "-c", SecBranch.UpdateCmd)
	// upCmd := exec.Command("bash", "-c", "git pull")
	upCmd.Dir = SecBranch.Dir
	err := upCmd.Run()
	if err != nil {
		return nil, err
	}
	// 获取分支
	cmd := exec.Command("bash", "-c", SecBranch.GetCmd)
	// cmd := exec.Command("bash", "-c", "git branch -r")
	cmd.Dir = SecBranch.Dir
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var branches [][]byte
	res := bytes.FieldsFunc(data, func(r rune) bool { return r == '\n' })
	head := []byte("HEAD")
	origin := []byte("origin/")
	for _, b := range res {
		if bytes.Contains(b, head) {
			continue
		}
		branches = append(branches, bytes.TrimPrefix(bytes.TrimSpace(b), origin))
	}

	return branches, nil
}

// InferDir 向上推断目标文件、文件夹所在路径
func InferDir(d string) string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for pwd != "/" {
		if _, err := os.Stat(path.Join(pwd, d)); err == nil || os.IsExist(err) {
			return pwd
		} else {
			pwd = filepath.Dir(pwd)
		}
	}
	return ""
}
