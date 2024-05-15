package chat

import (
	"gossip/internal/models"
	"log"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
)

type chatUser struct {
	service  *service
	alive    chan bool
	ingress  chan event
	handlers map[eventName]func(*chatUser, event)
	user     *models.User
	roomIds  map[uuid.UUID]bool
	conn     *websocket.Conn
}

func newChatUser(
	service *service,
	user *models.User,
	roomIds []uuid.UUID,
	conn *websocket.Conn,
) *chatUser {
	u := &chatUser{
		service:  service,
		alive:    make(chan bool),
		ingress:  make(chan event),
		handlers: make(map[eventName]func(*chatUser, event)),
		user:     user,
		roomIds:  make(map[uuid.UUID]bool),
		conn:     conn,
	}

	for _, roomId := range roomIds {
		u.roomIds[roomId] = true
	}

	u.handlers[MESSAGE] = (*chatUser).messageEventHandler

	go u.receiveEvents()
	go u.receiveWebsocket()

	return u
}

func (u *chatUser) receiveEvents() {
	defer (func() {
		close(u.ingress)
		close(u.alive)
	})()
	for {
		select {
		case <-u.alive:
			return
		case e := <-u.ingress:
			handler, ok := u.handlers[e.name()]
			if !ok {
				log.Printf("invalid event")
				continue
			}
			handler(u, e)
		}
	}
}

func (u *chatUser) receiveWebsocket() {
	for {
		select {
		case <-u.alive:
			return
		default:
			var payload payload
			if err := u.conn.ReadJSON(&payload); err != nil {
				u.disconnect()
				return
			}
			roomId, err := uuid.FromString(payload.RoomId)
			if err != nil {
				log.Printf("invalid room ID")
				continue
			}
			room, ok := u.service.chatRooms[roomId]
			if !ok {
				log.Printf("invalid room ID")
				continue
			}
			room.ingress <- newMessageEvent(payload)
		}
	}
}

func (u *chatUser) disconnect() error {
	u.alive <- false
	return u.conn.Close()
}

// handlers

func (u *chatUser) messageEventHandler(e event) {
	event := e.(*messageEvent)
	if err := u.conn.WriteJSON(event.payload); err != nil {
		u.disconnect()
		return
	}
}
