package chat

import (
	"errors"
	"fmt"
	"gossip/internal/utils"
	"log"
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
	alive      chan bool
	handlers   map[eventType]func(*service, event)
}

func InitService() *service {
	wsUpgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	s := &service{
		wsUpgrader: wsUpgrader,
		rooms:      make(map[string]*room),
		clients:    make(map[uuid.UUID]*client),
		ingress:    make(chan event),
		alive:      make(chan bool),
		handlers:   make(map[eventType]func(*service, event)),
	}
	s.handlers[NEW_CLIENT] = (*service).newClientHandler
	s.handlers[CLIENT_DISCONNECT] = (*service).clientDisconnectHandler
	s.handlers[NEW_ROOM] = (*service).newRoomHandler
	s.handlers[DESTROY_ROOM] = (*service).destroyRoomHandler
	go s.receiveEvents()
	return s
}

func (s *service) InitRoutes(router *chi.Mux) {
	chatRouter := s.chatRouter()
	router.Mount("/chat", chatRouter)
}

func (s *service) chatRouter() *chi.Mux {
	chatRouter := chi.NewRouter()

	// new client
	chatRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Username string `query:"username"`
			UserId   string `query:"userId"`
		}
		query, err := utils.GetURLQueryStruct[request](r.URL)
		if err != nil {
			log.Println(err.Error())
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		userId, err := uuid.FromString(query.UserId)
		if err != nil {
			log.Println(err.Error())
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		conn, err := s.wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err.Error())
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			Message string `json:"message"`
		}
		if err := conn.WriteJSON(response{
			Message: "client connected",
		}); err != nil {
			conn.Close()
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		client := newClient(userId, query.Username, conn, s)
		s.ingress <- makeNewClientEvent(client)
	})

	// new room
	chatRouter.Post("/rooms", func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Name string `json:"name"`
		}
		body, err := utils.ReadJSON[request](r)
		if err != nil {
			log.Println(err.Error())
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		room := newRoom(body.Name, s)
		s.ingress <- makeNewRoomEvent(room)

		type response struct {
			Message string `json:"message"`
		}
		utils.WriteJSON(w, http.StatusCreated, response{
			Message: fmt.Sprintf("room created %s", room.name),
		})
	})

	// destroy room
	chatRouter.Delete(
		"/rooms/{name}",
		func(w http.ResponseWriter, r *http.Request) {
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

			event := makeDestroyRoomEvent(room)
			room.ingress <- event
			s.ingress <- event

			type response struct {
				Message string `json:"message"`
			}
			utils.WriteJSON(w, http.StatusCreated, response{
				Message: fmt.Sprintf("room destroyed %s", room.name),
			})
		},
	)

	// client join room
	chatRouter.Post(
		"/rooms/{name}",
		func(w http.ResponseWriter, r *http.Request) {
			roomName := chi.URLParam(r, "name")

			type request struct {
				UserId string `json:"userId"`
			}
			body, err := utils.ReadJSON[request](r)
			if err != nil {
				log.Println(err.Error())
				utils.WriteError(w, http.StatusBadRequest, err)
				return
			}

			userId, err := uuid.FromString(body.UserId)
			if err != nil {
				log.Println(err.Error())
				utils.WriteError(w, http.StatusBadRequest, err)
				return
			}

			client, ok := s.clients[userId]
			if !ok {
				utils.WriteError(
					w,
					http.StatusBadRequest,
					clientDoesNotExistError,
				)
				return
			}

			room, ok := s.rooms[roomName]
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
	defer (func() {
		close(s.ingress)
		close(s.alive)
	})()
	for {
		select {
		case <-s.alive:
			return
		case e := <-s.ingress:
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

func (s *service) newClientHandler(e event) {
	event := e.(*newClientEvent)
	s.clients[event.client.userId] = event.client
	event.client.init()
}

func (s *service) clientDisconnectHandler(e event) {
	event := e.(*clientDisconnectEvent)
	delete(s.clients, event.client.userId)
}

func (s *service) newRoomHandler(e event) {
	event := e.(*newRoomEvent)
	s.rooms[event.room.name] = event.room
	event.room.init()
}

func (s *service) destroyRoomHandler(e event) {
	event := e.(*destroyRoomEvent)
	delete(s.rooms, event.room.name)
}
