package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type UserSession struct {
	Id        uuid.UUID `db:"id" json:"id"`
	UserId    uuid.UUID `db:"user_id" json:"userId"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
}
