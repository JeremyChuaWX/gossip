package chat

import (
	"errors"
	"log"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
)

var (
	roomDoesNotExistError   = errors.New("room does not exist")
	clientDoesNotExistError = errors.New("client does not exist")
)

type Service struct {
	wsUpgrader *websocket.Upgrader
	rooms      map[string]*room
	clients    map[uuid.UUID]*client
	ingress    chan event
	alive      chan bool
	handlers   map[eventType]func(*Service, event)
}

func InitService() *Service {
	wsUpgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	s := &Service{
		wsUpgrader: wsUpgrader,
		rooms:      make(map[string]*room),
		clients:    make(map[uuid.UUID]*client),
		ingress:    make(chan event, 100),
		alive:      make(chan bool),
		handlers:   make(map[eventType]func(*Service, event)),
	}
	s.handlers[NEW_CLIENT] = (*Service).newClientHandler
	s.handlers[CLIENT_DISCONNECT] = (*Service).clientDisconnectHandler
	s.handlers[NEW_ROOM] = (*Service).newRoomHandler
	s.handlers[DESTROY_ROOM] = (*Service).destroyRoomHandler
	s.handlers[MESSAGE] = (*Service).messageHandler
	go s.receiveEvents()
	return s
}

func (s *Service) UpgradeConnection(
	w http.ResponseWriter,
	r *http.Request,
) (*websocket.Conn, error) {
	conn, err := s.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	type response struct {
		Message string `json:"message"`
	}
	if err := conn.WriteJSON(response{
		Message: "client connected",
	}); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func (s *Service) NewClient(
	userId uuid.UUID,
	username string,
	conn *websocket.Conn,
) {
	client := newClient(userId, username, conn, s)
	s.ingress <- makeNewClientEvent(client)
}

func (s *Service) NewRoom(name string) {
	room := newRoom(name, s)
	s.ingress <- makeNewRoomEvent(room)
}

func (s *Service) GetRooms() []string {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	return rooms
}

func (s *Service) DestroyRoom(name string) error {
	room, ok := s.rooms[name]
	if !ok {
		return roomDoesNotExistError
	}

	event := makeDestroyRoomEvent(room)
	room.ingress <- event
	s.ingress <- event

	return nil
}

func (s *Service) ClientJoinRoom(userId uuid.UUID, roomName string) error {
	client, ok := s.clients[userId]
	if !ok {
		return clientDoesNotExistError
	}

	room, ok := s.rooms[roomName]
	if !ok {
		return roomDoesNotExistError
	}

	event := makeClientJoinRoomEvent(client, room)
	client.ingress <- event
	room.ingress <- event

	return nil
}

// run as goroutine
func (s *Service) receiveEvents() {
	defer (func() {
		close(s.ingress)
		close(s.alive)
	})()
	for {
		select {
		case <-s.alive:
			return
		case e := <-s.ingress:
			log.Printf("[service] event received: %v", e)
			handler, ok := s.handlers[e.name()]
			if !ok {
				log.Println("invalid event")
				continue
			}
			handler(s, e)
		}
	}
}

// handlers

func (s *Service) newClientHandler(e event) {
	event := e.(*newClientEvent)
	s.clients[event.client.userId] = event.client
	event.client.init()
}

func (s *Service) clientDisconnectHandler(e event) {
	event := e.(*clientDisconnectEvent)
	delete(s.clients, event.client.userId)
}

func (s *Service) newRoomHandler(e event) {
	event := e.(*newRoomEvent)
	s.rooms[event.room.name] = event.room
	event.room.init()
}

func (s *Service) destroyRoomHandler(e event) {
	event := e.(*destroyRoomEvent)
	delete(s.rooms, event.room.name)
}

func (s *Service) messageHandler(e event) {
	event := e.(*messageEvent)
	room := s.rooms[event.room]
	room.ingress <- event
}
