package chat

import (
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
}

func newClient(
	userId uuid.UUID,
	username string,
	conn *websocket.Conn,
	service *service,
) *client {
	return &client{
		userId:   userId,
		username: username,
		rooms:    make(map[*room]bool),
		conn:     conn,
		ingress:  make(chan event),
		service:  service,
	}
}

func (c *client) init() {
	go c.receiveEvents()
}

func (c *client) disconnect() {
	c.conn.Close()
	close(c.ingress)
	c.service.ingress <- makeClientDisconnectEvent(c)
}

// run as goroutine
func (c *client) receiveEvents() {
	for {
		e := <-c.ingress
		switch e := e.(type) {
		case *messageEvent:
			msg := e.toJSON()
			c.conn.WriteJSON(msg) // TODO: handle error
		case *clientJoinRoomEvent:
			c.rooms[e.room] = true
		case *clientLeaveRoomEvent:
			delete(c.rooms, e.room)
		}
	}
}
