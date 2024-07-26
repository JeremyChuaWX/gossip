package main

import (
	"context"
	"gossip/internal/adapters/postgres"
	"gossip/internal/adapters/redis"
	"gossip/internal/api/middlewares"
	"gossip/internal/api/routes"
	"gossip/internal/config"
	"gossip/internal/services/chat"
	"gossip/internal/services/message"
	"gossip/internal/services/room"
	"gossip/internal/services/roomuser"
	"gossip/internal/services/session"
	"gossip/internal/services/user"
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

	redis, err := redis.Init(ctx, config.RedisURL, "")
	if err != nil {
		slog.Error(err.Error())
		return
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

	messageService := &message.Service{
		PgPool: pgPool,
	}

	chatService, err := chat.NewService(
		userService,
		roomService,
		roomUserService,
		messageService,
	)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	middlewares := &middlewares.Middlewares{
		SessionService: sessionService,
	}

	router := &routes.Router{
		RoomService:     roomService,
		UserService:     userService,
		RoomUserService: roomUserService,
		SessionService:  sessionService,
		ChatService:     chatService,
		Middlewares:     middlewares,
	}

	slog.Info("server is running", "address", config.ServerAddress)
	err = router.Start(config.ServerAddress)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
