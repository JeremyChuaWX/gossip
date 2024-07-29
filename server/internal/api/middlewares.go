package api

import (
	"context"
	"errors"
	"gossip/internal/repository"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

var invalidSessionIdError = errors.New("invalid session ID")

func (api *Api) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionIdHeaderValue := r.Header.Get(SESSION_ID_HEADER)
		if sessionIdHeaderValue == "" {
			slog.Error("empty session ID header")
			errorToJSON(w, http.StatusUnauthorized, invalidSessionIdError)
			return
		}
		sessionId, err := uuid.FromString(sessionIdHeaderValue)
		if err != nil {
			slog.Error(
				"invalid session ID",
				"sessionIdHeaderValue",
				sessionIdHeaderValue,
			)
			errorToJSON(w, http.StatusUnauthorized, invalidSessionIdError)
			return
		}
		res, err := api.Repository.UserSessionFindOne(
			r.Context(),
			repository.UserSessionFindOneParams{SessionId: sessionId},
		)
		if err != nil {
			slog.Error("session ID not found", "sessionId", sessionId)
			errorToJSON(w, http.StatusUnauthorized, invalidSessionIdError)
			return
		}
		nextReq := r.WithContext(
			context.WithValue(r.Context(), USER_SESSION_CONTEXT_KEY, res),
		)
		next.ServeHTTP(w, nextReq)
	})
}
