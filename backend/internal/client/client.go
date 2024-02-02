package client

import (
	"gossip/internal/message"

	"github.com/gorilla/websocket"
)

type Client struct {
	Name string
	Room string
	Conn *websocket.Conn
}

func (c *Client) ReadMessages() {
	for {
		var msg message.MessageDTO
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			c.Conn.Close()
			return
		}
	}
}
