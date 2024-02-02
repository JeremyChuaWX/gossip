package handlers

import (
	"encoding/json"
	"fmt"
	"gossip/internal/manager"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Handlers struct {
	Manager    *manager.Manager
	WsUpgrader *websocket.Upgrader
}

func (h *Handlers) NewRoomHandler(w http.ResponseWriter, r *http.Request) {
	// parse URL query for room name
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		io.WriteString(w, "error parsing URL query")
		return
	}
	room, err := h.Manager.AddRoom(name)
	if err != nil {
		log.Println(err)
		return
	}
	go room.BroadcastMessages() // TODO: add context to cancel when room closes
	io.WriteString(w, fmt.Sprintf("room %s created", name))
}

func (h *Handlers) GetRoomsHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Rooms []string `json:"rooms"`
	}
	rooms := h.Manager.Rooms()
	body := Response{Rooms: rooms}
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func (h *Handlers) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade connection to websocket
	conn, err := h.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading connection", err)
		return
	}
	defer conn.Close()

	// parse URL param for room name
	roomName := chi.URLParam(r, "room")
	if len(roomName) == 0 {
		log.Println("error parsing URL param")
		return
	}

	// parse URL query for client name
	clientName := r.URL.Query().Get("username")
	if len(clientName) == 0 {
		log.Println("error parsing URL query")
		return
	}

	// get room by name
	room, err := h.Manager.Room(roomName)
	if err != nil {
		log.Println(err)
		return
	}

	// client join room
	client, err := room.AddClient(clientName, conn)
	if err != nil {
		log.Println(err)
		return
	}

	go client.ReadMessages()
}

func parseJSON[T any](r *http.Request) (T, error) {
	decoder := json.NewDecoder(r.Body)
	var obj T
	err := decoder.Decode(&obj)
	return obj, err
}
