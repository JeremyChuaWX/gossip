package chat

type room struct {
	name    string
	clients map[*client]bool
	ingress chan event
}

func newRoom(name string) *room {
	return &room{
		name:    name,
		clients: make(map[*client]bool),
		ingress: make(chan event),
	}
}

func (r *room) init() {
	go r.receiveEvents()
	go r.receiveMessages()
}

func (r *room) destroy() {
	// send event to chat service to handle destory
}

// run as goroutine
func (r *room) receiveEvents() {
	for {
		e := <-r.ingress
		switch e := e.(type) {
		case *messageEvent:
			r.broadcastMessages(e)
		case *clientJoinRoomEvent:
			r.clients[e.client] = true
		case *clientLeaveRoomEvent:
			delete(r.clients, e.client)
		}
	}
}

func (r *room) broadcastMessages(e *messageEvent) {
	for client := range r.clients {
		client.ingress <- e
	}
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
