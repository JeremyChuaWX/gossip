package models

import "time"

type Base struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	Base
	Username   string    `gorm:"unique;not null" json:"username"`
	Email      string    `json:"email"`
	Password   string    `gorm:"not null" json:"password"`
	Posts      []Post    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"posts"`
	Comments   []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"comments"`
	Subscribed []Post    `gorm:"many2many:subscriptions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"subscribed"`
	IsPublic   bool      `json:"is_public"`
}

type Post struct {
	Base
	UserID    int       `gorm:"not null" json:"user_id"`
	Comments  []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"comments"`
	Tags      []Tag     `gorm:"many2many:taggable;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"tags"`
	PostScore int       `json:"post_score"`
	Title     string    `gorm:"not null" json:"title"`
	Body      string    `gorm:"not null" json:"body"`
}

type Comment struct {
	Base
	UserID       int       `gorm:"not null" json:"user_id"`
	PostID       int       `gorm:"not null" json:"post_id"`
	ParentID     *int      `json:"parent_id"`
	Replies      []Comment `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CommentScore int       `json:"comment_score"`
	Body         string    `gorm:"not null" json:"body"`
}

type Tag struct {
	Base
	Name  string `gorm:"not null"`
	Posts []Post `gorm:"many2many:taggable;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"posts"`
}

type Subscriptions struct {
	Base
	UserID int `gorm:"primaryKey"`
	PostID int `gorm:"primaryKey"`
}

type Taggable struct {
	Base
	PostID int `gorm:"primaryKey"`
	TagID  int `gorm:"primaryKey"`
}
