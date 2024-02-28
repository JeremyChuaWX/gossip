package main

import (
	"context"
	"gossip/internal/adapters/postgres"
	"gossip/internal/chat"
	"gossip/internal/user"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const ADDRESS string = "server:3000"

func main() {
	// constants
	router := InitRouter()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// adapters
	pgPool, err := postgres.Init(
		ctx,
		"postgresql://admin:password123@postgres:5432/my_db?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer pgPool.Close()

	// user module
	userRepository := &user.Repository{
		PgPool: pgPool,
	}
	userService := &user.Service{
		Repository: userRepository,
	}
	userService.InitRoutes(router)

	// chat module
	chatService := chat.InitService()
	chatService.InitRoutes(router)

	// run server
	StartRouter(router)
}

func InitRouter() *chi.Mux {
	router := chi.NewMux()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	return router
}

func StartRouter(router *chi.Mux) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		log.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Printf("Logging err: %s\n", err.Error())
	}

	log.Println("running server on address", ADDRESS)
	log.Fatal(http.ListenAndServe(ADDRESS, router))
}