package routers

import (
	"gossip/backend/pkg/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func initTagRouter(router fiber.Router, DB *gorm.DB) {
	tagHandler := handlers.TagHandler{DB: DB}
	tagRouter := router.Group("/tags")

	tagRouter.Post("/create-tag", tagHandler.CreateTag)
	tagRouter.Get("/get-tags", tagHandler.GetAllTags)
	tagRouter.Get("/get-tag/:id", tagHandler.GetTagById)
	tagRouter.Put("/update-tag/:id", tagHandler.UpdateTag)
	tagRouter.Delete("/delete-tag/:id", tagHandler.DeleteTag)
}
