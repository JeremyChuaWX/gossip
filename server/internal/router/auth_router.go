package router

import (
	"gossip/internal/constants"
	_user "gossip/internal/domains/user"
	"gossip/internal/password"
	"gossip/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

func (router *Router) authRouter() *chi.Mux {
	authRouter := chi.NewMux()

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

		dto := _user.FindOneByUsernameDTO{
			Username: req.Username,
		}
		user, err := router.UserRepository.FindOneByUsername(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if err = password.Verify(req.Password, user.PasswordHash); err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		sessionId, err := router.SessionRepository.Create(
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

		dto := _user.CreateDTO{
			Username:     req.Username,
			PasswordHash: passwordHash,
		}
		user, err := router.UserRepository.Create(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		sessionId, err := router.SessionRepository.Create(
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
				Message: "signed up",
			},
			SessionId: sessionId,
		})
	})

	authRouter.Group(func(authRouter chi.Router) {
		authRouter.Use(router.Middlewares.AuthMiddleware())

		// sign out
		authRouter.Post(
			"/signout",
			func(w http.ResponseWriter, r *http.Request) {
				sessionId := r.Header.Get(constants.SESSION_ID_HEADER)
				if err := router.SessionRepository.Delete(r.Context(), sessionId); err != nil {
					utils.WriteError(w, http.StatusUnauthorized, err)
					return
				}
			},
		)

		// me
		authRouter.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			authUserId := r.Context().Value(constants.USER_ID_CONTEXT_KEY).(uuid.UUID)

			dto := _user.FindOneDTO{
				Id: authUserId,
			}
			user, err := router.UserRepository.FindOne(r.Context(), dto)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, err)
				return
			}

			type response struct {
				utils.BaseResponse
				User _user.User `json:"user"`
			}
			utils.WriteJSON(w, http.StatusOK, response{
				BaseResponse: utils.BaseResponse{
					Error:   false,
					Message: "current signed in user",
				},
				User: user,
			})
		})
	})

	return authRouter
}
