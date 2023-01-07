package main

import (
	"errors"
	"gossip/backend/pkg/config"
	"gossip/backend/pkg/database"
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/routers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func fiberError(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	if err = c.Status(code).JSON(models.ServerResponse{
		Error: true,
		Msg:   e.Message,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	return nil
}

func main() {
	var err error

	env, err := config.LoadEnv(".")
	if err != nil {
		log.Fatal(err)
	}

	DB := database.Initialise(&env)

	app := fiber.New(fiber.Config{
		ErrorHandler: fiberError,
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	routers.Initialise(app, DB)

	log.Fatal(app.Listen(":" + env.ServerPort))
}
