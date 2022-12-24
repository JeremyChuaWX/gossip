package routes

func (s server) InitialiseUserRoutes() {
	s.router.GET("/users/:id", s.handler.GetUser)
	s.router.POST("/users", s.handler.CreateUser)
	s.router.PUT("/users/:id", s.handler.UpdateUser)
	s.router.DELETE("/users/:id", s.handler.DeleteUser)
}
