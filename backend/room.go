package main

import (
	"fmt"
	"time"
)

type Room struct {
	Name     string
	Clients  map[string]*Client
	Messages chan RoomMessage
}

func NewRoom(name string) *Room {
	return &Room{
		Name:     name,
		Clients:  make(map[string]*Client),
		Messages: make(chan RoomMessage),
	}
}

type RoomMessage struct {
	Timestamp time.Time
	Client    *Client
	Body      string
}

func (rm RoomMessage) String() string {
	return fmt.Sprintf(
		"[%v] %s: %s",
		rm.Timestamp.Format(time.UnixDate),
		rm.Client.Name,
		rm.Body,
	)
}
