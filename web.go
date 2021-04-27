package main

import (
	"fmt"
	"lgc/router"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":6464", router.Mux)
	if err != nil {
		fmt.Println("服务启动失败:", err)
	}
}
