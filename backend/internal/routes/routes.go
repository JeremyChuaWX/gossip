package routes

import (
	"gossip/backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

type server struct {
	router  *gin.Engine
	handler handlers.Handler
}

func InitialiseRoutes(router *gin.Engine, handler handlers.Handler) {
	s := server{router, handler}

	s.InitialiseUserRoutes()
}
