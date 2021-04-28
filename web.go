package main

import (
	"embed"
	"fmt"
	"lgc/router"
	"net/http"
)

//go:embed static
var static embed.FS

func init() {
	router.InitRouter(&static)
}

func main() {
	err := http.ListenAndServe(":6464", router.Mux)
	if err != nil {
		fmt.Println("服务启动失败:", err)
	}
}
