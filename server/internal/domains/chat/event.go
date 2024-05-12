package chat

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
