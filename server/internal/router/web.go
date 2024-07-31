package router

import (
	"gossip/internal/repository"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

func (router *Router) webRouter() *chi.Mux {
	web := chi.NewMux()
	web.Group(router.webRouteGroup)
	web.Group(router.webAuthedRouteGroup)
	return web
}

func (router *Router) webRouteGroup(mux chi.Router) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := sessionFromContext(r.Context())
		if err != nil {
			http.ServeFile(w, r, "pages/index.html")
			return
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	})

	mux.Get("/signup", func(w http.ResponseWriter, r *http.Request) {
		_, err := sessionFromContext(r.Context())
		if err != nil {
			http.ServeFile(w, r, "pages/signup.html")
			return
		}
		prev := r.URL.Query().Get("prev")
		if prev != "" {
			http.Redirect(w, r, prev, http.StatusFound)
			return
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	})

	mux.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		_, err := sessionFromContext(r.Context())
		if err != nil {
			http.ServeFile(w, r, "pages/login.html")
			return
		}
		prev := r.URL.Query().Get("prev")
		if prev != "" {
			http.Redirect(w, r, prev, http.StatusFound)
			return
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	})
}

func (router *Router) webAuthedRouteGroup(mux chi.Router) {
	mux.Use(router.authMiddleware)

	mux.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		session := sessionFromContextSafe(r.Context())
		rooms, err := router.Repository.RoomFindManyByUserId(
			r.Context(),
			repository.RoomFindManyByUserIdParams{UserId: session.UserId},
		)
		if err != nil {
			slog.Error(
				"error finding rooms for user",
				"userSession",
				session,
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

	mux.Get("/rooms/join", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/join-room.html")
	})

	mux.Get("/rooms/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		session := sessionFromContextSafe(r.Context())
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
				UserId: session.UserId,
				RoomId: roomId,
			},
		)
		if err != nil || !isMember {
			slog.Error(
				"user not in room",
				"userId",
				session.UserId,
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
