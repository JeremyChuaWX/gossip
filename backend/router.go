package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Router struct {
	Mux *chi.Mux
}

func NewRouter(h *Handlers) *Router {
	mux := chi.NewRouter()

	// middleware
	mux.Use(cors.Handler(cors.Options{}))
	mux.Use(middleware.StripSlashes)
	mux.Use(middleware.CleanPath)
	mux.Use(middleware.Logger)

	// handlers
	mux.Post("/room", h.NewRoomHandler)
	mux.Get("/room", h.GetRoomsHandler)
	mux.HandleFunc("/room/{room}", h.JoinRoomHandler)

	return &Router{Mux: mux}
}
