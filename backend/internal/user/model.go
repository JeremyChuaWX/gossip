package user

import "github.com/gofrs/uuid/v5"

type User struct {
	Id       uuid.UUID `db:"id" json:"id"`
	Username string    `db:"username" json:"username"`
	Password string    `db:"password" json:"-"`
}
