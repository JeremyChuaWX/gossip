package main

import (
	"context"
	"gossip/internal/adapters/postgres"
	"gossip/internal/adapters/redis"
	"gossip/internal/chat"
	"gossip/internal/environment"
	"gossip/internal/middlewares"
	"gossip/internal/session"
	"gossip/internal/user"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	env, err := environment.Init()
	if err != nil {
		log.Fatal(err)
	}

	router := initRouter()
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

	middlewares := &middlewares.Middlewares{
		SessionRepository: sessionRepository,
	}

	userRepository := &user.Repository{
		PgPool: pgPool,
	}
	userService := &user.Service{
		UserRepository:    userRepository,
		SessionRepository: sessionRepository,
		Middlewares:       middlewares,
	}
	userService.InitRoutes(router)

	chatService := chat.InitService()
	chatService.InitRoutes(router)

	startRouter(router, env.ServerAddress)
}

func initRouter() *chi.Mux {
	router := chi.NewMux()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	return router
}

func startRouter(router *chi.Mux, address string) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		log.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Printf("Logging err: %s\n", err.Error())
	}

	log.Println("running server on address", address)
	log.Fatal(http.ListenAndServe(address, router))
}
