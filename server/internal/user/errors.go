package user

import "errors"

var (
	userForbiddenError    = errors.New("mismatch user ID")
	invalidSessionIdError = errors.New("invalid session ID")
)
