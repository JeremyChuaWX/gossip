package user

import (
	"encoding/json"
	"gossip/internal/password"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

type errorResponse struct {
	Error string `json:"error"`
}

type Service struct {
	Repository *Repository
}

func (s *Service) InitRoutes(router *chi.Mux) {
	userRouter := s.userRouter()
	router.Mount("/users", userRouter)
}

func (s *Service) userRouter() *chi.Mux {
	userRouter := chi.NewRouter()

	// create user
	userRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		req := readJSON[request](r)

		passwordHash, err := password.Hash(req.Password)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		dto := createDTO{
			username:     req.Username,
			passwordHash: passwordHash,
		}
		user, err := s.Repository.create(r.Context(), dto)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			Message string `json:"message"`
			User    User
		}
		writeJSON(w, http.StatusOK, response{
			Message: "created user",
			User:    user,
		})
	})

	// find one user
	userRouter.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		urlId := chi.URLParam(r, "id")
		id, err := uuid.FromString(urlId)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		dto := findOneDTO{
			id: id,
		}

		user, err := s.Repository.findOne(r.Context(), dto)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			Message string `json:"message"`
			User    User
		}
		writeJSON(w, http.StatusOK, response{
			Message: "created user",
			User:    user,
		})
	})

	// update user
	userRouter.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Username *string `json:"username,omitempty"`
			Password *string `json:"password,omitempty"`
		}
		req := readJSON[Request](r)

		urlId := chi.URLParam(r, "id")
		id, err := uuid.FromString(urlId)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		dto := updateDTO{
			id:           id,
			username:     req.Username,
			passwordHash: nil,
		}

		if req.Password != nil {
			passwordHash, err := password.Hash(*req.Password)
			dto.passwordHash = passwordHash
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
		}

		user, err := s.Repository.update(r.Context(), dto)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			Message string `json:"message"`
			User    User
		}
		writeJSON(w, http.StatusOK, response{
			Message: "created user",
			User:    user,
		})
	})

	// delete user
	userRouter.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		urlId := chi.URLParam(r, "id")
		id, err := uuid.FromString(urlId)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		dto := deleteDTO{
			id: id,
		}
		user, err := s.Repository.delete(r.Context(), dto)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			Message string `json:"message"`
			User    User
		}
		writeJSON(w, http.StatusOK, response{
			Message: "created user",
			User:    user,
		})
	})

	return userRouter
}

func readJSON[T any](r *http.Request) T {
	var res T
	json.NewDecoder(r.Body).Decode(&res)
	return res
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		json.NewEncoder(w).Encode(errorResponse{
			Error: err.Error(),
		})
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, errorResponse{
		Error: err.Error(),
	})
}
