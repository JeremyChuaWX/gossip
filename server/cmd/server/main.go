package main

import (
	"context"
	"gossip/internal/adapters/postgres"
	"gossip/internal/adapters/redis"
	"gossip/internal/config"
	"gossip/internal/domains/chat"
	"gossip/internal/domains/session"
	"gossip/internal/domains/user"
	"gossip/internal/middlewares"
	"gossip/internal/router"
	"log"
)

func main() {
	config, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgPool, err := postgres.Init(ctx, config.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pgPool.Close()

	redis, err := redis.Init(ctx, config.RedisURL, "")
	if err != nil {
		log.Fatal(err)
	}
	defer redis.Close()

	sessionRepository := &session.Repository{
		Redis: redis,
	}

	userRepository := &user.Repository{
		PgPool: pgPool,
	}

	chatService := chat.InitService()

	middlewares := &middlewares.Middlewares{
		SessionRepository: sessionRepository,
	}

	router := &router.Router{
		UserRepository:    userRepository,
		SessionRepository: sessionRepository,
		ChatService:       chatService,
		Middlewares:       middlewares,
	}

	log.Println("running server on address", config.ServerAddress)
	log.Fatal(router.Start(config.ServerAddress))
}
