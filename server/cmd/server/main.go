package main

import (
	"context"
	"gossip/internal/adapters/postgres"
	"gossip/internal/chat"
	"gossip/internal/config"
	"gossip/internal/repository"
	"gossip/internal/router"
	"log"
	"log/slog"
	"net/http"
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
		log.Fatal(err.Error())
	}
	defer pgPool.Close()

	repository := &repository.Repository{
		PgPool: pgPool,
	}

	chatService, err := chat.NewService(repository)
	if err != nil {
		log.Fatal(err.Error())
	}

	router, err := (&router.Router{
		Repository:  repository,
		ChatService: chatService,
	}).Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	slog.Info("server is running", "address", config.ServerAddress)
	if err := http.ListenAndServe(config.ServerAddress, router); err != nil {
		log.Fatal(err.Error())
	}
}
