package api

import (
	"gossip/backend/pkg/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initUserRouter(router *gin.RouterGroup, DB *gorm.DB) {
	userHandler := handlers.UserHandler{DB: DB}

	router.POST("/users", userHandler.SignUp)
	router.GET("/users/:id", userHandler.GetUserById)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)
}
