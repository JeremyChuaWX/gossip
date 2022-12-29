package routers

import (
	"gossip/backend/pkg/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initUserRouter(router *gin.RouterGroup, DB *gorm.DB) {
	userHandler := handlers.UserHandler{DB: DB}
	userRouter := router.Group("/users")

	userRouter.GET("/:id", userHandler.GetUserById)
	userRouter.PUT("/:id", userHandler.UpdateUser)
	userRouter.DELETE("/:id", userHandler.DeleteUser)
}
