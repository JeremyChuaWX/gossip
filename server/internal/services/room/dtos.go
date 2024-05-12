package room

import "github.com/gofrs/uuid/v5"

type CreateDTO struct {
	Name string
}

type FindOneDTO struct {
	Id uuid.UUID
}

type FindOneByUsernameDTO struct {
	Username string
}

type UpdateDTO struct {
	Id   uuid.UUID
	Name *string
}

type DeleteDTO struct {
	Id uuid.UUID
}
