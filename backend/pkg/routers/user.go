package routers

import (
	"gossip/backend/pkg/handlers"
	"gossip/backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func initUserRouter(router fiber.Router, DB *gorm.DB) {
	userHandler := handlers.UserHandler{DB: DB}
	userRouter := router.Group("/users")

	userRouter.Use(middlewares.Jwtware())

	userRouter.Get("/get-user/:id", userHandler.GetUserById)
	userRouter.Get("/get-me", userHandler.GetMe)
	userRouter.Put("/update-user/:id", userHandler.UpdateUser)
	userRouter.Delete("/delete-user/:id", userHandler.DeleteUser)
	userRouter.Put("/toggle-visibility/:id", userHandler.ToggleProfileVisibility)
}
