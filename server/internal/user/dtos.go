package user

import "github.com/gofrs/uuid/v5"

type CreateDTO struct {
	Username     string
	PasswordHash []byte
}

type findOneDTO struct {
	id uuid.UUID
}

type FindOneByUsernameDTO struct {
	Username string
}

type updateDTO struct {
	id           uuid.UUID
	username     *string
	passwordHash []byte
}

type deleteDTO struct {
	id uuid.UUID
}
