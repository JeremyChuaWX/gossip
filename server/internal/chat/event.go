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

func newMessageEvent(
	userId uuid.UUID,
	username string,
	message *message,
) (messageEvent, error) {
	roomId, err := uuid.FromString(message.RoomId)
	if err != nil {
		return messageEvent{}, err
	}
	message.UserId = userId.String()
	message.Username = username
	return messageEvent{
		payload: message,
		roomId:  roomId,
		userId:  userId,
	}, nil
}

type roomCreatedEvent struct {
	room *room
}

type userConnectedEvent struct {
	user *user
}

type userDisconnectedEvent struct {
	userId uuid.UUID
}

type userJoinedRoomEvent struct {
	roomId uuid.UUID
	userId uuid.UUID
}

type userLeftRoomEvent struct {
	roomId uuid.UUID
	userId uuid.UUID
}
