package router

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	Mux *http.ServeMux
)

func init() {
	Mux = http.NewServeMux()

	// 静态资源
	file := http.FileServer(http.Dir("static"))
	Mux.Handle("/lgc/", http.StripPrefix("/lgc/", file))

	Mux.HandleFunc("/addTask", addTask)
}

func addTask(w http.ResponseWriter, r *http.Request) {

	ip := strings.Split(r.RemoteAddr, ":")[0]
	if len(ip) < 8 {
		http.Error(w, "IP错误", http.StatusBadRequest)
		return
	}
	fmt.Println(ip)
}
