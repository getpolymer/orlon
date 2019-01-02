package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	Subscribe(NewSocketWriter(ws))
}

func serveWebSocket() {
	http.HandleFunc("/ws", serveWs)
}

func serveWebSite() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./web/index.html") })
	http.Handle("/assets/", http.FileServer(http.Dir("./web/assets")))
}

func setupLogger() {
	f, err := os.Create("./orlon.log")
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Println("Error setting up logger ", err.Error())
		return
	}
	log.SetOutput(f)
}

func main() {
	setupLogger()
	serveWebSite()
	serveWebSocket()
	go runPseudoTerminal()
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}
