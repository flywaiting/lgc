package com

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"lgc/util"
	"os"
	"path"
)

type cfg struct {
	Log      string   `json:"logDir"`
	Db       string   `json:"db"`
	WkDir    string   `json:"workDir"`
	Root     string   `json:"root"`
	SufTasks []string `json:"sufTasks"`
}

type cmd struct {
	branchs [][]string
	teams   []string
}

// 配置相关信息
var (
	cfgInfo *cfg
	cmdInfo *cmd
)

//go:embed cfg/cfg.json
var cfgBytes []byte

//go:embed cfg/sql.sql
var SqlStr string

func init() {
	cfgInfo = &cfg{}
	err := json.Unmarshal(cfgBytes, cfgInfo)
	util.ErrCheck(err)

	loadConf()

	cfgInfo.Log = path.Join(cfgInfo.Root, cfgInfo.Log)
	cfgInfo.Db = path.Join(cfgInfo.Root, cfgInfo.Db)

	for _, dir := range []string{cfgInfo.Log, cfgInfo.Db} {
		if _, err = os.Stat(dir); os.IsNotExist(err) {
			// os.Mkdir(logDir, os.ModeDir)
			err = os.MkdirAll(dir, 0755)
			util.ErrCheck(err)
		}
	}
}

func loadConf() {
	fn := "cfg.json"
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		data, err := os.ReadFile(fn)
		util.ErrCheck(err)
		err = json.Unmarshal(data, cfgInfo)
		util.ErrCheck(err)
	}
}

func LogPath(id int) string {
	return path.Join(cfgInfo.Log, fmt.Sprint(id, ".log"))
}

func DbPath() string {
	return path.Join(cfgInfo.Db, "db")
}

func SufTask() []string {
	return cfgInfo.SufTasks
}

func WkDir() string {
	return cfgInfo.WkDir
}

// 命令相关信息
func CmdInfo(refresh bool) (data []byte, err error) {
	if refresh || cmdInfo == nil {
		fn := ".mo.env.js"
		bytes, err := os.ReadFile(path.Join(cfgInfo.WkDir, fn))
		if err != nil {
			return nil, err
		}

		cmdInfo = &cmd{
			branchs: util.CatchBranchs(bytes),
			teams:   util.CatchTeams(bytes),
		}
	}

	data, err = json.Marshal(cmdInfo)
	return
}
