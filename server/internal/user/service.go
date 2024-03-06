package user

import (
	"errors"
	"gossip/internal/password"
	"gossip/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

var UserForbiddenError = errors.New("mismatch user ID")

type Service struct {
	Repository *Repository
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
		user, err := s.Repository.userFindOneByUsername(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if err = password.Verify(req.Password, user.PasswordHash); err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		sessionId, err := s.Repository.sessionCreate(
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
		user, err := s.Repository.userCreate(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		sessionId, err := s.Repository.sessionCreate(
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

	// me
	authRouter.Get("/me", func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(USER_ID_CONTEXT_KEY).(uuid.UUID)

		dto := userFindOneDTO{
			id: id,
		}
		user, err := s.Repository.userFindOne(r.Context(), dto)
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
				Message: "current signed in user",
			},
			User: user,
		})
	})

	return authRouter
}

func (s *Service) userRouter() *chi.Mux {
	userRouter := chi.NewRouter()
	userRouter.Use(AuthMiddleware(s.Repository))

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

		user, err := s.Repository.userFindOne(r.Context(), dto)
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

		user, err := s.Repository.userFindOneByUsername(r.Context(), dto)
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

		authUserId := r.Context().Value(USER_ID_CONTEXT_KEY).(uuid.UUID)
		if authUserId != id {
			utils.WriteError(w, http.StatusForbidden, UserForbiddenError)
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

		user, err := s.Repository.userUpdate(r.Context(), dto)
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

		authUserId := r.Context().Value(USER_ID_CONTEXT_KEY).(uuid.UUID)
		if authUserId != id {
			utils.WriteError(w, http.StatusForbidden, UserForbiddenError)
			return
		}

		dto := deleteDTO{
			id: id,
		}
		user, err := s.Repository.userDelete(r.Context(), dto)
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
