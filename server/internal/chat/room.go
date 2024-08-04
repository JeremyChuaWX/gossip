package chat

import (
	"context"
	"gossip/internal/repository"
	"log/slog"

	"github.com/gofrs/uuid/v5"
)

type room struct {
	service *Service
	ingress chan event
	userIds map[uuid.UUID]bool
}

func newRoom(service *Service, roomId uuid.UUID) (*room, error) {
	room := &room{
		service: service,
		ingress: make(chan event),
		userIds: make(map[uuid.UUID]bool),
	}
	results, err := service.repository.UsersFindManyByRoomId(
		context.Background(),
		repository.UsersFindManyByRoomIdParams{RoomId: roomId},
	)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		room.userIds[result.UserId] = true
	}
	go room.receiveEvents()
	return room, nil
}

func (room *room) receiveEvents() {
	for {
		event, ok := <-room.ingress
		slog.Info("room received event", "event", event)
		if !ok {
			return
		}
		room.eventHandler(event)
	}
}

// event management

func (room *room) eventHandler(event event) {
	switch event := event.(type) {
	case messageEvent:
		room.messageEventHandler(event)
	case userJoinedRoomEvent:
		room.userJoinedRoomEventHandler(event)
	case userLeftRoomEvent:
		room.userLeftRoomEventHandler(event)
	default:
		slog.Error("invalid event", "event", event)
	}
}

func (room *room) messageEventHandler(event messageEvent) {
	if err := room.service.repository.MessageSave(
		context.Background(),
		repository.MessageSaveParams{
			UserId: event.userId,
			RoomId: event.roomId,
			Body:   event.payload.Body,
		},
	); err != nil {
		slog.Error("error saving message", "message", event.payload)
	}
	for userId := range room.userIds {
		user, ok := room.service.users[userId]
		if !ok {
			slog.Error("user not found", "userId", userId)
			continue
		}
		slog.Info("user found", "userId", userId)
		user.send <- event.payload
	}
}

func (room *room) userJoinedRoomEventHandler(event userJoinedRoomEvent) {
	room.userIds[event.userId] = true
}

func (room *room) userLeftRoomEventHandler(event userLeftRoomEvent) {
	delete(room.userIds, event.userId)
}
