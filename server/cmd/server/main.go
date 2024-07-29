package main

import (
	"context"
	"gossip/internal/adapters/postgres"
	"gossip/internal/api"
	"gossip/internal/chat"
	"gossip/internal/config"
	"gossip/internal/repository"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config, err := config.Init()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgPool, err := postgres.Init(ctx, config.PostgresURL)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer pgPool.Close()

	repository := &repository.Repository{
		PgPool: pgPool,
	}

	chatService, err := chat.NewService(repository)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	api := &api.Api{
		Repository:  repository,
		ChatService: chatService,
	}

	// web := &web.Web{}

	router := chi.NewMux()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/api", api.InitRouter())
	// router.Mount("/", web.InitRouter())

	if err := chi.Walk(router, walkRoutes); err != nil {
		slog.Error("error walking router")
		return
	}

	slog.Info("server is running", "address", config.ServerAddress)
	if err := http.ListenAndServe(config.ServerAddress, router); err != nil {
		slog.Error(err.Error())
		return
	}
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
