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
	results, err := service.repository.RoomFindMany(context.Background())
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		room, err := newRoom(service, result.RoomId)
		if err != nil {
			return nil, err
		}
		service.rooms[result.RoomId] = room
	}
	go service.receiveEvents()
	return service, nil
}

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

// actor methods

func (service *Service) receiveEvents() {
	for {
		select {
		case <-service.alive:
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
	case userConnectedEvent:
		s.userConnectedEventHandler(event)
	case userDisconnectedEvent:
		s.userDisconnectedEventHandler(event)
	default:
		slog.Error("invalid event", "event", event)
	}
}

func (service *Service) userConnectedEventHandler(event userConnectedEvent) {
	service.users[event.user.userId] = event.user
}

func (service *Service) userDisconnectedEventHandler(
	event userDisconnectedEvent,
) {
	delete(service.users, event.userId)
}
