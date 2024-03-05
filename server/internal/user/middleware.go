package user

import (
	"context"
	"gossip/internal/utils"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type Middleware func(http.Handler) http.Handler

func AuthMiddleware(repository *Repository) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionId := r.Header.Get(SESSION_ID_HEADER)

			userIdStr, err := repository.sessionsGet(r.Context(), sessionId)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, err)
				return
			}

			userId, err := uuid.FromString(userIdStr)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, err)
				return
			}

			next.ServeHTTP(
				w,
				r.WithContext(
					context.WithValue(r.Context(), USER_ID_CONTEXT_KEY, userId),
				),
			)
		})
	}
}