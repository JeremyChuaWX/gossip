package user

import "github.com/gofrs/uuid/v5"

type createDTO struct {
	username     string
	passwordHash []byte
}

type findOneDTO struct {
	id uuid.UUID
}

type findOneByUsernameDTO struct {
	username string
}

type updateDTO struct {
	id           uuid.UUID
	username     *string
	passwordHash []byte
}

type deleteDTO struct {
	id uuid.UUID
}
