package routers

import (
	"gossip/backend/pkg/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func initPostRouter(router fiber.Router, DB *gorm.DB) {
	postHandler := handlers.PostHandler{DB: DB}
	postRouter := router.Group("/posts")

	postRouter.Post("/", postHandler.CreatePost)
	postRouter.Get("/", postHandler.GetAllPosts)
	postRouter.Get("/:id", postHandler.GetPostById)
	postRouter.Put("/:id", postHandler.UpdatePost)
	postRouter.Delete("/:id", postHandler.DeletePost)
}
