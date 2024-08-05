package chat

import (
	"context"
	"fmt"
	"gossip/internal/repository"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  BUFFER_SIZE,
	WriteBufferSize: BUFFER_SIZE,
}

type Service struct {
	ingress    chan event
	repository *repository.Repository
	users      map[uuid.UUID]*user
	rooms      map[uuid.UUID]*room
}

func NewService(repository *repository.Repository) (*Service, error) {
	service := &Service{
		ingress:    make(chan event),
		repository: repository,
		users:      make(map[uuid.UUID]*user),
		rooms:      make(map[uuid.UUID]*room),
	}
	service.initRooms()
	go service.receiveEvents()
	return service, nil
}

func (service *Service) initRooms() {
	results, err := service.repository.RoomFindMany(context.Background())
	if err != nil {
		slog.Error("error finding rooms", "error", err.Error())
	}
	if len(results) < 1 {
		slog.Error("no rooms found")
		return
	}
	for _, result := range results {
		room, err := newRoom(service, result.RoomId)
		if err != nil {
			slog.Error("error initing room", "roomId", result.RoomId)
		}
		service.rooms[result.RoomId] = room
		slog.Info("inited room", "roomId", result.RoomId)
	}
}

// API methods

func (service *Service) UserConnect(
	w http.ResponseWriter,
	r *http.Request,
	userId uuid.UUID,
	username string,
) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	user := newUser(service, conn, userId, username)
	service.ingress <- userConnectedEvent{user: user}
	return nil
}

func (service *Service) RoomCreate(roomId uuid.UUID) error {
	room, err := newRoom(service, roomId)
	if err != nil {
		return err
	}
	service.ingress <- roomCreatedEvent{room: room}
	return nil
}

func (service *Service) receiveEvents() {
	for {
		event, ok := <-service.ingress
		if !ok {
			return
		}
		service.eventHandler(event)
	}
}

// event management

func (s *Service) eventHandler(event event) {
	switch event := event.(type) {
	case roomCreatedEvent:
		s.roomCreatedEventHandler(event)
	case userConnectedEvent:
		s.userConnectedEventHandler(event)
	case userDisconnectedEvent:
		s.userDisconnectedEventHandler(event)
	default:
		slog.Error("invalid event", "event", event)
	}
}

func (service *Service) roomCreatedEventHandler(event roomCreatedEvent) {
	service.rooms[event.room.roomId] = event.room
}

func (service *Service) userConnectedEventHandler(event userConnectedEvent) {
	if user, ok := service.users[event.user.userId]; ok {
		slog.Info(
			"existing user connection found",
			"userId", user.userId,
			"address", fmt.Sprintf("%p", user),
		)
		user.disconnect()
	}
	service.users[event.user.userId] = event.user
	slog.Info(
		"connecting user",
		"userId", event.user.userId,
		"address", fmt.Sprintf("%p", event.user),
	)
	if user, ok := service.users[event.user.userId]; ok {
		slog.Info("user connected",
			"userId", user.userId,
			"address", fmt.Sprintf("%p", user),
		)
	}
}

func (service *Service) userDisconnectedEventHandler(
	event userDisconnectedEvent,
) {
	delete(service.users, event.userId)
}
