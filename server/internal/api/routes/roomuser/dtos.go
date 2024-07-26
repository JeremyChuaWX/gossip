package roomuser

import "github.com/gofrs/uuid/v5"

type UserJoinRoomDTO struct {
	UserId uuid.UUID
	RoomId uuid.UUID
}

type UserLeaveRoomDTO struct {
	UserId uuid.UUID
	RoomId uuid.UUID
}

type FindRoomIdsByUserIdDTO struct {
	UserId uuid.UUID
}

type FindUserIdsByRoomIdDTO struct {
	RoomId uuid.UUID
}
