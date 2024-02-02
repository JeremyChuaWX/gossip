package room

import (
	"errors"
	"gossip/internal/client"
	"gossip/internal/message"

	"github.com/gorilla/websocket"
)

var (
	ClientExistsError       = errors.New("client already exists in room")
	ClientDoesNotExistError = errors.New("client does not exist in room")
)

type Room struct {
	Name     string
	Clients  map[string]*client.Client
	Messages chan message.MessageDTO
}

func New(name string) *Room {
	return &Room{
		Name:     name,
		Clients:  make(map[string]*client.Client),
		Messages: make(chan message.MessageDTO),
	}
}

func (r *Room) AddClient(
	name string,
	conn *websocket.Conn,
) (*client.Client, error) {
	if _, ok := r.Clients[name]; ok {
		return nil, ClientExistsError
	}
	c := &client.Client{
		Name: name,
		Room: r.Name,
		Conn: conn,
	}
	r.Clients[c.Name] = c
	return c, nil
}

func (r *Room) Client(name string) (*client.Client, error) {
	client, ok := r.Clients[name]
	if !ok {
		return nil, ClientDoesNotExistError
	}
	return client, nil
}

func (r *Room) RemoveClient(name string) error {
	if _, ok := r.Clients[name]; !ok {
		return ClientDoesNotExistError
	}
	delete(r.Clients, name)
	return nil
}

func (r *Room) BroadcastMessages() {
	for {
		msg := <-r.Messages
		for client := range r.Clients {
			c := r.Clients[client]
			c.Conn.WriteJSON(msg)
		}
	}
}
