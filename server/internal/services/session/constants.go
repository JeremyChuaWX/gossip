package session

import "time"

type UserContextKey string

const (
	USER_ID_CONTEXT_KEY UserContextKey = "USER_ID"
	SESSION_ID_HEADER                  = "x-session-id"
	SESSION_EXPIRATION                 = time.Hour * 24 * 7
)
