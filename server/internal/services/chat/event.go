package chat

import (
	"github.com/gofrs/uuid/v5"
)

type eventName string

type event interface {
	name() eventName
}

const MESSAGE_EVENT eventName = "MESSAGE_EVENT"

type messageEvent struct {
	payload *message
	roomId  uuid.UUID
	userId  uuid.UUID
}

func newMessageEvent(message *message) (*messageEvent, error) {
	roomId, err := uuid.FromString(message.RoomId)
	if err != nil {
		return nil, err
	}
	userId, err := uuid.FromString(message.UserId)
	if err != nil {
		return nil, err
	}
	return &messageEvent{
		payload: message,
		roomId:  roomId,
		userId:  userId,
	}, nil
}

func (event *messageEvent) name() eventName {
	return MESSAGE_EVENT
}

const USER_CONNECT_EVENT eventName = "USER_CONNECT_EVENT"

type userConnectEvent struct {
	user *user
}

func (event *userConnectEvent) name() eventName {
	return USER_CONNECT_EVENT
}

const USER_DISCONNECT_EVENT eventName = "USER_DISCONNECT_EVENT"

type userDisconnectEvent struct {
	userId uuid.UUID
}

func (event *userDisconnectEvent) name() eventName {
	return USER_DISCONNECT_EVENT
}

const USER_JOIN_ROOM_EVENT eventName = "USER_JOIN_ROOM_EVENT"

type userJoinRoomEvent struct {
	roomId uuid.UUID
	userId uuid.UUID
}

func (event *userJoinRoomEvent) name() eventName {
	return USER_JOIN_ROOM_EVENT
}

const USER_LEAVE_ROOM_EVENT eventName = "USER_LEAVE_ROOM_EVENT"

type userLeaveRoomEvent struct {
	roomId uuid.UUID
	userId uuid.UUID
}

func (event *userLeaveRoomEvent) name() eventName {
	return USER_LEAVE_ROOM_EVENT
}
