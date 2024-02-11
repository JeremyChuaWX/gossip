package user

import "github.com/go-chi/chi/v5"

func InitRoutes(r *chi.Mux, s *Service) *chi.Mux {
	ur := chi.NewRouter()
	ur.Post("/", s.createHandler)
	ur.Get("/{id}", s.findOneHandler)
	ur.Patch("/{id}", s.updateHandler)
	ur.Delete("/{id}", s.deleteHandler)
	r.Mount("/users", ur)
	return r
}
