package chat

import (
	"context"
	"gossip/internal/models"
	userPackage "gossip/internal/services/user"
	"time"

	"github.com/gofrs/uuid/v5"
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
	userId uuid.UUID,
) (*user, error) {
	user := &user{
		service:  service,
		conn:     conn,
		ingress:  make(chan event),
		alive:    make(chan bool),
		send:     make(chan *message),
		handlers: make(map[eventName]func(*user, event)),
	}

	userModel, err := service.userService.FindOne(
		context.Background(),
		userPackage.FindOneDTO{UserId: userId},
	)
	if err != nil {
		return nil, err
	}
	user.model = &userModel

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
			user.alive <- false
			return
		}
		messageEvent, err := newMessageEvent(&message)
		if err != nil {
			// TODO: handle invalid message
			continue
		}
		room, ok := user.service.rooms[messageEvent.roomId]
		if !ok {
			// TODO: handle room not found
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
			}
			if err := user.conn.WriteJSON(message); err != nil {
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
