package api

import (
	"gossip/internal/chat"
	"gossip/internal/repository"

	"github.com/go-chi/chi/v5"
)

type Api struct {
	Repository  *repository.Repository
	ChatService *chat.Service
}

func (api *Api) InitRouter() *chi.Mux {
	router := chi.NewMux()
	router.Group(api.routesGroup)
	router.Group(api.authedRoutesGroup)
	return router
}
