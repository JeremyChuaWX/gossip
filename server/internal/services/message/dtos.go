package message

import "github.com/gofrs/uuid/v5"

type SaveDto struct {
	UserId uuid.UUID
	RoomId uuid.UUID
	Body   string
}

type FindManyByRoomIdDto struct {
	RoomId uuid.UUID
}
