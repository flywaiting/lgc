package com

import (
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
	Pattern  []string `json:"pattern"`
	SufTasks []string `json:"sufTasks"`
}

type cmd struct {
	Teams   []string   `json:"teams"`
	Pattern []string   `json:"patterns"`
	Branchs [][]string `json:"branchs"`
}

// 配置相关信息
var (
	cfgInfo *cfg
	cmdInfo *cmd
)

func InitCom(data *[]byte) {
	cfgInfo = &cfg{}
	err := json.Unmarshal(*data, cfgInfo)
	util.ErrCheck(err)

	loadConf()

	cfgInfo.Log = path.Join(cfgInfo.Root, cfgInfo.Log)
	cfgInfo.Db = path.Join(cfgInfo.Root, cfgInfo.Db)

	for _, dir := range []string{cfgInfo.Log, cfgInfo.Db} {
		os.RemoveAll(dir)
		err = os.MkdirAll(dir, 0755)
		util.ErrCheck(err)
		// if _, err = os.Stat(dir); os.IsNotExist(err) {
		// 	// os.Mkdir(logDir, os.ModeDir)
		// }
	}
}

func loadConf() {
	fn := path.Join(cfgInfo.WkDir, "cfg.json")
	if _, err := os.Stat(fn); !os.IsNotExist(err) {
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
			Branchs: util.CatchBranchs(bytes),
			Teams:   util.CatchTeams(bytes),
			Pattern: cfgInfo.Pattern,
		}
	}

	data, err = json.Marshal(cmdInfo)
	return
}
