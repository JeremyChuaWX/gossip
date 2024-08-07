package router

import (
	"context"
	"fmt"
	"gossip/internal/repository"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

func (router *Router) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionIdCookie, err := r.Cookie(SESSION_ID_COOKIE)
		if err != nil || sessionIdCookie.Value == "" {
			slog.Error("empty session ID cookie")
			next.ServeHTTP(w, r)
			return
		}
		sessionId, err := uuid.FromString(sessionIdCookie.Value)
		if err != nil {
			slog.Error(
				"invalid session ID",
				"sessionIdCookie.Value",
				sessionIdCookie.Value,
			)
			next.ServeHTTP(w, r)
			return
		}
		res, err := router.Repository.SessionFindOne(
			r.Context(),
			repository.SessionFindOneParams{SessionId: sessionId},
		)
		if err != nil {
			slog.Error("session ID not found", "sessionId", sessionId)
			next.ServeHTTP(w, r)
			return
		}
		nextReq := r.WithContext(
			context.WithValue(r.Context(), USER_SESSION_CONTEXT_KEY, res),
		)
		next.ServeHTTP(w, nextReq)
	})
}

func (router *Router) pagesAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := sessionFromContext(r.Context())
		if err != nil {
			url := fmt.Sprintf("/login?prev=%s", r.URL.Path)
			http.Redirect(w, r, url, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (router *Router) apiAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := sessionFromContext(r.Context())
		if err != nil {
			errorToJSON(w, http.StatusUnauthorized, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}
