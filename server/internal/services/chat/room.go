package chat

import (
	"gossip/internal/models"
	"log"

	"github.com/gofrs/uuid/v5"
)

type chatRoom struct {
	service  *service
	alive    chan bool
	ingress  chan event
	handlers map[eventName]func(*chatRoom, event)
	room     *models.Room
	userIds  map[uuid.UUID]bool
}

func newChatRoom(
	service *service,
	room *models.Room,
	userIds []uuid.UUID,
) *chatRoom {
	r := &chatRoom{
		service:  service,
		alive:    make(chan bool),
		ingress:  make(chan event),
		handlers: make(map[eventName]func(*chatRoom, event)),
		room:     room,
		userIds:  make(map[uuid.UUID]bool),
	}

	for _, userId := range userIds {
		r.userIds[userId] = true
	}

	r.handlers[MESSAGE] = (*chatRoom).messageEventHandler

	go r.receiveEvents()
	return r
}

func (r *chatRoom) receiveEvents() {
	defer (func() {
		close(r.ingress)
		close(r.alive)
	})()
	for {
		select {
		case <-r.alive:
			return
			// case e := <-r.ingress:
		}
	}
}

func (r *chatRoom) close() {
	r.alive <- false
}

// handlers

func (r *chatRoom) messageEventHandler(e event) {
	for userId := range r.userIds {
		user, ok := r.service.chatUsers[userId]
		if !ok {
			log.Printf("user %s not found", userId.String())
		}
		user.ingress <- e
	}
}
