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

	postRouter.Post("/create-post", middlewares.Jwtware(), postHandler.CreatePost)
	postRouter.Get("/get-posts", postHandler.GetAllPosts)
	postRouter.Get("/get-post/:id", postHandler.GetPostById)
	postRouter.Put("/update-post/:id", middlewares.Jwtware(), postHandler.UpdatePost)
	postRouter.Delete("/delete-post/:id", middlewares.Jwtware(), postHandler.DeletePost)
}
