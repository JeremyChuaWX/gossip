package user

import "github.com/gofrs/uuid/v5"

type CreateDTO struct {
	Username     string
	PasswordHash string
}

type FindOneDTO struct {
	UserId uuid.UUID
}

type FindOneByUsernameDTO struct {
	Username string
}

type UpdateDTO struct {
	UserId       uuid.UUID
	Username     *string
	PasswordHash *string
}

type DeleteDTO struct {
	UserId uuid.UUID
}
