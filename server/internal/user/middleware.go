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
			sessionCookie, err := r.Cookie("sessionId")
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, err)
				return
			}
			sessionId := sessionCookie.Value

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

			dto := userFindOneDTO{
				id: userId,
			}
			user, err := repository.userFindOne(r.Context(), dto)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, err)
				return
			}

			next.ServeHTTP(
				w,
				r.WithContext(
					context.WithValue(r.Context(), USER_CONTEXT_KEY, user),
				),
			)
		})
	}
}
