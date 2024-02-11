package user

import "github.com/gofrs/uuid/v5"

type createDTO struct {
	Username string
	Password []byte
}

type findOneDTO struct {
	Id uuid.UUID
}

type updateDTO struct {
	Id       uuid.UUID
	Username *string
	Password []byte
}

type deleteDTO struct {
	Id uuid.UUID
}
