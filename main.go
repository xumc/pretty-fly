package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
	"strconv"
	"path"
	"os"
)

type CtlData struct {
	Number string
}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var dir string
var staticHandler http.Handler

func init() {
	dir = path.Dir(os.Args[0])
	fmt.Printf("dir: %s", dir)
    staticHandler = http.FileServer(http.Dir(dir))
}


func main() {
	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := CtlData{
			Number: r.FormValue("number"),
		}
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/sprite.jpg", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ServeHTTP(w, r)
	})
	
	conns := make([]*websocket.Conn, 0)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		conns = append(conns, conn)

		fmt.Println("conns conunt: %d", len(conns))

		for {
			msgType, numStr, err := conn.ReadMessage()
			if err != nil {
				return
			}

			currentNum, err := strconv.Atoi(string(numStr))
			if err != nil {
				fmt.Printf("error happened: %s", err)
				return
			}
			nextNum := (currentNum) % len(conns)

			if err = conns[nextNum].WriteMessage(msgType, []byte("start")); err != nil {
				return
			}
		}
	})

	http.ListenAndServe(":80", nil)

}
