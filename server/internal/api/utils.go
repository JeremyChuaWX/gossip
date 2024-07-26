package api

import (
	"context"
	"gossip/internal/repository"
	"log"
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
