package routers

import (
	"gossip/backend/pkg/handlers"
	"gossip/backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func initPostRouter(router fiber.Router, DB *gorm.DB) {
	postHandler := handlers.PostHandler{DB: DB}
	postRouter := router.Group("/posts")

	postRouter.Post("/", middlewares.Jwtware(), postHandler.CreatePost)
	postRouter.Get("/", postHandler.GetAllPosts)
	postRouter.Get("/:id", postHandler.GetPostById)
	postRouter.Put("/:id", middlewares.Jwtware(), postHandler.UpdatePost)
	postRouter.Delete("/:id", middlewares.Jwtware(), postHandler.DeletePost)
}
