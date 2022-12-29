package routers

import (
	"gossip/backend/pkg/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initPostRouter(router *gin.RouterGroup, DB *gorm.DB) {
	postHandler := handlers.PostHandler{DB: DB}
	postRouter := router.Group("/posts")

	postRouter.POST("/posts", postHandler.CreatePost)
	postRouter.GET("/posts", postHandler.GetAllPosts)
	postRouter.GET("/posts/:id", postHandler.GetPostById)
	postRouter.PUT("/posts/:id", postHandler.UpdatePost)
	postRouter.DELETE("/posts/:id", postHandler.DeletePost)
}
