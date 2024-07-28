package main

import (
	"context"
	"gossip/internal/adapters/postgres"
	"gossip/internal/api/middlewares"
	"gossip/internal/api/routes"
	"gossip/internal/chat"
	"gossip/internal/config"
	"gossip/internal/repository"
	"log/slog"
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

	middlewares := &middlewares.Middlewares{
		Repository: repository,
	}

	router := &routes.Router{
		Repository:  repository,
		ChatService: chatService,
		Middlewares: middlewares,
	}

	slog.Info("server is running", "address", config.ServerAddress)
	err = router.Start(config.ServerAddress)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
