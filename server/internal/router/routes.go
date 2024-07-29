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
			slog.Error("error parsing home.html")
			return
		}
		err = t.Execute(w, rooms)
		if err != nil {
			slog.Error("error executing home.html template")
			return
		}
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
		t, err := template.ParseFiles("pages/room.html")
		if err != nil {
			slog.Error("error parsing room.html")
			return
		}
		err = t.Execute(w, room)
		if err != nil {
			slog.Error("error executing room.html template")
			return
		}
	})
}
