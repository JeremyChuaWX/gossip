package main

import (
	"context"
	"gossip/internal/adapters/postgres"
	"gossip/internal/user"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

const ADDRESS string = "127.0.0.1:3000"

func main() {
	// constants
	_ = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	router := chi.NewMux()
	ctx := context.Background()

	// adapters
	pgPool, err := postgres.Init(ctx, "")
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
