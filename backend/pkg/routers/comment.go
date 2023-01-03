package routers

import (
	"gossip/backend/pkg/handlers"
	"gossip/backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func initCommentRouter(router fiber.Router, DB *gorm.DB) {
	cmtHandler := handlers.CommentHandler{DB: DB}
	cmtRouter := router.Group("/comments")

	cmtRouter.Post("/create-post", middlewares.Jwtware(), cmtHandler.CreateComment)
	cmtRouter.Get("/get-comment/:id", cmtHandler.GetCommentById)
	cmtRouter.Put("/update-comment/:id", middlewares.Jwtware(), cmtHandler.UpdateComment)
	cmtRouter.Delete("/delete-comment/:id", middlewares.Jwtware(), cmtHandler.DeleteComment)
}
