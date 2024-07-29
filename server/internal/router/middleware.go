package router

import (
	"context"
	"gossip/internal/repository"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

func (router *Router) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionIdCookie, err := r.Cookie(SESSION_ID_COOKIE)
		if err != nil || sessionIdCookie.Value == "" {
			slog.Error("empty session ID cookie")
			return
		}
		sessionId, err := uuid.FromString(sessionIdCookie.Value)
		if err != nil {
			slog.Error(
				"invalid session ID",
				"sessionIdCookie.Value",
				sessionIdCookie.Value,
			)
			return
		}
		res, err := router.Repository.UserSessionFindOne(
			r.Context(),
			repository.UserSessionFindOneParams{SessionId: sessionId},
		)
		if err != nil {
			slog.Error("session ID not found", "sessionId", sessionId)
			return
		}
		nextReq := r.WithContext(
			context.WithValue(r.Context(), USER_SESSION_CONTEXT_KEY, res),
		)
		next.ServeHTTP(w, nextReq)
	})
}
