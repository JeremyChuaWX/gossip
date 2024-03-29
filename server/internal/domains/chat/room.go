package chat

import "log"

type room struct {
	name     string
	clients  map[*client]bool
	ingress  chan event
	alive    chan bool
	service  *Service
	handlers map[eventType]func(*room, event)
}

func newRoom(name string, service *Service) *room {
	r := &room{
		name:     name,
		clients:  make(map[*client]bool),
		ingress:  make(chan event, 100),
		alive:    make(chan bool),
		service:  service,
		handlers: make(map[eventType]func(*room, event)),
	}
	r.handlers[MESSAGE] = (*room).messageHandler
	r.handlers[CLIENT_JOIN_ROOM] = (*room).clientJoinRoomHandler
	r.handlers[CLIENT_LEAVE_ROOM] = (*room).clientLeaveRoomHandler
	r.handlers[DESTROY_ROOM] = (*room).destroyRoomHandler
	return r
}

func (r *room) init() {
	go r.receiveEvents()
}

// run as goroutine
func (r *room) receiveEvents() {
	defer (func() {
		close(r.ingress)
		close(r.alive)
	})()
	for {
		select {
		case <-r.alive:
			return
		case e := <-r.ingress:
			log.Printf("[room] event received: %v", e)
			handler, ok := r.handlers[e.name()]
			if !ok {
				log.Println("invalid event")
				continue
			}
			handler(r, e)
		}
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

func (r *room) destroyRoomHandler(e event) {
	for client := range r.clients {
		client.ingress <- makeClientLeaveRoomEvent(client, r)
	}
	r.alive <- false
}
