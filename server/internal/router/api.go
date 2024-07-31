package router

import (
	"gossip/internal/repository"
	"gossip/internal/utils/password"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

func (router *Router) apiRouter() *chi.Mux {
	api := chi.NewMux()
	api.Group(router.apiRouteGroup)
	api.Group(router.apiAuthedRouteGroup)
	return api
}

func (router *Router) apiRouteGroup(mux chi.Router) {
	mux.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
		body, err := readJSON[struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}](r)
		if err != nil {
			slog.Error("error parsing body")
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		passwordHash, err := password.Hash(body.Password)
		if err != nil {
			slog.Error("error hashing password", "body.Password", body.Password)
			errorToJSON(w, http.StatusBadRequest, err)
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
			slog.Error("error creating user")
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "signed up",
			Data: map[string]any{
				"user": map[string]any{
					"id": user.UserId,
				},
			},
		})
	})

	mux.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		body, err := readJSON[struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}](r)
		if err != nil {
			slog.Error("error parsing body")
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		user, err := router.Repository.UserFindOneByUsername(
			r.Context(),
			repository.UserFindOneByUsernameParams{Username: body.Username},
		)
		if err != nil {
			slog.Error("error finding user", "body.Username", body.Username)
			errorToJSON(w, http.StatusUnauthorized, err)
			return
		}
		err = password.Verify(body.Password, user.PasswordHash)
		if err != nil {
			slog.Error("error verifying password")
			errorToJSON(w, http.StatusUnauthorized, err)
			return
		}
		session, err := router.Repository.SessionCreate(
			r.Context(),
			repository.SessionCreateParams{UserId: user.UserId},
		)
		if err != nil {
			slog.Error("error creating session", "userId", user.UserId)
			errorToJSON(w, http.StatusUnauthorized, err)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     SESSION_ID_COOKIE,
			Value:    session.SessionId.String(),
			Path:     "/",
			Expires:  session.ExpiresOn,
			Secure:   true,
			HttpOnly: true,
		})
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "logged in",
			Data: map[string]any{
				"session": map[string]any{
					"id":        session.SessionId,
					"expiresOn": session.ExpiresOn,
				},
			},
		})
	})
}

func (router *Router) apiAuthedRouteGroup(mux chi.Router) {
	mux.Use(router.authMiddleware)

	mux.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		session := sessionFromContextSafe(r.Context())
		err := router.Repository.SessionDelete(
			r.Context(),
			repository.SessionDeleteParams{
				SessionId: session.SessionId,
			},
		)
		if err != nil {
			slog.Error("error deleting session", "session", session)
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     SESSION_ID_COOKIE,
			MaxAge:   0,
			Secure:   true,
			HttpOnly: true,
		})
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "logged out",
		})
	})

	mux.Get("/connect", func(w http.ResponseWriter, r *http.Request) {
		session := sessionFromContextSafe(r.Context())
		err := router.ChatService.UserConnect(
			w,
			r,
			session.UserId,
			session.Username,
		)
		if err != nil {
			slog.Error("error creating WS connection")
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
	})

	mux.Post("/rooms/create", func(w http.ResponseWriter, r *http.Request) {
		session := sessionFromContextSafe(r.Context())
		body, err := readJSON[struct {
			RoomName string `json:"roomName"`
		}](r)
		if err != nil {
			slog.Error("error parsing body")
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		room, err := router.Repository.RoomCreate(
			r.Context(),
			repository.RoomCreateParams{Name: body.RoomName},
		)
		if err != nil {
			slog.Error("error creating room")
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		err = router.Repository.UserJoinRoom(
			r.Context(),
			repository.UserJoinRoomParams{
				UserId: session.UserId,
				RoomId: room.RoomId,
			},
		)
		if err != nil {
			slog.Error("error joining room")
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "room created",
		})
	})

	mux.Post("/rooms/join", func(w http.ResponseWriter, r *http.Request) {
		session := sessionFromContextSafe(r.Context())
		body, err := readJSON[struct {
			RoomId string `json:"roomId"`
		}](r)
		if err != nil {
			slog.Error("error parsing body")
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		roomId, err := uuid.FromString(body.RoomId)
		if err != nil {
			slog.Error("error parsing roomId", "body.RoomId", body.RoomId)
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		err = router.Repository.UserJoinRoom(
			r.Context(),
			repository.UserJoinRoomParams{
				UserId: session.UserId,
				RoomId: roomId,
			},
		)
		if err != nil {
			slog.Error("error joining room")
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "room joined",
		})
	})

	mux.Post("/rooms/leave", func(w http.ResponseWriter, r *http.Request) {
		session := sessionFromContextSafe(r.Context())
		body, err := readJSON[struct {
			RoomId string `json:"roomId"`
		}](r)
		if err != nil {
			slog.Error("error parsing body")
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		roomId, err := uuid.FromString(body.RoomId)
		if err != nil {
			slog.Error("error parsing roomId", "body.RoomId", body.RoomId)
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		err = router.Repository.UserLeaveRoom(
			r.Context(),
			repository.UserLeaveRoomParams{
				UserId: session.UserId,
				RoomId: roomId,
			},
		)
		if err != nil {
			slog.Error("error leaving room")
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "room left",
		})
	})
}
