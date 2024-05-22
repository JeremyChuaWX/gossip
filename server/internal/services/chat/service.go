package chat

import (
	"errors"
	"gossip/internal/services/message"
	"gossip/internal/services/room"
	"gossip/internal/services/roomuser"
	"gossip/internal/services/user"
	"log"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
)

type Service struct {
	userService     *user.Service
	roomService     *room.Service
	roomUserService *roomuser.Service
	messageService  *message.Service
	chatUsers       map[uuid.UUID]*chatUser
	chatRooms       map[uuid.UUID]*chatRoom
}

func InitService(
	userService *user.Service,
	roomService *room.Service,
	roomUserService *roomuser.Service,
	messageService *message.Service,
) (*Service, error) {
	s := &Service{
		userService:     userService,
		roomService:     roomService,
		roomUserService: roomUserService,
		messageService:  messageService,
		chatUsers:       make(map[uuid.UUID]*chatUser),
		chatRooms:       make(map[uuid.UUID]*chatRoom),
	}

	ctx := context.Background()
	rooms, err := roomService.FindMany(ctx)
	if err != nil {
		return nil, err
	}
	errors := make(chan error, len(rooms))
	chatRooms := make(chan *chatRoom, len(rooms))
	for _, room := range rooms {
		go func() {
			roomUsers, err := roomUserService.FindUserIdsByRoomId(
				ctx,
				roomuser.FindUserIdsByRoomIdDTO{
					RoomId: room.Id,
				},
			)
			if err != nil {
				errors <- err
				return
			}
			chatRooms <- newChatRoom(s, &room, roomUsers)
		}()
	}

	counter := len(rooms)
	for {
		select {
		case err = <-errors:
			return nil, err
		case chatRoom := <-chatRooms:
			s.chatRooms[chatRoom.room.Id] = chatRoom
			counter--
		default:
			if counter <= 0 {
				return s, nil
			}
		}
	}
}

func (s *Service) userConnect(
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
			UserId: userId,
		},
	)
	if err != nil {
		return err
	}
	s.chatUsers[userId] = newChatUser(s, &_user, roomIds, conn)
	return nil
}

func (s *Service) userDisconnect(userId uuid.UUID) error {
	chatUser, ok := s.chatUsers[userId]
	if !ok {
		return errors.New("chat user not found")
	}
	delete(s.chatUsers, userId)
	return chatUser.disconnect()
}

func (s *Service) userJoinRoom(userId uuid.UUID, roomId uuid.UUID) {
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

func (s *Service) userLeaveRoom(userId uuid.UUID, roomId uuid.UUID) {
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
