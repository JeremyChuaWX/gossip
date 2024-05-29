package chat

import (
	"context"
	"gossip/internal/models"
	roomuserPackage "gossip/internal/services/roomuser"
	"log/slog"

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

func newRoom(service *Service, roomModel *models.Room) (*room, error) {
	room := &room{
		model:    roomModel,
		service:  service,
		ingress:  make(chan event),
		alive:    make(chan bool),
		handlers: make(map[eventName]func(*room, event)),

		userIds: make(map[uuid.UUID]bool),
	}

	roomuserModels, err := service.roomuserService.FindUserIdsByRoomId(
		context.Background(),
		roomuserPackage.FindUserIdsByRoomIdDTO{RoomId: roomModel.Id},
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
	r.handlers[USER_JOIN_ROOM_EVENT] = (*room).userJoinRoomEventHandler
	r.handlers[USER_LEAVE_ROOM_EVENT] = (*room).userLeaveRoomEventHandler
}

func (room *room) messageEventHandler(event event) {
	messageEvent, ok := event.(*messageEvent)
	if !ok {
		slog.Error("error typecasting to message event", "event", event)
		return
	}
	for userId := range room.userIds {
		user, ok := room.service.users[userId]
		if !ok {
			slog.Error("user not found", "userId", userId)
			continue
		}
		user.send <- messageEvent.payload
	}
}

func (room *room) userJoinRoomEventHandler(event event) {
	userJoinRoomEvent, ok := event.(*userJoinRoomEvent)
	if !ok {
		slog.Error("error typecasting to user join room event", "event", event)
		return
	}
	room.userIds[userJoinRoomEvent.userId] = true
}

func (room *room) userLeaveRoomEventHandler(event event) {
	userLeaveRoomEvent, ok := event.(*userLeaveRoomEvent)
	if !ok {
		slog.Error("error typecasting to user leave room event", "event", event)
		return
	}
	delete(room.userIds, userLeaveRoomEvent.userId)
}
