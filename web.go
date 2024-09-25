package main

import (
	"fmt"
	"lgc/router"
	"net/http"
)

// var (
// 	//go:embed cfg/cfg.json
// 	cfg []byte
// 	//go:embed cfg/sql.sql
// 	sql string
// 	//go:embed static
// 	static embed.FS
// )

func init() {
	router.InitRouter()
	// com.InitCom()
	// server.InitDB()
}

func main() {
	err := http.ListenAndServe(":6464", router.Mux)
	if err != nil {
		fmt.Println("服务启动失败:", err)
	}
}
