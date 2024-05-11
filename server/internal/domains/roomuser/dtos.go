package roomuser

import "github.com/gofrs/uuid/v5"

type CreateDTO struct {
	RoomId uuid.UUID
	UserId uuid.UUID
}

type FindManyByRoomIdDTO struct {
	RoomId uuid.UUID
}

type FindManyByUserIdDTO struct {
	UserId uuid.UUID
}

type DeleteDTO struct {
	RoomId uuid.UUID
	UserId uuid.UUID
}
