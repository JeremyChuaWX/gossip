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

	cmtRouter.Post("/", middlewares.Jwtware(), cmtHandler.CreateComment)
	cmtRouter.Get("/:id", cmtHandler.GetCommentById)
	cmtRouter.Put("/:id", middlewares.Jwtware(), cmtHandler.UpdateComment)
	cmtRouter.Delete("/:id", middlewares.Jwtware(), cmtHandler.DeleteComment)
}
