package routers

import (
	"gossip/backend/pkg/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func initPostRouter(router fiber.Router, DB *gorm.DB) {
	postHandler := handlers.PostHandler{DB: DB}
	postRouter := router.Group("/posts")

	postRouter.Post("/posts", postHandler.CreatePost)
	postRouter.Get("/posts", postHandler.GetAllPosts)
	postRouter.Get("/posts/:id", postHandler.GetPostById)
	postRouter.Put("/posts/:id", postHandler.UpdatePost)
	postRouter.Delete("/posts/:id", postHandler.DeletePost)
}
