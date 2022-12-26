package main

import (
	"gossip/backend/internal/db"
	"gossip/backend/internal/handlers"
	"gossip/backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	DB := db.Init()
	router := gin.Default()
	handler := handlers.New(DB)

	routes.InitialiseRoutes(router, handler)

	router.Run("localhost:3000")
}
