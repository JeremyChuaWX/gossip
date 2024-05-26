package chat

import (
	"gossip/internal/models"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

type user struct {
	model    *models.User
	service  *Service
	ingress  chan event
	alive    chan bool
	handlers map[eventName]func(*user, event)

	conn *websocket.Conn
	send chan *message
}

func newUser(
	service *Service,
	conn *websocket.Conn,
	userModel *models.User,
) (*user, error) {
	user := &user{
		model:    userModel,
		service:  service,
		ingress:  make(chan event),
		alive:    make(chan bool),
		handlers: make(map[eventName]func(*user, event)),

		conn: conn,
		send: make(chan *message),
	}

	user.registerEventHandlers()

	go user.receiveEvents()
	go user.readPump()
	go user.writePump()

	return user, nil
}

func (user *user) readPump() {
	user.conn.SetReadLimit(MAX_MESSAGE_SIZE)
	user.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
	user.conn.SetPongHandler(func(string) error {
		user.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
		return nil
	})
	for {
		var message message
		if err := user.conn.ReadJSON(&message); err != nil {
			slog.Error("error reading JSON", "message", message)
			user.alive <- false
			return
		}
		messageEvent, err := newMessageEvent(&message)
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
	ticker := time.NewTicker(PING_PERIOD)
	defer ticker.Stop()
	for {
		select {
		case message, ok := <-user.send:
			if !ok {
				// TODO: handle user send channel closed
				slog.Error("user send channel closed")
			}
			if err := user.conn.WriteJSON(message); err != nil {
				slog.Error("error writing JSON", "message", message)
				user.alive <- false
				return
			}
		case <-ticker.C:
			user.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
			if err := user.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// actor methods

func (user *user) receiveEvents() {
	defer user.disconnect()
	for {
		select {
		case <-user.alive:
			return
		case event, ok := <-user.ingress:
			if !ok {
				return
			}
			handler, ok := user.handlers[event.name()]
			if !ok {
				continue
			}
			handler(user, event)
		}
	}
}

func (user *user) disconnect() {
	user.conn.Close()
	user.service.ingress <- &userDisconnectEvent{userId: user.model.Id}
}

// event management

func (u *user) registerEventHandlers() {
}
