package api

import (
	"context"
	"encoding/json"
	"gossip/internal/repository"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"
)

func SessionIdFromHeader(header http.Header) (uuid.UUID, error) {
	sessionId := header.Get(SESSION_ID_HEADER)
	return uuid.FromString(sessionId)
}

func UserSessionFromContext(
	ctx context.Context,
) repository.UserSessionFindOneResult {
	return ctx.Value(USER_SESSION_CONTEXT_KEY).(repository.UserSessionFindOneResult)
}

func WalkRoutes(
	method string,
	route string,
	handler http.Handler,
	middlewares ...func(http.Handler) http.Handler,
) error {
	route = strings.Replace(route, "/*/", "/", -1)
	log.Printf("%s %s\n", method, route)
	return nil
}

type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ReadJSON[T any](r *http.Request) (T, error) {
	var res T
	err := json.NewDecoder(r.Body).Decode(&res)
	return res, err
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("error writing JSON body", "error", err.Error())
	}
}

func ErrorToJSON(w http.ResponseWriter, status int, err error) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	encodingErr := json.NewEncoder(w).Encode(BaseResponse{
		Success: false,
		Message: err.Error(),
	})
	slog.Error("error writing JSON body", "error", encodingErr.Error())
}
