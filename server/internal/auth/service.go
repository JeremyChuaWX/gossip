package auth

import (
	"gossip/internal/password"
	"gossip/internal/user"
	"gossip/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Service struct {
	UserRepository    *user.Repository
	SessionRepository *Repository
}

func (s *Service) InitRoutes(router *chi.Mux) {
	authRouter := s.authRouter()
	router.Mount("/auth", authRouter)
}

func (s *Service) authRouter() *chi.Mux {
	authRouter := chi.NewRouter()

	// sign in
	authRouter.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		req, err := utils.ReadJSON[request](r)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		dto := user.FindOneByUsernameDTO{
			Username: req.Username,
		}
		user, err := s.UserRepository.FindOneByUsername(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if err = password.Compare([]byte(user.PasswordHash), req.Password); err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		sessionId, err := s.SessionRepository.createSession(
			r.Context(),
			user.Id.String(),
		)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			utils.BaseResponse
			SessionId string `json:"sessionId"`
		}
		utils.WriteJSON(w, http.StatusOK, response{
			BaseResponse: utils.BaseResponse{
				Error:   false,
				Message: "signed in",
			},
			SessionId: sessionId,
		})
	})

	// sign up
	authRouter.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		type response struct {
			utils.BaseResponse
			User user.User `json:"user"`
		}

		req, err := utils.ReadJSON[request](r)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		passwordHash, err := password.Hash(req.Password)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		dto := user.CreateDTO{
			Username:     req.Username,
			PasswordHash: passwordHash,
		}
		user, err := s.UserRepository.Create(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, response{
			BaseResponse: utils.BaseResponse{
				Error:   false,
				Message: "signed up",
			},
			User: user,
		})
	})

	return authRouter
}
