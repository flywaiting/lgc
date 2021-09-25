package cfg

import (
	"encoding/json"
	"os"
	"sync"
)

type mcfg struct {
}

var cfg mcfg
var loadCfg sync.Once

func Cfg() *mcfg {
	loadCfg.Do(func() {
		f, err := os.OpenFile("cfg.json", os.O_RDONLY, 0664)
		if err != nil {
			return
		}
		defer f.Close()

		decoder := json.NewDecoder(f)
		decoder.Decode(&cfg)
	})

	return &cfg
}
