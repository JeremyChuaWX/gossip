package api

import (
	"gossip/backend/pkg/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initTagRouter(router *gin.RouterGroup, DB *gorm.DB) {
	tagHandler := handlers.TagHandler{DB: DB}

	router.POST("/tags", tagHandler.CreateTag)
	router.GET("/tags", tagHandler.GetAllTags)
	router.GET("/tags/:id", tagHandler.GetTagById)
	router.PUT("/tags/:id", tagHandler.UpdateTag)
	router.DELETE("/tags/:id", tagHandler.DeleteTag)
}
