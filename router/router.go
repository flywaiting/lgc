package router

import (
	"fmt"
	"lgc/com"
	"lgc/server"
	"lgc/util"
	"net/http"
	"strings"
)

var (
	Mux *http.ServeMux
)

func InitRouter() {
	Mux = http.NewServeMux()

	// fsys, err := fs.Sub(static, "static")
	// util.ErrCheck(err)
	// Mux.Handle("/lgc/", http.StripPrefix("/lgc/", http.FileServer(http.FS(fsys))))
	Mux.Handle("/lgc/", http.StripPrefix("/lgc/", http.FileServer(http.Dir("static"))))
	Mux.HandleFunc("/addTask", addTask)
	Mux.HandleFunc("/taskInfo", taskInfo)
	Mux.HandleFunc("/cmdInfo", cmdInfo)
	Mux.HandleFunc("/loop", loop)
	Mux.HandleFunc("/stopTask", stopTask)
	Mux.HandleFunc("/removeTask", removeTask)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if len(ip) < 8 {
		http.Error(w, "IP获取失败", http.StatusBadRequest)
		return
	}

	err := server.AddTask(r.Body, ip)
	if err != nil {
		fmt.Println(err.Error())
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

func stopTask(w http.ResponseWriter, r *http.Request) {
	id := util.Atoi(r.Body)
	if id == 0 {
		http.Error(w, "ID获取失败", http.StatusBadRequest)
		return
	}
	server.StopTask(com.Kill, id)
}

func removeTask(w http.ResponseWriter, r *http.Request) {
	id := util.Atoi(r.Body)
	if id == 0 {
		http.Error(w, "ID获取失败", http.StatusBadRequest)
		return
	}
	server.StopTask(com.Remove, id)
}

func loop(w http.ResponseWriter, r *http.Request) {
	id := util.Atoi(r.Body)
	ing := server.IsRunning(id)

	code := http.StatusOK
	if ing {
		code = http.StatusAccepted
	}
	w.WriteHeader(code)
}
