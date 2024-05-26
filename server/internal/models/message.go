package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Message struct {
	Id        uuid.UUID `db:"id" json:"id"`
	UserId    uuid.UUID `db:"user_id" json:"userId"`
	RoomId    uuid.UUID `db:"room_id" json:"roomId"`
	Body      string    `db:"body" json:"body"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
}
