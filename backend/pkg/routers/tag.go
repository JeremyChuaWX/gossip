package routers

import (
	"gossip/backend/pkg/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func initTagRouter(router fiber.Router, DB *gorm.DB) {
	tagHandler := handlers.TagHandler{DB: DB}
	tagRouter := router.Group("/tags")

	tagRouter.Post("/tags", tagHandler.CreateTag)
	tagRouter.Get("/tags", tagHandler.GetAllTags)
	tagRouter.Get("/tags/:id", tagHandler.GetTagById)
	tagRouter.Put("/tags/:id", tagHandler.UpdateTag)
	tagRouter.Delete("/tags/:id", tagHandler.DeleteTag)
}
