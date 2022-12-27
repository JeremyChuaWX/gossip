package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(engine *gin.Engine, DB *gorm.DB) {
	apiGroup := engine.Group("/api")

	initUserRouter(apiGroup, DB)
	initPostRouter(apiGroup, DB)
	initCommentRouter(apiGroup, DB)
}
