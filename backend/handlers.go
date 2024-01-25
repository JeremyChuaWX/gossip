package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Handlers struct {
	Rooms      map[string]*Room
	WsUpgrader *websocket.Upgrader
}

type WebsocketDTO struct {
	Message string `json:"message"`
}

func (h *Handlers) NewRoomHandler(w http.ResponseWriter, r *http.Request) {
	// parse URL query for room name
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		io.WriteString(w, "error parsing URL query")
		return
	}

	room := NewRoom(name)
	h.Rooms[name] = room
	go BroadcastMessages(room) // TODO: add context to cancel when room closes

	io.WriteString(w, fmt.Sprintf("room %s created", name))
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
	room, ok := h.Rooms[roomName]
	if !ok {
		log.Printf("room %s does not exist", roomName)
		return
	}

	// client join room
	_, ok = room.Clients[clientName]
	if ok {
		log.Printf("client %s already joined room", clientName)
		return
	}
	client := &Client{
		Name: clientName,
		Room: room.Name,
		Conn: conn,
	}
	room.Clients[client.Name] = client

	// get requests from connection and push on channel
	for {
		var dto WebsocketDTO
		if err := conn.ReadJSON(&dto); err != nil {
			log.Println(err)
			delete(room.Clients, client.Name)
			return
		}
		rm := RoomMessage{
			Timestamp: time.Now(),
			Client:    client,
			Body:      dto.Message,
		}
		room.Messages <- rm
	}
}

func BroadcastMessages(room *Room) {
	for {
		rm := <-room.Messages
		formattedMessage := rm.String()
		dto := WebsocketDTO{
			Message: formattedMessage,
		}
		for client := range room.Clients {
			conn := room.Clients[client].Conn
			if err := conn.WriteJSON(dto); err != nil {
				log.Println(err)
				conn.Close()
				delete(room.Clients, client)
			}
		}
	}
}

func parseJSON[T any](r *http.Request) (T, error) {
	decoder := json.NewDecoder(r.Body)
	var obj T
	err := decoder.Decode(&obj)
	return obj, err
}
