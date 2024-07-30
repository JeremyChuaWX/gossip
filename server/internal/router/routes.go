package router

import (
	"gossip/internal/repository"
	"gossip/internal/utils/password"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

func (router *Router) routeGroup(mux chi.Router) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/index.html")
	})

	mux.Get("/signup", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/signup.html")
	})

	mux.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
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
		user, err := router.Repository.UserCreate(
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

	mux.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/login.html")
	})

	mux.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		body, err := readJSON[struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}](r)
		if err != nil {
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		user, err := router.Repository.UserFindOneByUsername(
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
		session, err := router.Repository.UserSessionCreate(
			r.Context(),
			repository.UserSessionCreateParams{UserId: user.UserId},
		)
		if err != nil {
			errorToJSON(w, http.StatusUnauthorized, err)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     SESSION_ID_COOKIE,
			Value:    session.SessionId.String(),
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

func (router *Router) authedRouteGroup(mux chi.Router) {
	mux.Use(router.authMiddleware)

	mux.Get("/connect", func(w http.ResponseWriter, r *http.Request) {
		userSession := userSessionFromContext(r.Context())
		err := router.ChatService.UserConnect(
			w,
			r,
			userSession.UserId,
			userSession.Username,
		)
		if err != nil {
			slog.Error("error creating WS connection")
			return
		}
	})

	mux.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		userSession := userSessionFromContext(r.Context())
		rooms, err := router.Repository.RoomFindManyByUserId(
			r.Context(),
			repository.RoomFindManyByUserIdParams{UserId: userSession.UserId},
		)
		if err != nil {
			slog.Error(
				"error finding rooms for user",
				"userSession",
				userSession,
			)
			return
		}
		t, err := template.ParseFiles("pages/home.html")
		if err != nil {
			slog.Error("error parsing home.html", "error", err)
			return
		}
		err = t.Execute(w, rooms)
		if err != nil {
			slog.Error("error executing home.html template", "error", err)
			return
		}
	})

	mux.Get("/rooms/create", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/create-room.html")
	})

	mux.Post("/rooms/create", func(w http.ResponseWriter, r *http.Request) {
		userSession := userSessionFromContext(r.Context())
		body, err := readJSON[struct {
			RoomName string `json:"roomName"`
		}](r)
		if err != nil {
			slog.Error("invalid body for create room", "error", err)
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		room, err := router.Repository.RoomCreate(
			r.Context(),
			repository.RoomCreateParams{Name: body.RoomName},
		)
		if err != nil {
			slog.Error("error creating room", "error", err)
			errorToJSON(w, http.StatusInternalServerError, err)
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
			slog.Error("error joining room", "error", err)
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "room created",
		})
	})

	mux.Get("/rooms/join", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/join-room.html")
	})

	mux.Post("/rooms/join", func(w http.ResponseWriter, r *http.Request) {
		userSession := userSessionFromContext(r.Context())
		body, err := readJSON[struct {
			RoomId string `json:"roomId"`
		}](r)
		if err != nil {
			slog.Error("invalid body for create room", "error", err)
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		roomId, err := uuid.FromString(body.RoomId)
		if err != nil {
			slog.Error("invalid room ID", "body", body)
			errorToJSON(w, http.StatusBadRequest, err)
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
			slog.Error("error joining room", "error", err)
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "room joined",
		})
	})

	mux.Get("/rooms/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		userSession := userSessionFromContext(r.Context())
		roomIdParamValue := chi.URLParam(r, "roomId")
		if roomIdParamValue == "" {
			slog.Error("invalid room ID")
			return
		}
		roomId, err := uuid.FromString(roomIdParamValue)
		if err != nil {
			slog.Error("invalid room ID", "roomIdParamValue", roomIdParamValue)
			return
		}
		isMember, err := router.Repository.UserCheckRoomMembership(
			r.Context(),
			repository.UserCheckRoomMembershipParams{
				UserId: userSession.UserId,
				RoomId: roomId,
			},
		)
		if err != nil || !isMember {
			slog.Error(
				"user not in room",
				"userId",
				userSession.UserId,
				"roomId",
				roomId,
			)
			return
		}
		room, err := router.Repository.RoomFindOne(
			r.Context(),
			repository.RoomFindOneParams{RoomId: roomId},
		)
		if err != nil {
			slog.Error("error finding room", "roomId", roomId)
			return
		}
		messages, err := router.Repository.MessagesFindManyByRoomId(
			r.Context(),
			repository.MessagesFindManyByRoomIdParams{RoomId: roomId},
		)
		if err != nil {
			slog.Error("error room messages", "roomId", roomId)
		}
		t, err := template.ParseFiles("pages/room.html")
		if err != nil {
			slog.Error("error parsing room.html", "error", err)
			return
		}
		err = t.Execute(w, map[string]any{
			"name":     room.Name,
			"messages": messages,
		})
		if err != nil {
			slog.Error("error executing room.html template", "error", err)
			return
		}
	})
}
