package user

import "time"

const SESSION_EXPIRATION = time.Hour * 24 * 7

type UserContextKey string

const USER_ID_CONTEXT_KEY UserContextKey = "USER_ID"
