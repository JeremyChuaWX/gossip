package chat

import "log"

type room struct {
	name     string
	clients  map[*client]bool
	ingress  chan event
	service  *service
	handlers map[eventType]func(event)
}

func newRoom(name string, service *service) *room {
	r := &room{
		name:     name,
		clients:  make(map[*client]bool),
		ingress:  make(chan event),
		service:  service,
		handlers: make(map[eventType]func(event)),
	}
	r.handlers[MESSAGE] = r.messageHandler
	r.handlers[CLIENT_JOIN_ROOM] = r.clientJoinRoomHandler
	r.handlers[CLIENT_LEAVE_ROOM] = r.clientLeaveRoomHandler
	return r
}

func (r *room) init() {
	go r.receiveEvents()
	go r.receiveMessages()
}

func (r *room) destroy() {
	close(r.ingress)
	r.service.ingress <- makeRemoveRoomEvent(r)
}

// run as goroutine
func (r *room) receiveMessages() {
	for {
		var msg messageJSON
		for client := range r.clients {
			client.conn.ReadJSON(&msg) // TODO: handle error
			e := msg.toEvent()
			r.ingress <- e
		}
	}
}

// run as goroutine
func (r *room) receiveEvents() {
	for {
		e := <-r.ingress
		handler, ok := r.handlers[e.name()]
		if !ok {
			log.Panic("no handler")
		}
		handler(e)
	}
}

// handlers

func (r *room) messageHandler(e event) {
	for client := range r.clients {
		client.ingress <- e
	}
}

func (r *room) clientJoinRoomHandler(e event) {
	event := e.(*clientJoinRoomEvent)
	r.clients[event.client] = true
}

func (r *room) clientLeaveRoomHandler(e event) {
	event := e.(*clientLeaveRoomEvent)
	delete(r.clients, event.client)
}
