package models

import "github.com/gofrs/uuid/v5"

type User struct {
	Id           uuid.UUID `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
}
