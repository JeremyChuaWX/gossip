package main

import (
	"gossip/backend/pkg/config"
	"gossip/backend/pkg/database"
	"gossip/backend/pkg/routers"
	"gossip/backend/pkg/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	var err error

	env, err := config.LoadConfig(".")
	if err != nil {
		panic("Cannot load environment variables")
	}

	DB := database.Initialise(&env)

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.CustomErrorHandler,
	})

	routers.Init(app, DB)

	log.Fatal(app.Listen(":" + env.ServerPort))
}
