package chat

import "time"

type eventType string

const (
	NEW_CLIENT        eventType = "NEW_CLIENT"
	CLIENT_DISCONNECT eventType = "CLIENT_DISCONNECT"

	NEW_ROOM     eventType = "NEW_ROOM"
	DESTROY_ROOM eventType = "DESTROY_ROOM"

	MESSAGE           eventType = "MESSAGE"
	CLIENT_JOIN_ROOM  eventType = "CLIENT_JOIN_ROOM"
	CLIENT_LEAVE_ROOM eventType = "CLIENT_LEAVE_ROOM"
)

type event interface {
	name() eventType
}

// new client

type newClientEvent struct {
	_name  eventType
	client *client
}

func makeNewClientEvent(client *client) *newClientEvent {
	return &newClientEvent{
		_name:  NEW_CLIENT,
		client: client,
	}
}

func (e *newClientEvent) name() eventType {
	return e._name
}

// client disconnect

type clientDisconnectEvent struct {
	_name  eventType
	client *client
}

func makeClientDisconnectEvent(client *client) *clientDisconnectEvent {
	return &clientDisconnectEvent{
		_name:  CLIENT_DISCONNECT,
		client: client,
	}
}

func (e *clientDisconnectEvent) name() eventType {
	return e._name
}

// new room

type newRoomEvent struct {
	_name eventType
	room  *room
}

func makeNewRoomEvent(room *room) event {
	return &newRoomEvent{
		_name: NEW_ROOM,
		room:  room,
	}
}

func (e *newRoomEvent) name() eventType {
	return e._name
}

// remove room

type destroyRoomEvent struct {
	_name eventType
	room  *room
}

func makeDestroyRoomEvent(room *room) *destroyRoomEvent {
	return &destroyRoomEvent{
		_name: DESTROY_ROOM,
		room:  room,
	}
}

func (e *destroyRoomEvent) name() eventType {
	return e._name
}

// message

type messageEvent struct {
	_name     eventType
	timestamp time.Time
	message   string
}

func (e *messageEvent) name() eventType {
	return e._name
}

func (e *messageEvent) toJSON() messageJSON {
	return messageJSON{
		Timestamp: e.timestamp,
		Message:   e.message,
	}
}

type messageJSON struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

func (m *messageJSON) toEvent() *messageEvent {
	return &messageEvent{
		_name:     MESSAGE,
		timestamp: m.Timestamp,
		message:   m.Message,
	}
}

// client join room

type clientJoinRoomEvent struct {
	_name  eventType
	client *client
	room   *room
}

func makeClientJoinRoomEvent(client *client, room *room) *clientJoinRoomEvent {
	return &clientJoinRoomEvent{
		_name:  CLIENT_JOIN_ROOM,
		client: client,
		room:   room,
	}
}

func (e *clientJoinRoomEvent) name() eventType {
	return e._name
}

// client leave room

type clientLeaveRoomEvent struct {
	_name  eventType
	client *client
	room   *room
}

func makeClientLeaveRoomEvent(
	client *client,
	room *room,
) *clientLeaveRoomEvent {
	return &clientLeaveRoomEvent{
		_name:  CLIENT_LEAVE_ROOM,
		client: client,
		room:   room,
	}
}

func (e *clientLeaveRoomEvent) name() eventType {
	return e._name
}
