package room

import "github.com/gofrs/uuid/v5"

type Room struct {
	Id   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
}
