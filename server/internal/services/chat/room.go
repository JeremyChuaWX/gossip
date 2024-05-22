package chat

import (
	"context"
	"gossip/internal/models"
	"gossip/internal/services/message"
	"log"

	"github.com/gofrs/uuid/v5"
)

type chatRoom struct {
	service  *Service
	alive    chan bool
	ingress  chan event
	handlers map[eventName]func(*chatRoom, event)
	room     *models.Room
	userIds  map[uuid.UUID]bool
}

func newChatRoom(
	service *Service,
	room *models.Room,
	roomUsers []models.RoomUser,
) *chatRoom {
	r := &chatRoom{
		service:  service,
		alive:    make(chan bool),
		ingress:  make(chan event),
		handlers: make(map[eventName]func(*chatRoom, event)),
		room:     room,
		userIds:  make(map[uuid.UUID]bool),
	}

	for _, roomUser := range roomUsers {
		r.userIds[roomUser.UserId] = true
	}

	r.handlers[MESSAGE] = (*chatRoom).messageEventHandler
	r.handlers[USER_JOIN_ROOM] = (*chatRoom).userJoinRoomEventHandler
	r.handlers[USER_LEAVE_ROOM] = (*chatRoom).userLeaveRoomEventHandler

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
		case e := <-r.ingress:
			handler, ok := r.handlers[e.name()]
			if !ok {
				log.Println("invalid event")
				continue
			}
			handler(r, e)
		}
	}
}

func (r *chatRoom) close() {
	r.alive <- false
}

// handlers

func (r *chatRoom) messageEventHandler(e event) {
	event := e.(*messageEvent)
	ctx := context.Background()
	r.service.messageService.Save(ctx, message.SaveDto{
		UserId: event.userId,
		RoomId: event.roomId,
		Body:   event.payload.Body,
	})
	for userId := range r.userIds {
		// NOTE: skip the sender
		if userId == event.userId {
			continue
		}
		user, ok := r.service.chatUsers[userId]
		if !ok {
			log.Printf("user %s not found", userId.String())
		}
		user.ingress <- event
	}
}

func (r *chatRoom) userJoinRoomEventHandler(e event) {
	event := e.(*userJoinRoomEvent)
	if r.room.Id != event.roomId {
		log.Println("wrong room")
		return
	}
	if _, ok := r.userIds[event.userId]; ok {
		log.Println("user already in room")
		return
	}
	r.userIds[event.userId] = true
}

func (r *chatRoom) userLeaveRoomEventHandler(e event) {
	event := e.(*userLeaveRoomEvent)
	if r.room.Id != event.roomId {
		log.Println("wrong room")
		return
	}
	if _, ok := r.userIds[event.userId]; !ok {
		log.Println("user not in room")
		return
	}
	delete(r.userIds, event.userId)
}
