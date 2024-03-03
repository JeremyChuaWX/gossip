package user

import "github.com/gofrs/uuid/v5"

type userCreateDTO struct {
	username     string
	passwordHash []byte
}

type userFindOneDTO struct {
	id uuid.UUID
}

type userFindOneByUsernameDTO struct {
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
