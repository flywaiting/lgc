package router

import (
	"embed"
	"lgc/com"
	"lgc/server"
	"net/http"
	"strings"
)

var (
	Mux *http.ServeMux
)

func InitRouter(static *embed.FS) {
	Mux = http.NewServeMux()

	Mux.Handle("/lgc/", http.StripPrefix("/lgc/", http.FileServer(http.FS(static))))
	Mux.HandleFunc("/addTask", addTask)
	Mux.HandleFunc("/taskInfo", taskInfo)
	Mux.HandleFunc("/cmdInfo", cmdInfo)
	Mux.HandleFunc("/loop", loop)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if len(ip) < 8 {
		http.Error(w, "IP获取失败", http.StatusBadRequest)
		return
	}

	err := server.AddTask(r.Body, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	server.Run()
}

func taskInfo(w http.ResponseWriter, r *http.Request) {
	data, err := server.TaskInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Write(data)
}

func cmdInfo(w http.ResponseWriter, r *http.Request) {
	data, err := com.CmdInfo(r.ContentLength > 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	w.Write(data)
}

func loop(w http.ResponseWriter, r *http.Request) {

}
