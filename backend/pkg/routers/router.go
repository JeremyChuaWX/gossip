package routers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Initialise(app *fiber.App, DB *gorm.DB) {
	api := app.Group("/api")

	initAuthRouter(api, DB)
	initUserRouter(api, DB)
	initPostRouter(api, DB)
	initCommentRouter(api, DB)
	initTagRouter(api, DB)
	initTaggableRouter(api, DB)
	initSubscriptionRouter(api, DB)
}
