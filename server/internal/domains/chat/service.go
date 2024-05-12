package chat

import (
	"errors"
	"gossip/internal/domains/room"
	"gossip/internal/domains/user"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
)

type service struct {
	chatUsers map[uuid.UUID]*chatUser
	chatRooms map[uuid.UUID]*chatRoom
}

func InitService() *service {
	s := &service{
		chatUsers: make(map[uuid.UUID]*chatUser),
		chatRooms: make(map[uuid.UUID]*chatRoom),
	}
	return s
}

func (s *service) chatUserConnect(
	conn *websocket.Conn,
	user *user.User,
	rooms []room.Room,
) error {
	_, ok := s.chatUsers[user.Id]
	if ok {
		return errors.New("chat user already connected")
	}
	chatUser := &chatUser{
		service: s,
		ingress: make(chan event),
		user:    user,
		roomIds: make(map[uuid.UUID]bool),
		conn:    conn,
	}
	for _, room := range rooms {
		chatUser.roomIds[room.Id] = true
	}
	s.chatUsers[user.Id] = chatUser
	return nil
}

func (s *service) chatUserDisconnect(chatUserId uuid.UUID) error {
	chatUser, ok := s.chatUsers[chatUserId]
	if !ok {
		return errors.New("chat user not found")
	}
	return chatUser.disconnect()
}
