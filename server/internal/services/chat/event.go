package chat

import "github.com/gofrs/uuid/v5"

type eventName string

const (
	USER_CONNECT    eventName = "USER_CONNECT"
	USER_DISCONNECT eventName = "USER_DISCONNECT"

	ROOM_CREATE  eventName = "ROOM_CREATE"
	ROOM_DESTROY eventName = "ROOM_DESTROY"

	MESSAGE         eventName = "MESSAGE"
	USER_JOIN_ROOM  eventName = "USER_JOIN_ROOM"
	USER_LEAVE_ROOM eventName = "USER_LEAVE_ROOM"
)

type event interface {
	name() eventName
}

// message event

type payload struct {
	RoomId string `json:"roomId"`
	UserId string `json:"userId"`
	Body   string `json:"body"`
}

type messageEvent struct {
	payload payload
}

func (e *messageEvent) name() eventName {
	return MESSAGE
}

func newMessageEvent(payload payload) event {
	return &messageEvent{
		payload: payload,
	}
}

// user join room

type userJoinRoomEvent struct {
	userId uuid.UUID
	roomId uuid.UUID
}

func (e *userJoinRoomEvent) name() eventName {
	return USER_JOIN_ROOM
}

func newUserJoinRoomEvent(userId uuid.UUID, roomId uuid.UUID) event {
	return &userJoinRoomEvent{
		userId: userId,
		roomId: roomId,
	}
}

// user leave room

type userLeaveRoomEvent struct {
	userId uuid.UUID
	roomId uuid.UUID
}

func (e *userLeaveRoomEvent) name() eventName {
	return USER_LEAVE_ROOM
}

func newUserLeaveRoomEvent(userId uuid.UUID, roomId uuid.UUID) event {
	return &userLeaveRoomEvent{
		userId: userId,
		roomId: roomId,
	}
}
