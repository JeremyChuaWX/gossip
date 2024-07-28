package middlewares

import (
	"context"
	"errors"
	"gossip/internal/api"
	"gossip/internal/repository"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

var invalidSessionIdError = errors.New("invalid session ID")

type Middlewares struct {
	Repository *repository.Repository
}

func (m *Middlewares) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionIdHeaderValue := r.Header.Get(api.SESSION_ID_HEADER)
		if sessionIdHeaderValue == "" {
			slog.Error("empty session ID header")
			api.ErrorToJSON(
				w,
				http.StatusUnauthorized,
				invalidSessionIdError,
			)
			return
		}
		sessionId, err := uuid.FromString(sessionIdHeaderValue)
		if err != nil {
			slog.Error("invalid session ID", "sessionId", sessionId)
			api.ErrorToJSON(
				w,
				http.StatusUnauthorized,
				invalidSessionIdError,
			)
			return
		}
		res, err := m.Repository.UserSessionFindOne(
			r.Context(),
			repository.UserSessionFindOneParams{SessionId: sessionId},
		)
		if err != nil {
			slog.Error("session ID not found", "sessionId", sessionId)
			api.ErrorToJSON(
				w,
				http.StatusUnauthorized,
				invalidSessionIdError,
			)
			return
		}
		nextReq := r.WithContext(
			context.WithValue(
				r.Context(),
				api.USER_SESSION_CONTEXT_KEY,
				res,
			),
		)
		next.ServeHTTP(w, nextReq)
	})
}
