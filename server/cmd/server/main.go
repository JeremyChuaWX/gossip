package main

import (
	"context"
	"gossip/internal/adapters/postgres"
	"gossip/internal/adapters/redis"
	"gossip/internal/domains/chat"
	"gossip/internal/domains/session"
	"gossip/internal/domains/user"
	"gossip/internal/environment"
	"gossip/internal/middlewares"
	"gossip/internal/router"
	"log"
)

func main() {
	env, err := environment.Init()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgPool, err := postgres.Init(ctx, env.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pgPool.Close()

	redis, err := redis.Init(ctx, env.RedisURL, "")
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

	log.Println("running server on address", env.ServerAddress)
	log.Fatal(router.Start(env.ServerAddress))
}
