package api

import (
	"gossip/backend/pkg/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initCommentRouter(router *gin.RouterGroup, DB *gorm.DB) {
	cmtHandler := handlers.CommentHandler{DB: DB}

	router.POST("/comments", cmtHandler.CreateComment)
	router.GET("/comments/:id", cmtHandler.GetCommentById)
	router.PUT("/comments/:id", cmtHandler.UpdateComment)
	router.DELETE("/comments/:id", cmtHandler.DeleteComment)
}
