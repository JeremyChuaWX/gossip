package routers

import (
	"gossip/backend/pkg/handlers"
	"gossip/backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func initTaggableRouter(router fiber.Router, DB *gorm.DB) {
	taggableHandler := handlers.TaggableHandler{DB: DB}
	taggableRouter := router.Group("/taggable")

	taggableRouter.Use(middlewares.Jwtware())

	taggableRouter.Post("/create-taggable", taggableHandler.CreateTaggable)
}
