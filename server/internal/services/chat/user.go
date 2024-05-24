package chat

import (
	"gossip/internal/models"
	"log/slog"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
)

type chatUser struct {
	service  *Service
	alive    chan bool
	ingress  chan event
	handlers map[eventName]func(*chatUser, event)
	user     *models.User
	roomIds  map[uuid.UUID]bool
	conn     *websocket.Conn
}

func newChatUser(
	service *Service,
	user *models.User,
	roomUsers []models.RoomUser,
	conn *websocket.Conn,
) *chatUser {
	u := &chatUser{
		service:  service,
		alive:    make(chan bool),
		ingress:  make(chan event),
		handlers: make(map[eventName]func(*chatUser, event)),
		user:     user,
		roomIds:  make(map[uuid.UUID]bool),
		conn:     conn,
	}
	for _, roomUser := range roomUsers {
		u.roomIds[roomUser.RoomId] = true
	}
	u.registerHandlers()
	go u.receiveEvents()
	go u.receiveWebsocket()

	return u
}

func (u *chatUser) receiveEvents() {
	for {
		select {
		case <-u.alive:
			return
		case e := <-u.ingress:
			handler, ok := u.handlers[e.name()]
			if !ok {
				slog.Error("invalid event")
				continue
			}
			handler(u, e)
		}
	}
}

func (u *chatUser) receiveWebsocket() {
	for {
		select {
		case <-u.alive:
			return
		default:
			var payload payload
			if err := u.conn.ReadJSON(&payload); err != nil {
				return
			}
			roomId, err := uuid.FromString(payload.RoomId)
			if err != nil {
				slog.Error("invalid room ID")
				continue
			}
			room, ok := u.service.chatRooms[roomId]
			if !ok {
				slog.Error("room not found", "roomId", roomId.String())
				continue
			}
			userId, err := uuid.FromString(payload.UserId)
			if err != nil {
				slog.Error("invalid user ID")
				continue
			}
			if userId != u.user.Id {
				slog.Error(
					"wrong user",
					"userId",
					userId.String(),
					"connUserId",
					u.user.Id.String(),
				)
				continue
			}
			room.ingress <- newMessageEvent(roomId, userId, payload)
		}
	}
}

// handlers

func (u *chatUser) registerHandlers() {
	u.handlers[MESSAGE] = (*chatUser).messageEventHandler
	u.handlers[USER_JOIN_ROOM] = (*chatUser).userJoinRoomEventHandler
	u.handlers[USER_LEAVE_ROOM] = (*chatUser).userLeaveRoomEventHandler
}

func (u *chatUser) messageEventHandler(e event) {
	event := e.(*messageEvent)
	if err := u.conn.WriteJSON(event.payload); err != nil {
		return
	}
}

func (u *chatUser) userJoinRoomEventHandler(e event) {
	event := e.(*userJoinRoomEvent)
	if u.user.Id != event.userId {
		slog.Error("wrong user")
		return
	}
	if _, ok := u.roomIds[event.roomId]; ok {
		slog.Error("user already in room")
		return
	}
	u.roomIds[event.roomId] = true
}

func (u *chatUser) userLeaveRoomEventHandler(e event) {
	event := e.(*userJoinRoomEvent)
	if u.user.Id != event.userId {
		slog.Error("wrong user")
		return
	}
	if _, ok := u.roomIds[event.roomId]; !ok {
		slog.Error("user not in room")
		return
	}
	delete(u.roomIds, event.roomId)
}
