package chat

import "time"

type message struct {
	RoomId    string    `json:"roomId"`
	UserId    string    `json:"userId"`
	Username  string    `json:"username"`
	Body      string    `json:"body"`
	Timestamp time.Time `json:"timestamp"`
}
