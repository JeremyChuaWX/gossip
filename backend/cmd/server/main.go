package main

import (
	"context"
	"gossip/internal/adapters/postgres"
	"gossip/internal/user"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

const ADDRESS string = "server:3000"

func main() {
	// constants
	_ = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	router := InitRouter()
	ctx := context.Background()

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

	// run server
	log.Println("running server on address", ADDRESS)
	log.Fatal(http.ListenAndServe(ADDRESS, router))
}

func InitRouter() *chi.Mux {
	router := chi.NewMux()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	return router
}
