package api

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {

	//pwd, err := os.Getwd()
	//fmt.Println(pwd, err)
	http.ServeFile(w, r, "./static/index.html")
	//b, err := os.ReadFile(pwd + "/static/index.html")
	//fmt.Println(err)
	//if err == nil {
	//	w.Write(b)
	//}
	//tpl, err := template.ParseFiles(pwd + "/static/index.html")
	//if err != nil {
	//	fmt.Println("index 解析错误")
	//	return
	//}
	//tpl.Execute(w, nil)
}

var up = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Up2Ws(w http.ResponseWriter, r *http.Request) {
	conn, err := up.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go msg(conn)

}

func msg(conn *websocket.Conn) {
	//defer  func(){
	//	conn.Close()
	//}
	defer conn.Close()
	i := 1
	for {
		t, m, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(t, string(m))
		conn.WriteMessage(i, []byte("hi, ws..."))
		//i++
	}
}

//func FS(w http)
