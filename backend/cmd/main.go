package main

import (
	"gossip/backend/pkg/api"
	"gossip/backend/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	DB := db.Init()
	engine := gin.Default()

	api.Init(engine, DB)

	engine.Run("localhost:3000")
}
