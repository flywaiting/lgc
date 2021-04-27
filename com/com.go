package com

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

var info *cfg

//go:embed cfg/cfg.json
var cfgCtx []byte

func init() {
	info = &cfg{}
	err := json.Unmarshal(cfgCtx, info)
	util.ErrCheck(err)

	loadConf()

	info.Log = path.Join(info.Root, info.Log)
	info.Db = path.Join(info.Root, info.Db)

	for _, dir := range []string{info.Log, info.Db} {
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
		data, err := ioutil.ReadFile(fn)
		util.ErrCheck(err)
		err = json.Unmarshal(data, info)
		util.ErrCheck(err)
	}
}

func LogPath(id int) string {
	return path.Join(info.Log, fmt.Sprint(id, ".log"))
}

func DbPath() string {
	return path.Join(info.Db, "db")
}

func SufTask() []string {
	return info.SufTasks
}
