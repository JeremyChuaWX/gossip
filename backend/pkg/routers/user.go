package routers

import (
	"gossip/backend/pkg/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func initUserRouter(router fiber.Router, DB *gorm.DB) {
	userHandler := handlers.UserHandler{DB: DB}
	userRouter := router.Group("/users")

	userRouter.Get("/:id", userHandler.GetUserById)
	userRouter.Put("/:id", userHandler.UpdateUser)
	userRouter.Delete("/:id", userHandler.DeleteUser)
}
