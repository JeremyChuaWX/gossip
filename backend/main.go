package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const ADDRESS string = "127.0.0.1:3000"

func main() {
	wsUpgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	handlers := &Handlers{
		Rooms:      make(map[string]*Room),
		WsUpgrader: wsUpgrader,
	}
	router := NewRouter(handlers)

	log.Println("running server on address", ADDRESS)
	log.Fatal(http.ListenAndServe(ADDRESS, router.Mux))
}
