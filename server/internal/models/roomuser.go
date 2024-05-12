package models

import "github.com/gofrs/uuid/v5"

type RoomUser struct {
	RoomId uuid.UUID `db:"room_id" json:"roomId"`
	UserId uuid.UUID `db:"user_id" json:"userId"`
}
