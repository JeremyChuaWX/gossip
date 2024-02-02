package message

import "time"

type Message struct {
	Timestamp time.Time
	Client    string
	Body      string
}

type MessageDTO struct {
	Client  string `json:"client"`
	Message string `json:"message"`
}
