package middlewares

import (
	"context"
	"errors"
	"gossip/internal/api"
	"gossip/internal/services/session"
	"gossip/internal/utils/httpjson"
	"net/http"
)

var invalidSessionIdError = errors.New("invalid session ID")

type Middlewares struct {
	SessionService *session.Service
}

func (m *Middlewares) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.Header.Get(api.SESSION_ID_HEADER)
		if sessionId == "" {
			httpjson.WriteError(
				w,
				http.StatusUnauthorized,
				invalidSessionIdError,
			)
			return
		}
		userId, err := m.SessionService.Get(r.Context(), sessionId)
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		nextReq := r.WithContext(
			context.WithValue(
				r.Context(),
				api.USER_ID_CONTEXT_KEY,
				userId,
			),
		)
		next.ServeHTTP(w, nextReq)
	})
}
