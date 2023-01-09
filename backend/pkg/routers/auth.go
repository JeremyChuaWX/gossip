package routers

import (
	"gossip/backend/pkg/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func initAuthRouter(router fiber.Router, DB *gorm.DB) {
	authHandler := handlers.AuthHandler{DB: DB}
	authRouter := router.Group("/auth")

	authRouter.Post("/signup", authHandler.SignUp)
	authRouter.Post("/signin", authHandler.SignIn)
	authRouter.Get("/signout", authHandler.SignOut)
	authRouter.Get("/refresh", authHandler.RefreshAccessToken)
}
