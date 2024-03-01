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

	authRouter := s.authRouter()
	router.Mount("/auth", authRouter)
}

func (s *Service) userRouter() *chi.Mux {
	userRouter := chi.NewRouter()

	// create user
	userRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		req, err := utils.ReadJSON[request](r)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		passwordHash, err := password.Hash(req.Password)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		dto := createDTO{
			username:     req.Username,
			passwordHash: passwordHash,
		}
		user, err := s.Repository.create(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			Message string `json:"message"`
			User    User   `json:"user"`
		}
		utils.WriteJSON(w, http.StatusCreated, response{
			Message: "created user",
			User:    user,
		})
	})

	// find one user
	userRouter.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		urlId := chi.URLParam(r, "id")
		id, err := uuid.FromString(urlId)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
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
			Message string `json:"message"`
			User    User   `json:"user"`
		}
		utils.WriteJSON(w, http.StatusOK, response{
			Message: "found user",
			User:    user,
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
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		urlId := chi.URLParam(r, "id")
		id, err := uuid.FromString(urlId)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
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
			Message string `json:"message"`
			User    User   `json:"user"`
		}
		utils.WriteJSON(w, http.StatusOK, response{
			Message: "updated user",
			User:    user,
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
			Message string `json:"message"`
			User    User   `json:"user"`
		}
		utils.WriteJSON(w, http.StatusOK, response{
			Message: "deleted user",
			User:    user,
		})
	})

	return userRouter
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

		dto := findOneByUsernameDTO{
			username: req.Username,
		}
		user, err := s.Repository.findOneByUsername(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if err = password.Compare([]byte(user.PasswordHash), req.Password); err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		type response struct {
			Message string `json:"message"`
			User    User   `json:"user"`
		}
		utils.WriteJSON(w, http.StatusCreated, response{
			Message: "signed in",
			User:    user,
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

		dto := createDTO{
			username:     req.Username,
			passwordHash: passwordHash,
		}
		user, err := s.Repository.create(r.Context(), dto)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			Message string `json:"message"`
			User    User   `json:"user"`
		}
		utils.WriteJSON(w, http.StatusCreated, response{
			Message: "signed up",
			User:    user,
		})
	})

	return authRouter
}
