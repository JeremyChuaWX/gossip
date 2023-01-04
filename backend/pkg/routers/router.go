package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func Initialise(app *fiber.App, DB *gorm.DB) {
	api := app.Group("/api")

	api.Use(logger.New())
	api.Use(cors.New(cors.Config{
		AllowOrigins: "localhost:3000",
	}))

	initAuthRouter(api, DB)
	initUserRouter(api, DB)
	initPostRouter(api, DB)
	initCommentRouter(api, DB)
	initTagRouter(api, DB)
	initTaggableRouter(api, DB)
	initSubscriptionRouter(api, DB)
}
