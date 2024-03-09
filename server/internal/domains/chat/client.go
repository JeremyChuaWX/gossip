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
	alive    chan bool
	service  *Service
	handlers map[eventType]func(*client, event)
}

func newClient(
	userId uuid.UUID,
	username string,
	conn *websocket.Conn,
	service *Service,
) *client {
	c := &client{
		userId:   userId,
		username: username,
		rooms:    make(map[*room]bool),
		conn:     conn,
		ingress:  make(chan event, 100),
		alive:    make(chan bool),
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

func (c *client) disconnect() {
	c.conn.Close()
	for room := range c.rooms {
		room.ingress <- makeClientLeaveRoomEvent(c, room)
	}
	c.service.ingress <- makeClientDisconnectEvent(c)
	c.alive <- false
}

// run as goroutine
func (c *client) receiveEvents() {
	defer (func() {
		close(c.ingress)
		close(c.alive)
	})()
	for {
		select {
		case <-c.alive:
			return
		case e := <-c.ingress:
			log.Printf("[client] event received: %v", e)
			handler, ok := c.handlers[e.name()]
			if !ok {
				log.Println("invalid event")
				continue
			}
			handler(c, e)
		default:
			var msgJSON messageJSON
			if err := c.conn.ReadJSON(&msgJSON); err != nil {
				c.disconnect()
			}
			msgEvent := msgJSON.toEvent(c.userId)
			log.Printf("[room] websocket message received: %v", msgEvent)
			c.service.ingress <- msgEvent
		}
	}
}

// handlers

func (c *client) messageHandler(e event) {
	event := e.(*messageEvent)
	msg := event.toJSON()
	log.Printf("[client] writing websocket message: %v", msg)
	if err := c.conn.WriteJSON(msg); err != nil {
		c.disconnect()
	}
}

func (c *client) clientJoinRoomHandler(e event) {
	event := e.(*clientJoinRoomEvent)
	c.rooms[event.room] = true
}

func (c *client) clientLeaveRoomHandler(e event) {
	event := e.(*clientLeaveRoomEvent)
	delete(c.rooms, event.room)
}
