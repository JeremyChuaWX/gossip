package chat

import (
	"github.com/gofrs/uuid/v5"
)

type event interface{}

type messageEvent struct {
	payload *message
	roomId  uuid.UUID
	userId  uuid.UUID
}

func newMessageEvent(message *message) (messageEvent, error) {
	roomId, err := uuid.FromString(message.RoomId)
	if err != nil {
		return messageEvent{}, err
	}
	userId, err := uuid.FromString(message.UserId)
	if err != nil {
		return messageEvent{}, err
	}
	return messageEvent{
		payload: message,
		roomId:  roomId,
		userId:  userId,
	}, nil
}

type userConnectEvent struct {
	user *user
}

type userDisconnectEvent struct {
	userId uuid.UUID
}

type userJoinRoomEvent struct {
	roomId uuid.UUID
	userId uuid.UUID
}

type userLeaveRoomEvent struct {
	roomId uuid.UUID
	userId uuid.UUID
}