package routers

import (
	"gossip/backend/pkg/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initAuthRouter(router *gin.RouterGroup, DB *gorm.DB) {
	authHandler := handlers.AuthHandler{DB: DB}
	authRouter := router.Group("/auth")

	authRouter.POST("/signup", authHandler.SignUp)
	authRouter.POST("/signin", authHandler.SignIn)
	authRouter.POST("/signout", authHandler.SignOut)
}
