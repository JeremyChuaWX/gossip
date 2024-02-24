package chat

import (
	"log"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
)

type client struct {
	userId   uuid.UUID
	username string
	rooms    map[*room]bool
	conn     *websocket.Conn
	ingress  chan event
	service  *service
	handlers map[eventType]func(*client, event)
}

func newClient(
	userId uuid.UUID,
	username string,
	conn *websocket.Conn,
	service *service,
) *client {
	c := &client{
		userId:   userId,
		username: username,
		rooms:    make(map[*room]bool),
		conn:     conn,
		ingress:  make(chan event),
		service:  service,
		handlers: make(map[eventType]func(*client, event)),
	}
	c.handlers[MESSAGE] = (*client).messageHandler
	c.handlers[CLIENT_JOIN_ROOM] = (*client).clientJoinRoomHandler
	c.handlers[CLIENT_LEAVE_ROOM] = (*client).clientLeaveRoomHandler
	return c
}

func (c *client) init() {
	go c.receiveEvents()
}

// func (c *client) disconnect() {
// 	c.conn.Close()
// 	close(c.ingress)
// 	c.service.ingress <- makeClientDisconnectEvent(c)
// }

// run as goroutine
func (c *client) receiveEvents() {
	for {
		e, ok := <-c.ingress
		if !ok {
			continue
		}
		handler, ok := c.handlers[e.name()]
		if !ok {
			log.Println("invalid event")
			continue
		}
		handler(c, e)
	}
}

// handlers

func (c *client) messageHandler(e event) {
	event := e.(*messageEvent)
	msg := event.toJSON()
	c.conn.WriteJSON(msg) // TODO: handle error
}

func (c *client) clientJoinRoomHandler(e event) {
	event := e.(*clientJoinRoomEvent)
	c.rooms[event.room] = true
}

func (c *client) clientLeaveRoomHandler(e event) {
	event := e.(*clientLeaveRoomEvent)
	delete(c.rooms, event.room)
}
