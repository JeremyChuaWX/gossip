package user

import (
	"gossip/internal/password"
	"gossip/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

type Service struct {
	Repository *Repository
}

func (s *Service) InitRoutes(router *chi.Mux) {
	userRouter := s.userRouter()
	router.Mount("/users", userRouter)
}

func (s *Service) userRouter() *chi.Mux {
	userRouter := chi.NewRouter()

	// find one user
	userRouter.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		urlId := chi.URLParam(r, "id")
		id, err := uuid.FromString(urlId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		dto := findOneDTO{
			id: id,
		}

		user, err := s.Repository.findOne(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			utils.BaseResponse
			User User `json:"user"`
		}
		utils.WriteJSON(w, http.StatusOK, response{
			BaseResponse: utils.BaseResponse{
				Error:   false,
				Message: "found user",
			},
			User: user,
		})
	})

	// find one user by username
	userRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Username string `query:"username"`
		}
		req, err := utils.GetURLQueryStruct[request](r.URL)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		dto := FindOneByUsernameDTO{
			Username: req.Username,
		}

		user, err := s.Repository.FindOneByUsername(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			utils.BaseResponse
			User User `json:"user"`
		}
		utils.WriteJSON(w, http.StatusOK, response{
			BaseResponse: utils.BaseResponse{
				Error:   false,
				Message: "found user",
			},
			User: user,
		})
	})

	// update user
	userRouter.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Username *string `json:"username,omitempty"`
			Password *string `json:"password,omitempty"`
		}
		req, err := utils.ReadJSON[Request](r)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		urlId := chi.URLParam(r, "id")
		id, err := uuid.FromString(urlId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
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
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}
		}

		user, err := s.Repository.update(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			utils.BaseResponse
			User User `json:"user"`
		}
		utils.WriteJSON(w, http.StatusOK, response{
			BaseResponse: utils.BaseResponse{
				Error:   false,
				Message: "updated user",
			},
			User: user,
		})
	})

	// delete user
	userRouter.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		urlId := chi.URLParam(r, "id")
		id, err := uuid.FromString(urlId)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		dto := deleteDTO{
			id: id,
		}
		user, err := s.Repository.delete(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			utils.BaseResponse
			User User `json:"user"`
		}
		utils.WriteJSON(w, http.StatusOK, response{
			BaseResponse: utils.BaseResponse{
				Error:   false,
				Message: "deleted user",
			},
			User: user,
		})
	})

	return userRouter
}
