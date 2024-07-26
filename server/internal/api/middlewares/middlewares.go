package middlewares

import (
	"context"
	"errors"
	"gossip/internal/api"
	"gossip/internal/utils/httpjson"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/redis/go-redis/v9"
)

var invalidSessionIdError = errors.New("invalid session ID")

type Middlewares struct {
	Redis *redis.Client
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
		userIdString, err := m.Redis.Get(r.Context(), sessionId).Result()
		if err != nil {
			httpjson.WriteError(
				w,
				http.StatusUnauthorized,
				invalidSessionIdError,
			)
			return
		}
		userId, err := uuid.FromString(userIdString)
		if err != nil {
			httpjson.WriteError(
				w,
				http.StatusUnauthorized,
				invalidSessionIdError,
			)
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
