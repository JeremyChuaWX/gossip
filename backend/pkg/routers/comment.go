package routers

import (
	"gossip/backend/pkg/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func initCommentRouter(router fiber.Router, DB *gorm.DB) {
	cmtHandler := handlers.CommentHandler{DB: DB}
	cmtRouter := router.Group("/comments")

	cmtRouter.Post("", cmtHandler.CreateComment)
	cmtRouter.Get("/:id", cmtHandler.GetCommentById)
	cmtRouter.Put("/:id", cmtHandler.UpdateComment)
	cmtRouter.Delete("/:id", cmtHandler.DeleteComment)
}
