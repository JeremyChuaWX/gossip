package chat

import (
	"errors"
	"gossip/internal/services/roomuser"
	"gossip/internal/services/user"
	"log"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
)

type service struct {
	userService     user.Service
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
	userId uuid.UUID,
) error {
	_, ok := s.chatUsers[userId]
	if ok {
		return errors.New("chat user already connected")
	}
	roomIds, err := s.roomUserService.FindRoomIdsByUserId(
		ctx,
		roomuser.FindRoomIdsByUserIdDTO{
			UserId: userId,
		},
	)
	if err != nil {
		return err
	}
	_user, err := s.userService.FindOne(
		ctx,
		user.FindOneDTO{
			Id: userId,
		},
	)
	if err != nil {
		return err
	}
	s.chatUsers[userId] = newChatUser(s, &_user, roomIds, conn)
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
