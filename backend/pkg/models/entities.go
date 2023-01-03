package models

import (
	"database/sql"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"
)

type Base struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = shortuuid.New()

	return
}

type User struct {
	Base
	Username   string    `gorm:"unique;not null" json:"username"`
	Email      string    `json:"email"`
	Password   string    `gorm:"not null" json:"-"`
	Posts      []Post    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"posts"`
	Comments   []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"comments"`
	Subscribed []Post    `gorm:"many2many:subscriptions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"subscribed"`
	IsPublic   bool      `json:"is_public"`
}

type Post struct {
	Base
	UserID    string    `gorm:"not null" json:"user_id"`
	Comments  []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"comments"`
	Tags      []Tag     `gorm:"many2many:taggable;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"tags"`
	PostScore int       `json:"post_score"`
	Title     string    `gorm:"not null" json:"title"`
	Body      string    `gorm:"not null" json:"body"`
}

type Comment struct {
	Base
	UserID       string         `gorm:"not null" json:"user_id"`
	PostID       string         `gorm:"not null" json:"post_id"`
	ParentID     sql.NullString `json:"parent_id"`
	Replies      []Comment      `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"replies"`
	CommentScore int            `json:"comment_score"`
	Body         string         `gorm:"not null" json:"body"`
}

type Tag struct {
	Base
	Name  string `gorm:"not null" json:"name"`
	Posts []Post `gorm:"many2many:taggable;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"posts"`
}

type Subscription struct {
	Base
	UserID string `gorm:"not null" json:"user_id"`
	PostID string `gorm:"not null" json:"post_id"`
}

type Taggable struct {
	Base
	PostID string `gorm:"not null" json:"post_id"`
	TagID  string `gorm:"not null" json:"tag_id"`
}
