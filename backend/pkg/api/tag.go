package api

import (
	"gossip/backend/pkg/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initTagRouter(router *gin.RouterGroup, DB *gorm.DB) {
	tagHandler := handlers.TagHandler{DB: DB}

	router.POST("/posts", tagHandler.CreateTag)
	router.GET("/posts", tagHandler.GetAllTags)
	router.GET("/posts/:id", tagHandler.GetTagById)
	router.PUT("/posts/:id", tagHandler.UpdateTag)
	router.DELETE("/posts/:id", tagHandler.DeleteTag)
}
