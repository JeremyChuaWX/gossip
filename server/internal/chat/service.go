package chat

import (
	"context"
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
	alive      chan bool
	repository *repository.Repository
	users      map[uuid.UUID]*user
	rooms      map[uuid.UUID]*room
}

func NewService(repository *repository.Repository) (*Service, error) {
	service := &Service{
		ingress:    make(chan event),
		alive:      make(chan bool),
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
	user, err := newUser(service, conn, userId, username)
	if err != nil {
		return err
	}
	service.ingress <- userConnectedEvent{user: user}
	return nil
}

func (service *Service) RoomCreate(roomId uuid.UUID) {
	service.ingress <- roomCreatedEvent{roomId: roomId}
}

// actor methods

func (service *Service) receiveEvents() {
	for {
		select {
		case <-service.alive:
			service.disconnect()
			return
		case event, ok := <-service.ingress:
			if !ok {
				return
			}
			service.eventHandler(event)
		}
	}
}

func (service *Service) disconnect() {
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
	room, err := newRoom(service, event.roomId)
	if err != nil {
		slog.Error("error creating room", "roomId", event.roomId)
	}
	service.rooms[event.roomId] = room
}

func (service *Service) userConnectedEventHandler(event userConnectedEvent) {
	service.users[event.user.userId] = event.user
}

func (service *Service) userDisconnectedEventHandler(
	event userDisconnectedEvent,
) {
	delete(service.users, event.userId)
}
