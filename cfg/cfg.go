package cfg

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var cfg map[string]interface{}
var loadCfg sync.Once

func Cfg() map[string]interface{} {
	pwd, _ := os.Getwd()
	fmt.Println(pwd)
	loadCfg.Do(func() {
		f, err := os.OpenFile("./cfg.json", os.O_RDONLY, 0664)
		if err != nil {
			fmt.Println("open failed")
			return
		}
		defer f.Close()

		decoder := json.NewDecoder(f)
		decoder.Decode(&cfg)
		fmt.Println("parse")
	})

	fmt.Printf("cfg: %v\n", cfg)

	return cfg
}
