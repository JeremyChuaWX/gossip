package chat

import (
	"context"
	"gossip/internal/models"
	messagePackage "gossip/internal/services/message"
	roomPackage "gossip/internal/services/room"
	roomuserPackage "gossip/internal/services/roomuser"
	userPackage "gossip/internal/services/user"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  BUFFER_SIZE,
	WriteBufferSize: BUFFER_SIZE,
}

type Service struct {
	ingress  chan event
	alive    chan bool
	handlers map[eventName]func(*Service, event)

	userService     *userPackage.Service
	roomService     *roomPackage.Service
	roomuserService *roomuserPackage.Service
	messageService  *messagePackage.Service

	users map[uuid.UUID]*user
	rooms map[uuid.UUID]*room
}

func NewService(
	userService *userPackage.Service,
	roomService *roomPackage.Service,
	roomuserService *roomuserPackage.Service,
	messageService *messagePackage.Service,
) (*Service, error) {
	service := &Service{
		ingress:  make(chan event),
		alive:    make(chan bool),
		handlers: make(map[eventName]func(*Service, event)),

		userService:     userService,
		roomService:     roomService,
		roomuserService: roomuserService,
		messageService:  messageService,

		users: make(map[uuid.UUID]*user),
		rooms: make(map[uuid.UUID]*room),
	}

	roomModels, err := roomService.FindMany(context.Background())
	if err != nil {
		return nil, err
	}
	for _, roomModel := range roomModels {
		room, err := newRoom(service, &roomModel)
		if err != nil {
			return nil, err
		}
		service.rooms[roomModel.Id] = room
	}

	service.registerEventHandlers()

	go service.receiveEvents()

	return service, nil
}

func (service *Service) UserConnect(
	w http.ResponseWriter,
	r *http.Request,
	userModel *models.User,
) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	user, err := newUser(service, conn, userModel)
	if err != nil {
		return err
	}
	service.ingress <- &userConnectEvent{user: user}
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
			handler, ok := service.handlers[event.name()]
			if !ok {
				continue
			}
			handler(service, event)
		}
	}
}

func (service *Service) disconnect() {
}

// event management

func (s *Service) registerEventHandlers() {
	s.handlers[USER_CONNECT_EVENT] = (*Service).userConnectEventHandler
	s.handlers[USER_DISCONNECT_EVENT] = (*Service).userDisconnectEventHandler
}

func (service *Service) userConnectEventHandler(event event) {
	userConnectEvent, ok := event.(*userConnectEvent)
	if !ok {
		return
	}
	service.users[userConnectEvent.user.model.Id] = userConnectEvent.user
}

func (service *Service) userDisconnectEventHandler(event event) {
	userDisconnectEvent, ok := event.(*userDisconnectEvent)
	if !ok {
		return
	}
	delete(service.users, userDisconnectEvent.userId)
}
