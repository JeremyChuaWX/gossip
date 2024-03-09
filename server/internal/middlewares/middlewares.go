package middlewares

import (
	"context"
	"errors"
	"gossip/internal/constants"
	"gossip/internal/domains/session"
	"gossip/internal/utils"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

var invalidSessionIdError = errors.New("invalid session ID")

type Middlewares struct {
	SessionRepository *session.Repository
}

type middleware func(http.Handler) http.Handler

func (m *Middlewares) AuthMiddleware() middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionId := r.Header.Get(constants.SESSION_ID_HEADER)
			if sessionId == "" {
				utils.WriteError(
					w,
					http.StatusUnauthorized,
					invalidSessionIdError,
				)
				return
			}

			userIdStr, err := m.SessionRepository.Get(r.Context(), sessionId)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, err)
				return
			}

			userId, err := uuid.FromString(userIdStr)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, err)
				return
			}

			nextReq := r.WithContext(
				context.WithValue(
					r.Context(),
					constants.USER_ID_CONTEXT_KEY,
					userId,
				),
			)

			next.ServeHTTP(w, nextReq)
		})
	}
}
