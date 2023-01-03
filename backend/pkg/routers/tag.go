package routers

import (
	"gossip/backend/pkg/handlers"
	"gossip/backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func initTagRouter(router fiber.Router, DB *gorm.DB) {
	tagHandler := handlers.TagHandler{DB: DB}
	tagRouter := router.Group("/tags")

	tagRouter.Post("/", tagHandler.CreateTag)
	tagRouter.Get("/", tagHandler.GetAllTags)
	tagRouter.Get("/:id", tagHandler.GetTagById)
	tagRouter.Put("/:id", tagHandler.UpdateTag)
	tagRouter.Delete("/:id", tagHandler.DeleteTag)
	tagRouter.Post("/tagpost", middlewares.Jwtware(), tagHandler.TagPost)
}
