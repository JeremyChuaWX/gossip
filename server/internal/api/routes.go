package api

import (
	"gossip/internal/repository"
	"gossip/internal/utils/password"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

func (api *Api) routesGroup(router chi.Router) {
	router.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
		body, err := readJSON[struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}](r)
		if err != nil {
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		passwordHash, err := password.Hash(body.Password)
		if err != nil {
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		user, err := api.Repository.UserCreate(
			r.Context(),
			repository.UserCreateParams{
				Username:     body.Username,
				PasswordHash: passwordHash,
			},
		)
		if err != nil {
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

	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		body, err := readJSON[struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}](r)
		if err != nil {
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		user, err := api.Repository.UserFindOneByUsername(
			r.Context(),
			repository.UserFindOneByUsernameParams{Username: body.Username},
		)
		if err != nil {
			errorToJSON(w, http.StatusUnauthorized, err)
			return
		}
		err = password.Verify(body.Password, user.PasswordHash)
		if err != nil {
			errorToJSON(w, http.StatusUnauthorized, err)
			return
		}
		session, err := api.Repository.UserSessionCreate(
			r.Context(),
			repository.UserSessionCreateParams{UserId: user.UserId},
		)
		if err != nil {
			errorToJSON(w, http.StatusUnauthorized, err)
			return
		}
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

func (api *Api) authedRoutesGroup(router chi.Router) {
	router.Use(api.authMiddleware)

	router.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		sessionId, err := sessionIdFromHeader(r.Header)
		if err != nil {
			errorToJSON(w, http.StatusUnauthorized, err)
			return
		}
		err = api.Repository.UserSessionDelete(
			r.Context(),
			repository.UserSessionDeleteParams{SessionId: sessionId},
		)
		if err != nil {
			errorToJSON(w, http.StatusUnauthorized, err)
			return
		}
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "logged out",
		})
	})

	router.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		userSession := userSessionFromContext(r.Context())
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "logged in user",
			Data: map[string]any{
				"user": map[string]any{
					"id":       userSession.UserId,
					"username": userSession.Username,
				},
			},
		})
	})

	router.Get("/rooms", func(w http.ResponseWriter, r *http.Request) {
		userSession := userSessionFromContext(r.Context())
		rooms, err := api.Repository.RoomFindManyByUserId(
			r.Context(),
			repository.RoomFindManyByUserIdParams{
				UserId: userSession.UserId,
			},
		)
		if err != nil {
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		formattedRooms := []map[string]any{}
		for _, room := range rooms {
			formattedRooms = append(formattedRooms, map[string]any{
				"id":   room.RoomId,
				"name": room.Name,
			})
		}
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "user rooms found",
			Data: map[string]any{
				"rooms": formattedRooms,
			},
		})
	})

	router.Get("/rooms/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		roomIdParam := chi.URLParam(r, "roomId")
		roomId, err := uuid.FromString(roomIdParam)
		if err != nil {
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		userSession := userSessionFromContext(r.Context())
		isMember, err := api.Repository.UserCheckRoomMembership(
			r.Context(),
			repository.UserCheckRoomMembershipParams{
				UserId: userSession.UserId,
				RoomId: roomId,
			},
		)
		if err != nil {
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		if !isMember {
			writeJSON(w, http.StatusForbidden, baseResponse{
				Success: true,
				Message: "user not in room",
			})
			return
		}
		room, err := api.Repository.RoomFindOne(
			r.Context(),
			repository.RoomFindOneParams{RoomId: roomId},
		)
		if err != nil {
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		messages, err := api.Repository.MessagesFindManyByRoomId(
			r.Context(),
			repository.MessagesFindManyByRoomIdParams{RoomId: roomId},
		)
		if err != nil {
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		formattedMessages := []map[string]any{}
		for _, message := range messages {
			formattedMessages = append(formattedMessages, map[string]any{
				"id":        message.MessageId,
				"userId":    message.UserId,
				"username":  message.Username,
				"body":      message.Body,
				"timestamp": message.Timestamp,
			})
		}
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "room messages found",
			Data: map[string]any{
				"name":     room.Name,
				"messages": formattedMessages,
			},
		})
	})

	router.Post("/rooms/create", func(w http.ResponseWriter, r *http.Request) {
		userSession := userSessionFromContext(r.Context())
		body, err := readJSON[struct {
			RoomName string `json:"roomName"`
		}](r)
		if err != nil {
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		room, err := api.Repository.RoomCreate(
			r.Context(),
			repository.RoomCreateParams{Name: body.RoomName},
		)
		if err != nil {
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		err = api.Repository.UserJoinRoom(
			r.Context(),
			repository.UserJoinRoomParams{
				UserId: userSession.UserId,
				RoomId: room.RoomId,
			},
		)
		if err != nil {
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		api.ChatService.RoomCreate(room.RoomId)
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "room created",
			Data: map[string]any{
				"room": map[string]any{
					"id": room.RoomId,
				},
			},
		})
	})

	router.Post("/rooms/join", func(w http.ResponseWriter, r *http.Request) {
		userSession := userSessionFromContext(r.Context())
		body, err := readJSON[struct {
			RoomId string `json:"roomId"`
		}](r)
		if err != nil {
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		roomId, err := uuid.FromString(body.RoomId)
		if err != nil {
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		err = api.Repository.UserJoinRoom(
			r.Context(),
			repository.UserJoinRoomParams{
				UserId: userSession.UserId,
				RoomId: roomId,
			},
		)
		if err != nil {
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "user joined room",
		})
	})

	router.Post("/rooms/leave", func(w http.ResponseWriter, r *http.Request) {
		userSession := userSessionFromContext(r.Context())
		body, err := readJSON[struct {
			RoomId string `json:"roomId"`
		}](r)
		if err != nil {
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		roomId, err := uuid.FromString(body.RoomId)
		if err != nil {
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		err = api.Repository.UserLeaveRoom(
			r.Context(),
			repository.UserLeaveRoomParams{
				UserId: userSession.UserId,
				RoomId: roomId,
			},
		)
		if err != nil {
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "user left room",
		})
	})

	router.Get("/rooms/connect", func(w http.ResponseWriter, r *http.Request) {
		userSession := userSessionFromContext(r.Context())
		err := api.ChatService.UserConnect(
			w,
			r,
			userSession.UserId,
			userSession.Username,
		)
		if err != nil {
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
	})
}
