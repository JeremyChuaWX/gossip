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
	userRouter.Put("/update-me", userHandler.UpdateMe)
	userRouter.Delete("/delete-me", userHandler.DeleteMe)
	userRouter.Put("/toggle-visibility", userHandler.ToggleProfileVisibility)
}
