package main

import (
	"errors"
	"gossip/backend/pkg/initialisers"
	"gossip/backend/pkg/routers"
	"log"

	"github.com/gofiber/fiber/v2"
)

var app *fiber.App

func init() {
	config, err := initialisers.LoadConfig(".")
	if err != nil {
		panic("Cannot load environment variables")
	}

	DB := initialisers.ConnectDB(&config)
	initialisers.MigrateDB(DB)

	app = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			if err = c.Status(code).JSON(e); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
			}

			return nil
		},
	})

	routers.Init(app, DB)
}

func main() {
	var err error

	config, err := initialisers.LoadConfig(".")
	if err != nil {
		panic("Cannot load environment variables")
	}

	log.Fatal(app.Listen(":" + config.ServerPort))
}
