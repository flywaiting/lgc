package util

import "github.com/BurntSushi/toml"

type Config struct {
	Git     GitConfig     `toml:"git"`
	Product PorductConfig `toml:"product"`
}

type GitConfig struct {
	TmpRepository string     `toml:"tmpRepository"`
	Apply         [][]string `toml:"apply"`
}
type PorductConfig struct {
	Teams   []string `toml:"teams"`
	Pattern []string `toml:"pattern"`
}

func InitConfig() (config Config) {
	// var
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		panic("配置文件解析错误")
	}

	return
}
