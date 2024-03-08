package user

import (
	"errors"
	"gossip/internal/constants"
	"gossip/internal/middlewares"
	"gossip/internal/password"
	"gossip/internal/session"
	"gossip/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

var userForbiddenError = errors.New("mismatch user ID")

type Service struct {
	UserRepository    *Repository
	SessionRepository *session.Repository
	Middlewares       *middlewares.Middlewares
}

func (s *Service) InitRoutes(router *chi.Mux) {
	authRouter := s.authRouter()
	router.Mount("/auth", authRouter)

	userRouter := s.userRouter()
	router.Mount("/users", userRouter)
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

		dto := userFindOneByUsernameDTO{
			username: req.Username,
		}
		user, err := s.UserRepository.userFindOneByUsername(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if err = password.Verify(req.Password, user.PasswordHash); err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		sessionId, err := s.SessionRepository.Create(
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

		dto := userCreateDTO{
			username:     req.Username,
			passwordHash: passwordHash,
		}
		user, err := s.UserRepository.userCreate(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		sessionId, err := s.SessionRepository.Create(
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
		authRouter.Use(s.Middlewares.AuthMiddleware())

		// sign out
		authRouter.Post(
			"/signout",
			func(w http.ResponseWriter, r *http.Request) {
				sessionId := r.Header.Get(constants.SESSION_ID_HEADER)
				if err := s.SessionRepository.Delete(r.Context(), sessionId); err != nil {
					utils.WriteError(w, http.StatusUnauthorized, err)
					return
				}
			},
		)

		// me
		authRouter.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			authUserId := r.Context().Value(constants.USER_ID_CONTEXT_KEY).(uuid.UUID)

			dto := userFindOneDTO{
				id: authUserId,
			}
			user, err := s.UserRepository.userFindOne(r.Context(), dto)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, err)
				return
			}

			type response struct {
				utils.BaseResponse
				User User `json:"user"`
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

func (s *Service) userRouter() *chi.Mux {
	userRouter := chi.NewRouter()
	userRouter.Use(s.Middlewares.AuthMiddleware())

	// find one user
	userRouter.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		urlId := chi.URLParam(r, "id")
		id, err := uuid.FromString(urlId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		dto := userFindOneDTO{
			id: id,
		}

		user, err := s.UserRepository.userFindOne(r.Context(), dto)
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

		dto := userFindOneByUsernameDTO{
			username: req.Username,
		}

		user, err := s.UserRepository.userFindOneByUsername(r.Context(), dto)
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
		urlId := chi.URLParam(r, "id")
		id, err := uuid.FromString(urlId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		authUserId := r.Context().Value(constants.USER_ID_CONTEXT_KEY).(uuid.UUID)
		if authUserId != id {
			utils.WriteError(w, http.StatusForbidden, userForbiddenError)
			return
		}

		type Request struct {
			Username *string `json:"username,omitempty"`
			Password *string `json:"password,omitempty"`
		}
		req, err := utils.ReadJSON[Request](r)
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
			dto.passwordHash = &passwordHash
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}
		}

		user, err := s.UserRepository.userUpdate(r.Context(), dto)
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

		authUserId := r.Context().Value(constants.USER_ID_CONTEXT_KEY).(uuid.UUID)
		if authUserId != id {
			utils.WriteError(w, http.StatusForbidden, userForbiddenError)
			return
		}

		dto := deleteDTO{
			id: id,
		}
		user, err := s.UserRepository.userDelete(r.Context(), dto)
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
