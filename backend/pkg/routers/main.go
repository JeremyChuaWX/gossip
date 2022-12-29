package routers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Init(app *fiber.App, DB *gorm.DB) {
	api := app.Group("/api")

	initAuthRouter(api, DB)
	initUserRouter(api, DB)
	initPostRouter(api, DB)
	initCommentRouter(api, DB)
	initTagRouter(api, DB)
}
