package main

import (
	"errors"
	"gossip/backend/pkg/config"
	"gossip/backend/pkg/database"
	"gossip/backend/pkg/routers"
	"log"

	"github.com/gofiber/fiber/v2"
)

func fiberError(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	if err = c.Status(code).JSON(e); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	return nil
}

func main() {
	var err error

	env, err := config.LoadEnv(".")
	if err != nil {
		panic("Cannot load env")
	}

	DB := database.Initialise(&env)

	app := fiber.New(fiber.Config{
		ErrorHandler: fiberError,
	})

	routers.Initialise(app, DB)

	log.Fatal(app.Listen(":" + env.ServerPort))
}
