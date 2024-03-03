package user

import "time"

const SESSION_EXPIRATION = time.Hour * 24 * 7

type UserContextKey string

const USER_CONTEXT_KEY UserContextKey = "user"
