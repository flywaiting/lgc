package util

import "github.com/BurntSushi/toml"

type Config struct {
	Git       Git       `toml:"git"`
	Project   Project   `toml:"project"`
	Workspace Workspace `toml:"workspace"`
}

type Git struct {
	TmpRepository string     `toml:"tmpRepository"`
	Apply         [][]string `toml:"apply"`
}
type Project struct {
	Teams   []string `toml:"teams"`
	Pattern []string `toml:"pattern"`
	Root    string   `toml:"root"`   // 项目路径
	Config  string   `toml:"config"` // 项目配置文件
	Flag    string   `toml:"flag"`   // 识别标识
}
type Workspace struct {
	Root string `toml:"root"` // 程序运行位置
	Log  string `toml:"log"`  // 日志目录
}

func InitConfig() *Config {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		panic("配置文件解析错误")
	}

	return &config
}
