package routers

import (
	"gossip/backend/pkg/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initTagRouter(router *gin.RouterGroup, DB *gorm.DB) {
	tagHandler := handlers.TagHandler{DB: DB}
	tagRouter := router.Group("/tags")

	tagRouter.POST("/tags", tagHandler.CreateTag)
	tagRouter.GET("/tags", tagHandler.GetAllTags)
	tagRouter.GET("/tags/:id", tagHandler.GetTagById)
	tagRouter.PUT("/tags/:id", tagHandler.UpdateTag)
	tagRouter.DELETE("/tags/:id", tagHandler.DeleteTag)
}
