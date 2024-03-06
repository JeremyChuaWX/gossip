package user

import "github.com/gofrs/uuid/v5"

type userCreateDTO struct {
	username     string
	passwordHash string
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
	passwordHash *string
}

type deleteDTO struct {
	id uuid.UUID
}
