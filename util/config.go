package util

import "github.com/BurntSushi/toml"

type Config struct {
	Git     GitConfig     `toml:"git"`
	Product PorductConfig `toml:"product"`
	Server  ServerConfig  `toml:"server"`
}

type GitConfig struct {
	TmpRepository string     `toml:"tmpRepository"`
	Apply         [][]string `toml:"apply"`
}
type PorductConfig struct {
	Teams      []string `toml:"teams"`
	Pattern    []string `toml:"pattern"`
	Workspace  string   `toml:"workspace"`  // 工作目录
	ConfigFile string   `toml:"configFile"` // 配置文件
	InsertFlag string   `toml:"insertFlag"`
}
type ServerConfig struct {
	Root string `toml:"root"` // 根目录
	Log  string `toml:"log"`  // 日志目录
}

func InitConfig() *Config {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		panic("配置文件解析错误")
	}

	return &config
}
