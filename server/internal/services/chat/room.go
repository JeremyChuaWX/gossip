package chat

import (
	"context"
	"gossip/internal/models"
	roomPackage "gossip/internal/services/room"
	roomuserPackage "gossip/internal/services/roomuser"

	"github.com/gofrs/uuid/v5"
)

type room struct {
	model    *models.Room
	service  *Service
	ingress  chan event
	alive    chan bool
	handlers map[eventName]func(*room, event)

	userIds map[uuid.UUID]bool
}

func newRoom(service *Service, roomId uuid.UUID) (*room, error) {
	room := &room{
		service:  service,
		ingress:  make(chan event),
		alive:    make(chan bool),
		handlers: make(map[eventName]func(*room, event)),
		userIds:  make(map[uuid.UUID]bool),
	}

	roomModel, err := service.roomService.FindOne(
		context.Background(),
		roomPackage.FindOneDTO{RoomId: roomId},
	)
	if err != nil {
		return nil, err
	}
	room.model = &roomModel

	roomuserModels, err := service.roomuserService.FindUserIdsByRoomId(
		context.Background(),
		roomuserPackage.FindUserIdsByRoomIdDTO{RoomId: roomId},
	)
	if err != nil {
		return nil, err
	}
	for _, roomuserModel := range roomuserModels {
		room.userIds[roomuserModel.UserId] = true
	}

	room.registerEventHandlers()

	go room.receiveEvents()

	return room, nil
}

// actor methods

func (room *room) receiveEvents() {
	defer room.disconnect()
	for {
		select {
		case <-room.alive:
			return
		case event, ok := <-room.ingress:
			if !ok {
				return
			}
			handler, ok := room.handlers[event.name()]
			if !ok {
				continue
			}
			handler(room, event)
		}
	}
}

func (room *room) disconnect() {
}

// event management

func (r *room) registerEventHandlers() {
	r.handlers[MESSAGE_EVENT] = (*room).messageEventHandler
}

func (room *room) messageEventHandler(event event) {
	messageEvent, ok := event.(*messageEvent)
	if !ok {
		return
	}
	for userId := range room.userIds {
		user, ok := room.service.users[userId]
		if !ok {
			continue
		}
		user.send <- messageEvent.payload
	}
}
