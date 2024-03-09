package user

import "github.com/gofrs/uuid/v5"

type CreateDTO struct {
	Username     string
	PasswordHash string
}

type FindOneDTO struct {
	Id uuid.UUID
}

type FindOneByUsernameDTO struct {
	Username string
}

type UpdateDTO struct {
	Id           uuid.UUID
	Username     *string
	PasswordHash *string
}

type DeleteDTO struct {
	Id uuid.UUID
}
