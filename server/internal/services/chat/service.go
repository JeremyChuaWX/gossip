package chat

import (
	"errors"
	"gossip/internal/models"
	"gossip/internal/services/roomuser"
	"log"

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

func (s *service) userConnect(
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

func (s *service) userDisconnect(userId uuid.UUID) error {
	chatUser, ok := s.chatUsers[userId]
	if !ok {
		return errors.New("chat user not found")
	}
	return chatUser.disconnect()
}

func (s *service) userJoinRoom(userId uuid.UUID, roomId uuid.UUID) {
	chatUser, ok := s.chatUsers[userId]
	if !ok {
		log.Println("user not found")
		return
	}
	chatRoom, ok := s.chatRooms[roomId]
	if !ok {
		log.Println("room not found")
		return
	}
	event := newUserJoinRoomEvent(userId, roomId)
	chatUser.ingress <- event
	chatRoom.ingress <- event
}

func (s *service) userLeaveRoom(userId uuid.UUID, roomId uuid.UUID) {
	chatUser, ok := s.chatUsers[userId]
	if !ok {
		log.Println("user not found")
		return
	}
	chatRoom, ok := s.chatRooms[roomId]
	if !ok {
		log.Println("room not found")
		return
	}
	event := newUserLeaveRoomEvent(userId, roomId)
	chatUser.ingress <- event
	chatRoom.ingress <- event
}
