package chat

import (
	room "gossip/internal/domains/room"
	"gossip/internal/domains/user"

	"github.com/gofrs/uuid/v5"
)

type chatRoom struct {
	service *service
	ingress chan event
	room    room.Room
	userIds map[uuid.UUID]bool
}

func newChatRoom(room room.Room, users []*user.User) *chatRoom {
	r := &chatRoom{
		room: room,
	}
	for _, user := range users {
		r.userIds[user.Id] = true
	}
	return r
}
