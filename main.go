package main

import (
	"lgc/api"
	"net/http"
)

func main() {
	demo()
}

func demo() {
	//http.FileServer("")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", api.Home)
	http.HandleFunc("/ws", api.Up2Ws)

	http.ListenAndServe(":8080", nil)
}
