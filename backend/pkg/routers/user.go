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

	userRouter.Get("/get-user/:id", userHandler.GetUserById)
	userRouter.Get("/get-me", middlewares.Jwtware(), userHandler.GetMe)
	userRouter.Put("/update-me", middlewares.Jwtware(), userHandler.UpdateMe)
	userRouter.Delete("/delete-me", middlewares.Jwtware(), userHandler.DeleteMe)
	userRouter.Get("/toggle-visibility", middlewares.Jwtware(), userHandler.ToggleProfileVisibility)
}
