package chat

import (
	"context"
	"gossip/internal/models"
	roomuserPackage "gossip/internal/services/roomuser"
	"log/slog"

	"github.com/gofrs/uuid/v5"
)

type room struct {
	model   *models.Room
	service *Service
	ingress chan event
	alive   chan bool

	userIds map[uuid.UUID]bool
}

func newRoom(service *Service, roomModel *models.Room) (*room, error) {
	room := &room{
		model:   roomModel,
		service: service,
		ingress: make(chan event),
		alive:   make(chan bool),

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
			room.eventHandler(event)
		}
	}
}

func (room *room) disconnect() {
}

// event management

func (room *room) eventHandler(event event) {
	switch event := event.(type) {
	case messageEvent:
		room.messageEventHandler(event)
	case userJoinRoomEvent:
		room.userJoinRoomEventHandler(event)
	case userLeaveRoomEvent:
		room.userLeaveRoomEventHandler(event)
	default:
		slog.Error("invalid event", "event", event)
	}
}

func (room *room) messageEventHandler(event messageEvent) {
	for userId := range room.userIds {
		user, ok := room.service.users[userId]
		if !ok {
			slog.Error("user not found", "userId", userId)
			continue
		}
		user.send <- event.payload
	}
}

func (room *room) userJoinRoomEventHandler(event userJoinRoomEvent) {
	room.userIds[event.userId] = true
}

func (room *room) userLeaveRoomEventHandler(event userLeaveRoomEvent) {
	delete(room.userIds, event.userId)
}
