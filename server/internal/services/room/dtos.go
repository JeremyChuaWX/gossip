package room

import "github.com/gofrs/uuid/v5"

type CreateDTO struct {
	Name string
}

type FindOneDTO struct {
	RoomId uuid.UUID
}

type FindOneByUsernameDTO struct {
	Username string
}

type UpdateDTO struct {
	RoomId   uuid.UUID
	Name *string
}

type DeleteDTO struct {
	RoomId uuid.UUID
}
