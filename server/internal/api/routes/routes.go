package routes

import (
	"gossip/internal/api"
	"gossip/internal/api/middlewares"
	"gossip/internal/chat"
	"gossip/internal/repository"
	"gossip/internal/utils/httpjson"
	"gossip/internal/utils/password"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gofrs/uuid/v5"
)

type Router struct {
	Repository  *repository.Repository
	ChatService *chat.Service
	Middlewares *middlewares.Middlewares
}

func (r *Router) Start(address string) error {
	router := chi.NewMux()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Group(r.registerRoutes)
	router.Group(r.registerAuthedRoutes)

	if err := chi.Walk(router, api.WalkRoutes); err != nil {
		return err
	}

	return http.ListenAndServe(address, router)
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
		user, err := router.Repository.UserCreate(
			r.Context(),
			repository.UserCreateParams{
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
				User any `json:"user"`
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
		user, err := router.Repository.UserFindOneByUsername(
			r.Context(),
			repository.UserFindOneByUsernameParams{Username: body.Username},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		err = password.Verify(body.Password, user.PasswordHash)
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		session, err := router.Repository.UserSessionCreate(
			r.Context(),
			repository.UserSessionCreateParams{UserId: user.UserId},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		httpjson.Write(w, http.StatusOK, httpjson.BaseResponse{
			Success: true,
			Message: "logged in",
			Data: struct {
				Session any `json:"session"`
			}{
				Session: session,
			},
		})
	})
}

func (router *Router) registerAuthedRoutes(mux chi.Router) {
	mux.Use(router.Middlewares.AuthMiddleware)

	mux.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		sessionId, err := api.SessionIdFromHeader(r.Header)
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		err = router.Repository.UserSessionDelete(
			r.Context(),
			repository.UserSessionDeleteParams{SessionId: sessionId},
		)
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
		userSession := api.UserSessionFromContext(r.Context())
		httpjson.Write(w, http.StatusOK, httpjson.BaseResponse{
			Success: true,
			Message: "logged in user",
			Data: struct {
				User any `json:"user"`
			}{
				User: userSession,
			},
		})
	})

	mux.Get("/rooms", func(w http.ResponseWriter, r *http.Request) {
		userSession := api.UserSessionFromContext(r.Context())
		rooms, err := router.Repository.RoomFindManyByUserId(
			r.Context(),
			repository.RoomFindManyByUserIdParams{
				UserId: userSession.UserId,
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
				Rooms any
			}{
				Rooms: rooms,
			},
		})
	})

	mux.Post("/rooms/create", func(w http.ResponseWriter, r *http.Request) {
		userSession := api.UserSessionFromContext(r.Context())
		body, err := httpjson.Read[struct {
			RoomName string `json:"roomName"`
		}](r)
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, err)
			return
		}
		room, err := router.Repository.RoomCreate(
			r.Context(),
			repository.RoomCreateParams{Name: body.RoomName},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		err = router.Repository.UserJoinRoom(
			r.Context(),
			repository.UserJoinRoomParams{
				UserId: userSession.UserId,
				RoomId: room.RoomId,
			},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		httpjson.Write(w, http.StatusOK, httpjson.BaseResponse{
			Success: true,
			Message: "room created",
			Data: struct {
				Room any `json:"room"`
			}{
				Room: room,
			},
		})
	})

	mux.Post("/rooms/join", func(w http.ResponseWriter, r *http.Request) {
		userSession := api.UserSessionFromContext(r.Context())
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
		err = router.Repository.UserJoinRoom(
			r.Context(),
			repository.UserJoinRoomParams{
				UserId: userSession.UserId,
				RoomId: roomId,
			},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		httpjson.Write(w, http.StatusOK, httpjson.BaseResponse{
			Success: true,
			Message: "user joined room",
		})
	})

	mux.Post("/rooms/leave", func(w http.ResponseWriter, r *http.Request) {
		userSession := api.UserSessionFromContext(r.Context())
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
		err = router.Repository.UserLeaveRoom(
			r.Context(),
			repository.UserLeaveRoomParams{
				UserId: userSession.UserId,
				RoomId: roomId,
			},
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		httpjson.Write(w, http.StatusOK, httpjson.BaseResponse{
			Success: true,
			Message: "user left room",
		})
	})

	mux.Get("/rooms/connect", func(w http.ResponseWriter, r *http.Request) {
		userSession := api.UserSessionFromContext(r.Context())
		err := router.ChatService.UserConnect(
			w,
			r,
			userSession.UserId,
			userSession.Username,
		)
		if err != nil {
			httpjson.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	})
}
