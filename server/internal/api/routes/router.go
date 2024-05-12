package router

import (
	"gossip/internal/domains/chat"
	"gossip/internal/domains/room"
	"gossip/internal/domains/roomuser"
	"gossip/internal/domains/session"
	"gossip/internal/domains/user"
	"gossip/internal/middlewares"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	RoomRepository     *room.Repository
	UserRepository     *user.Repository
	RoomUserRepository *roomuser.Repository
	SessionRepository  *session.Repository
	ChatService        *chat.Service
	Middlewares        *middlewares.Middlewares
}

func (r *Router) Start(address string) error {
	router := chi.NewMux()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/auth", r.authRouter())
	router.Mount("/users", r.userRouter())
	router.Mount("/chat", r.chatRouter())

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
