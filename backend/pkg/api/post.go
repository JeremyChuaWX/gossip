package api

import (
	"gossip/backend/pkg/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initPostRouter(router *gin.RouterGroup, DB *gorm.DB) {
	postHandler := handlers.PostHandler{DB: DB}

	router.POST("/posts", postHandler.CreatePost)
	router.GET("/posts", postHandler.GetAllPosts)
	router.GET("/posts/:id", postHandler.GetPost)
	router.PUT("/posts/:id", postHandler.UpdatePost)
	router.DELETE("/posts/:id", postHandler.DeletePost)
}
