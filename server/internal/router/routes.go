package router

import (
	"gossip/internal/repository"
	"gossip/internal/utils/password"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (router *Router) routeGroup(mux chi.Router) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/index.html")
	})

	mux.Get("/signup", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/signup.html")
	})

	mux.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
		body, err := readJSON[struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}](r)
		if err != nil {
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		passwordHash, err := password.Hash(body.Password)
		if err != nil {
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		user, err := router.Repository.UserCreate(
			r.Context(),
			repository.UserCreateParams{
				Username:     body.Username,
				PasswordHash: passwordHash,
			},
		)
		if err != nil {
			errorToJSON(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "signed up",
			Data: map[string]any{
				"user": map[string]any{
					"id": user.UserId,
				},
			},
		})
	})

	mux.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/login.html")
	})

	mux.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		body, err := readJSON[struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}](r)
		if err != nil {
			errorToJSON(w, http.StatusBadRequest, err)
			return
		}
		user, err := router.Repository.UserFindOneByUsername(
			r.Context(),
			repository.UserFindOneByUsernameParams{Username: body.Username},
		)
		if err != nil {
			errorToJSON(w, http.StatusUnauthorized, err)
			return
		}
		err = password.Verify(body.Password, user.PasswordHash)
		if err != nil {
			errorToJSON(w, http.StatusUnauthorized, err)
			return
		}
		session, err := router.Repository.UserSessionCreate(
			r.Context(),
			repository.UserSessionCreateParams{UserId: user.UserId},
		)
		if err != nil {
			errorToJSON(w, http.StatusUnauthorized, err)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     SESSION_ID_COOKIE,
			Value:    session.SessionId.String(),
			Secure:   true,
			HttpOnly: true,
		})
		writeJSON(w, http.StatusOK, baseResponse{
			Success: true,
			Message: "logged in",
			Data: map[string]any{
				"session": map[string]any{
					"id":        session.SessionId,
					"expiresOn": session.ExpiresOn,
				},
			},
		})
	})
}

func (router *Router) authedRouteGroup(mux chi.Router) {
	mux.Use(router.authMiddleware)

	mux.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/home.html")
	})
}
