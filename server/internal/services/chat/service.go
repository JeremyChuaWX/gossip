package chat

import (
	"errors"
	"gossip/internal/models"
	"gossip/internal/services/message"
	"gossip/internal/services/room"
	"gossip/internal/services/roomuser"
	"gossip/internal/services/user"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	ingress         chan event
	wsUpgrader      *websocket.Upgrader
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
		ingress: make(chan event),
		wsUpgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		userService:     userService,
		roomService:     roomService,
		roomUserService: roomUserService,
		messageService:  messageService,
		chatUsers:       make(map[uuid.UUID]*chatUser),
		chatRooms:       make(map[uuid.UUID]*chatRoom),
	}
	if err := s.initRooms(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Service) UserConnect(
	ctx context.Context,
	user *models.User,
	w http.ResponseWriter,
	r *http.Request,
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
	conn, err := s.upgradeConn(w, r)
	if err != nil {
		return err
	}
	s.chatUsers[user.Id] = newChatUser(s, user, roomIds, conn)
	return nil
}

func (s *Service) UserDisconnect(userId uuid.UUID) error {
	_, ok := s.chatUsers[userId]
	if !ok {
		return errors.New("chat user not found")
	}
	delete(s.chatUsers, userId)
	return nil
}

func (s *Service) UserJoinRoom(userId uuid.UUID, roomId uuid.UUID) {
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

func (s *Service) UserLeaveRoom(userId uuid.UUID, roomId uuid.UUID) {
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

func (s *Service) initRooms() error {
	rooms, err := s.roomService.FindMany(context.Background())
	if err != nil {
		return err
	}
	group, _ := errgroup.WithContext(context.Background())
	group.SetLimit(len(rooms))
	for _, room := range rooms {
		group.Go(func() error {
			roomUsers, err := s.roomUserService.FindUserIdsByRoomId(
				context.Background(),
				roomuser.FindUserIdsByRoomIdDTO{
					RoomId: room.Id,
				},
			)
			if err != nil {
				return err
			}
			s.chatRooms[room.Id] = newChatRoom(s, &room, roomUsers)
			return nil
		})
	}
	return group.Wait()
}

func (s *Service) upgradeConn(
	w http.ResponseWriter,
	r *http.Request,
) (*websocket.Conn, error) {
	conn, err := s.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
		return nil
	})
	return conn, nil
}
