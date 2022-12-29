package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(engine *gin.Engine, DB *gorm.DB) {
	router := engine.Group("/api")

	initAuthRouter(router, DB)
	initUserRouter(router, DB)
	initPostRouter(router, DB)
	initCommentRouter(router, DB)
	initTagRouter(router, DB)
}
