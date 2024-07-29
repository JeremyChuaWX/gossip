package api

import (
	"context"
	"encoding/json"
	"gossip/internal/repository"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

func sessionIdFromHeader(header http.Header) (uuid.UUID, error) {
	sessionId := header.Get(SESSION_ID_HEADER)
	return uuid.FromString(sessionId)
}

func userSessionFromContext(
	ctx context.Context,
) repository.UserSessionFindOneResult {
	return ctx.Value(USER_SESSION_CONTEXT_KEY).(repository.UserSessionFindOneResult)
}

type baseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func readJSON[T any](r *http.Request) (T, error) {
	var res T
	err := json.NewDecoder(r.Body).Decode(&res)
	return res, err
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("error writing JSON body", "error", err.Error())
	}
}

func errorToJSON(w http.ResponseWriter, status int, err error) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	encodingErr := json.NewEncoder(w).Encode(baseResponse{
		Success: false,
		Message: err.Error(),
	})
	slog.Error("error writing JSON body", "error", encodingErr.Error())
}
