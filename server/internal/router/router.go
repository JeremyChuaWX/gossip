package router

import (
	"gossip/internal/chat"
	"gossip/internal/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	Repository  *repository.Repository
	ChatService *chat.Service
}

func (router *Router) Init() (*chi.Mux, error) {
	mux := chi.NewMux()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Group(router.routeGroup)
	mux.Group(router.authedRouteGroup)

	if err := chi.Walk(mux, walkRoutes); err != nil {
		return nil, err
	}

	mux.Handle(
		"/static/*",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
	)

	return mux, nil
}
