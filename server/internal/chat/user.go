package chat

import (
	"log/slog"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
)

type user struct {
	userId   uuid.UUID
	username string
	service  *Service
	ingress  chan event
	alive    chan bool

	conn *websocket.Conn
	send chan *message
}

func newUser(
	service *Service,
	conn *websocket.Conn,
	userId uuid.UUID,
	username string,
) (*user, error) {
	user := &user{
		userId:   userId,
		username: username,
		service:  service,
		ingress:  make(chan event),
		alive:    make(chan bool),

		conn: conn,
		send: make(chan *message),
	}
	user.conn.SetReadLimit(MAX_MESSAGE_SIZE)
	go user.receiveEvents()
	go user.readPump()
	go user.writePump()
	return user, nil
}

func (user *user) readPump() {
	for {
		var message message
		if err := user.conn.ReadJSON(&message); err != nil {
			slog.Error(
				"error reading JSON",
				"error",
				err.Error(),
				"message",
				message,
			)
			user.alive <- false
			return
		}
		slog.Info("readPump message", "message", message)
		messageEvent, err := newMessageEvent(
			user.userId,
			user.username,
			&message,
		)
		if err != nil {
			slog.Error("error creating message event", "message", message)
			continue
		}
		room, ok := user.service.rooms[messageEvent.roomId]
		if !ok {
			slog.Error("room not found", "message", message)
			continue
		}
		room.ingress <- messageEvent
	}
}

func (user *user) writePump() {
	defer user.conn.Close()
	for {
		message, ok := <-user.send
		if !ok {
			slog.Error("user send channel closed")
			return
		}
		slog.Info("writePump message", "message", message)
		if err := user.conn.WriteJSON(message); err != nil {
			slog.Error(
				"error writing JSON",
				"error",
				err.Error(),
				"message",
				message,
			)
			return
		}
	}
}

// actor methods

func (user *user) receiveEvents() {
	for {
		select {
		case <-user.alive:
			user.disconnect()
			return
		case event, ok := <-user.ingress:
			if !ok {
				return
			}
			user.eventHandler(event)
		}
	}
}

func (user *user) disconnect() {
	user.service.ingress <- userDisconnectedEvent{userId: user.userId}
	user.conn.Close()
	close(user.ingress)
	close(user.alive)
	close(user.send)
}

// event management

func (user *user) eventHandler(event event) {
	switch event := event.(type) {
	default:
		slog.Error("invalid event", "event", event)
	}
}
