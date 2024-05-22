package routes

import (
	"gossip/internal/api"
	"gossip/internal/api/middlewares"
	"gossip/internal/models"
	"gossip/internal/services/chat"
	"gossip/internal/services/room"
	"gossip/internal/services/roomuser"
	"gossip/internal/services/session"
	"gossip/internal/services/user"
	"gossip/internal/utils/httpjson"
	"gossip/internal/utils/password"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gofrs/uuid/v5"
)

type Router struct {
	RoomService     *room.Service
	UserService     *user.Service
	RoomUserService *roomuser.Service
	SessionService  *session.Service
	ChatService     *chat.Service
	Middlewares     *middlewares.Middlewares
}

func (r *Router) Start(address string) error {
	router := chi.NewMux()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Group(r.registerRoutes)
	router.Group(r.registerAuthedRoutes)

	if err := chi.Walk(router, walkRoutes); err != nil {
		return err
	}

	return http.ListenAndServe(address, router)
}

func walkRoutes(
	method string,
	route string,
	handler http.Handler,
	middlewares ...func(http.Handler) http.Handler,
) error {
	route = strings.Replace(route, "/*/", "/", -1)
	log.Printf("%s %s\n", method, route)
	return nil
}

func (router *Router) registerRoutes(mux chi.Router) {
	mux.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
		body, err := httpjson.Read[struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}](r)
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, err)
			return
		}
		passwordHash, err := password.Hash(body.Password)
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, err)
			return
		}
		user, err := router.UserService.Create(
			r.Context(),
			user.CreateDTO{
				Username:     body.Username,
				PasswordHash: passwordHash,
			},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		httpjson.Write(w, http.StatusOK, httpjson.BaseResponse{
			Success: true,
			Message: "signed up",
			Data: struct {
				User models.User `json:"user"`
			}{
				User: user,
			},
		})
	})

	mux.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		body, err := httpjson.Read[struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}](r)
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, err)
			return
		}
		user, err := router.UserService.FindOneByUsername(
			r.Context(),
			user.FindOneByUsernameDTO{
				Username: body.Username,
			},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		sessionId, err := router.SessionService.Create(r.Context(), user.Id)
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		err = password.Verify(body.Password, user.PasswordHash)
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		httpjson.Write(w, http.StatusOK, httpjson.BaseResponse{
			Success: true,
			Message: "logged in",
			Data: struct {
				SessionId string `json:"sessionId"`
			}{
				SessionId: sessionId,
			},
		})
	})
}

func (router *Router) registerAuthedRoutes(mux chi.Router) {
	mux.Use(router.Middlewares.AuthMiddleware)

	mux.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.Header.Get(api.SESSION_ID_HEADER)
		err := router.SessionService.Delete(r.Context(), sessionId)
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		httpjson.Write(w, http.StatusOK, httpjson.BaseResponse{
			Success: true,
			Message: "logged out",
		})
	})

	mux.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(api.USER_ID_CONTEXT_KEY).(uuid.UUID)
		user, err := router.UserService.FindOne(r.Context(), user.FindOneDTO{
			UserId: userId,
		})
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		httpjson.Write(w, http.StatusOK, httpjson.BaseResponse{
			Success: true,
			Message: "logged in user",
			Data: struct {
				User models.User `json:"user"`
			}{
				User: user,
			},
		})
	})

	mux.Get("/users/rooms", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(api.USER_ID_CONTEXT_KEY).(uuid.UUID)
		user, err := router.UserService.FindOne(r.Context(), user.FindOneDTO{
			UserId: userId,
		})
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		rooms, err := router.RoomUserService.FindRoomsByUserId(
			r.Context(),
			roomuser.FindRoomIdsByUserIdDTO{
				UserId: user.Id,
			},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		httpjson.Write(w, http.StatusOK, httpjson.BaseResponse{
			Success: true,
			Message: "user rooms found",
			Data: struct {
				Rooms []models.Room `json:"rooms"`
			}{
				Rooms: rooms,
			},
		})
	})

	mux.Post("/users/join-room", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(api.USER_ID_CONTEXT_KEY).(uuid.UUID)
		user, err := router.UserService.FindOne(
			r.Context(),
			user.FindOneDTO{
				UserId: userId,
			},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		body, err := httpjson.Read[struct {
			RoomId string `json:"roomId"`
		}](r)
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, err)
			return
		}
		roomId, err := uuid.FromString(body.RoomId)
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, err)
			return
		}
		_, err = router.RoomUserService.UserJoinRoom(
			r.Context(),
			roomuser.UserJoinRoomDTO{
				UserId: user.Id,
				RoomId: roomId,
			},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, err)
			return
		}
		httpjson.Write(w, http.StatusOK, httpjson.BaseResponse{
			Success: true,
			Message: "user joined room",
		})
	})

	mux.Post("/users/leave-room", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(api.USER_ID_CONTEXT_KEY).(uuid.UUID)
		user, err := router.UserService.FindOne(
			r.Context(),
			user.FindOneDTO{
				UserId: userId,
			},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		body, err := httpjson.Read[struct {
			RoomId string `json:"roomId"`
		}](r)
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, err)
			return
		}
		roomId, err := uuid.FromString(body.RoomId)
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, err)
			return
		}
		_, err = router.RoomUserService.UserLeaveRoom(
			r.Context(),
			roomuser.UserLeaveRoomDTO{
				UserId: user.Id,
				RoomId: roomId,
			},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, err)
			return
		}
		httpjson.Write(w, http.StatusOK, httpjson.BaseResponse{
			Success: true,
			Message: "user left room",
		})
	})

	mux.Post("/rooms", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(api.USER_ID_CONTEXT_KEY).(uuid.UUID)
		_, err := router.UserService.FindOne(
			r.Context(),
			user.FindOneDTO{
				UserId: userId,
			},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		body, err := httpjson.Read[struct {
			RoomName string `json:"roomName"`
		}](r)
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, err)
			return
		}
		room, err := router.RoomService.Create(r.Context(), room.CreateDTO{
			Name: body.RoomName,
		})
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		httpjson.Write(w, http.StatusOK, httpjson.BaseResponse{
			Success: true,
			Message: "room created",
			Data: struct {
				Room models.Room `json:"room"`
			}{
				Room: room,
			},
		})
	})
}
