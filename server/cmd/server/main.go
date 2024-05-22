package main

import (
	"context"
	"gossip/internal/adapters/postgres"
	"gossip/internal/adapters/redis"
	"gossip/internal/api/middlewares"
	"gossip/internal/api/routes"
	"gossip/internal/config"
	"gossip/internal/services/room"
	"gossip/internal/services/roomuser"
	"gossip/internal/services/session"
	"gossip/internal/services/user"
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

	sessionService := &session.Service{
		Redis: redis,
	}

	userService := &user.Service{
		PgPool: pgPool,
	}

	roomService := &room.Service{
		PgPool: pgPool,
	}

	roomUserService := &roomuser.Service{
		PgPool: pgPool,
	}

	// chatService := chat.InitService(userRepository, room)

	middlewares := &middlewares.Middlewares{
		SessionService: sessionService,
	}

	router := &routes.Router{
		RoomService:     roomService,
		UserService:     userService,
		RoomUserService: roomUserService,
		SessionService:  sessionService,
		Middlewares:     middlewares,
	}

	log.Println("running server on address", config.ServerAddress)
	log.Fatal(router.Start(config.ServerAddress))
}
