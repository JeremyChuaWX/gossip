package routers

import (
	"gossip/backend/pkg/handlers"
	"gossip/backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func initSubscriptionRouter(router fiber.Router, DB *gorm.DB) {
	subscriptionHandler := handlers.SubscriptionHandler{DB: DB}
	subscriptionRouter := router.Group("/subscriptions")

	subscriptionRouter.Use(middlewares.Jwtware())

	subscriptionRouter.Post("/create-subscription", subscriptionHandler.CreateSubscription)
	subscriptionRouter.Delete("/delete-subscription", subscriptionHandler.DeleteSubscription)
}
