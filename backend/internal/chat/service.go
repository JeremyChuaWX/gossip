package chat

import (
	"errors"
	"fmt"
	"gossip/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
)

var (
	roomDoesNotExistError   = errors.New("room does not exist")
	clientDoesNotExistError = errors.New("client does not exist")
)

type service struct {
	wsUpgrader *websocket.Upgrader
	rooms      map[string]*room
	clients    map[uuid.UUID]*client
	ingress    chan event
}

func InitService() *service {
	wsUpgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	service := &service{
		wsUpgrader: wsUpgrader,
		rooms:      make(map[string]*room),
		clients:    make(map[uuid.UUID]*client),
		ingress:    make(chan event),
	}
	go service.receiveEvents()
	return service
}

func (s *service) InitRoutes(router *chi.Mux) {
	chatRouter := s.chatRouter()
	router.Mount("/chat", chatRouter)
}

func (s *service) chatRouter() *chi.Mux {
	chatRouter := chi.NewRouter()

	// new client
	chatRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}

		client := newClient(uuid.UUID{}, "", conn, s) // TODO: insert user info
		s.ingress <- makeNewClientEvent(client)

		type response struct {
			Message string `json:"message"`
		}
		utils.WriteJSON(w, http.StatusCreated, response{
			Message: fmt.Sprintf("client connected %s", client.userId.String()),
		})
	})

	// new room
	chatRouter.Post("/rooms", func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Name string `json:"name"`
		}
		req, err := utils.ReadJSON[request](r)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}

		room := newRoom(req.Name, s)
		s.ingress <- makeNewRoomEvent(room)

		type response struct {
			Message string `json:"message"`
		}
		utils.WriteJSON(w, http.StatusCreated, response{
			Message: fmt.Sprintf("room created %s", room.name),
		})
	})

	// client join room
	chatRouter.Post(
		"/rooms/{name}",
		func(w http.ResponseWriter, r *http.Request) {
			// TODO: userId from JWT
			userId := uuid.UUID{}

			client, ok := s.clients[userId]
			if !ok {
				utils.WriteError(
					w,
					http.StatusBadRequest,
					clientDoesNotExistError,
				)
				return
			}

			name := chi.URLParam(r, "name")

			room, ok := s.rooms[name]
			if !ok {
				utils.WriteError(
					w,
					http.StatusBadRequest,
					roomDoesNotExistError,
				)
				return
			}

			event := makeClientJoinRoomEvent(client, room)
			client.ingress <- event
			room.ingress <- event

			type response struct {
				Message string `json:"message"`
			}
			utils.WriteJSON(w, http.StatusCreated, response{
				Message: fmt.Sprintf("joined room %s", room.name),
			})
		},
	)

	return chatRouter
}

// run as goroutine
func (s *service) receiveEvents() {
	for {
		e := <-s.ingress
		switch e := e.(type) {
		case *newClientEvent:
			s.clients[e.client.userId] = e.client
			e.client.init()
		case *clientDisconnectEvent:
			delete(s.clients, e.client.userId)
		case *newRoomEvent:
			s.rooms[e.room.name] = e.room
			e.room.init()
		case *removeRoomEvent:
			delete(s.rooms, e.room.name)
		}
	}
}
