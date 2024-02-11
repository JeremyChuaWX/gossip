package user

import (
	"encoding/json"
	"gossip/internal/password"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

type Service struct {
	Repository *Repository
}

func (s *Service) createHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	req := readJSON[Request](r)

	password, err := password.Hash(req.Password)
	if err != nil {
		// writeError
		return
	}

	dto := createDTO{
		Username: req.Username,
		Password: password,
	}
	user, err := s.Repository.create(r.Context(), dto)
	if err != nil {
		// writeError
		return
	}

	// writeResponse
}

func (s *Service) findOneHandler(w http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "id")
	id, err := uuid.FromString(urlId)
	if err != nil {
		// writeError
		return
	}

	dto := findOneDTO{
		Id: id,
	}

	user, err := s.Repository.findOne(r.Context(), dto)
	if err != nil {
		// writeError
		return
	}

	// writeResponse
}

func (s *Service) updateHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Username *string `json:"username,omitempty"`
		Password *string `json:"password,omitempty"`
	}
	req := readJSON[Request](r)

	urlId := chi.URLParam(r, "id")
	id, err := uuid.FromString(urlId)
	if err != nil {
		// writeError
		return
	}

	dto := updateDTO{
		Id:       id,
		Username: req.Username,
		Password: nil,
	}

	if req.Password != nil {
		password, err := password.Hash(*req.Password)
		dto.Password = password
		if err != nil {
			// writeError
			return
		}
	}

	user, err := s.Repository.update(r.Context(), dto)
	if err != nil {
		// writeError
		return
	}

	// writeResponse
}

func (s *Service) deleteHandler(w http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "id")
	id, err := uuid.FromString(urlId)
	if err != nil {
		// writeError
		return
	}
	dto := deleteDTO{
		Id: id,
	}
	user, err := s.Repository.delete(r.Context(), dto)
	if err != nil {
		// writeError
		return
	}
	// writeResponse
}

func readJSON[T any](r *http.Request) T {
	var res T
	json.NewDecoder(r.Body).Decode(&res)
	return res
}

func writeJSON(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode("")
}
