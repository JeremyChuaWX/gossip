package main

import (
	"gossip/backend/pkg/initialisers"
	"gossip/backend/pkg/routers"

	"github.com/gin-gonic/gin"
)

var server *gin.Engine

func init() {
	config, err := initialisers.LoadConfig(".")
	if err != nil {
		panic("Cannot load environment variables")
	}

	DB := initialisers.ConnectDB(&config)
	initialisers.MigrateDB(DB)

	server = gin.Default()

	routers.Init(server, DB)
}

func main() {
	var err error

	config, err := initialisers.LoadConfig(".")
	if err != nil {
		panic("Cannot load environment variables")
	}

	if err = server.Run(":" + config.ServerPort); err != nil {
		panic("Cannot run server on port " + config.ServerPort)
	}
}
