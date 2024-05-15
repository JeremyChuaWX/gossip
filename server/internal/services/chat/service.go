package chat

import (
	"errors"
	"gossip/internal/models"
	"gossip/internal/services/roomuser"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
)

type service struct {
	roomUserService roomuser.Service
	chatUsers       map[uuid.UUID]*chatUser
	chatRooms       map[uuid.UUID]*chatRoom
}

func InitService(roomUserService roomuser.Service) *service {
	s := &service{
		roomUserService: roomUserService,
		chatUsers:       make(map[uuid.UUID]*chatUser),
		chatRooms:       make(map[uuid.UUID]*chatRoom),
	}
	// TODO: init chat rooms
	return s
}

func (s *service) chatUserConnect(
	ctx context.Context,
	conn *websocket.Conn,
	user *models.User,
) error {
	_, ok := s.chatUsers[user.Id]
	if ok {
		return errors.New("chat user already connected")
	}
	roomIds, err := s.roomUserService.FindRoomIdsByUserId(
		ctx,
		roomuser.FindRoomIdsByUserIdDTO{
			UserId: user.Id,
		},
	)
	if err != nil {
		return err
	}
	s.chatUsers[user.Id] = newChatUser(s, user, roomIds, conn)
	return nil
}

func (s *service) chatUserDisconnect(chatUserId uuid.UUID) error {
	chatUser, ok := s.chatUsers[chatUserId]
	if !ok {
		return errors.New("chat user not found")
	}
	return chatUser.disconnect()
}
