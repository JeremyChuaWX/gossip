package routers

import (
	"gossip/backend/pkg/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initCommentRouter(router *gin.RouterGroup, DB *gorm.DB) {
	cmtHandler := handlers.CommentHandler{DB: DB}
	cmtRouter := router.Group("/comments")

	cmtRouter.POST("", cmtHandler.CreateComment)
	cmtRouter.GET("/:id", cmtHandler.GetCommentById)
	cmtRouter.PUT("/:id", cmtHandler.UpdateComment)
	cmtRouter.DELETE("/:id", cmtHandler.DeleteComment)
}
