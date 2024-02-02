package manager

import (
	"errors"
	"gossip/internal/room"
)

var (
	RoomExistsError       = errors.New("room already exists in manager")
	RoomDoesNotExistError = errors.New("room does not exist in manager")
)

type Manager struct {
	rooms map[string]*room.Room
}

func New() *Manager {
	return &Manager{
		rooms: make(map[string]*room.Room),
	}
}

func (m *Manager) AddRoom(name string) (*room.Room, error) {
	if _, ok := m.rooms[name]; ok {
		return nil, RoomExistsError
	}
	r := room.New(name)
	m.rooms[r.Name] = r
	return r, nil
}

func (m *Manager) Room(name string) (*room.Room, error) {
	room, ok := m.rooms[name]
	if !ok {
		return nil, RoomDoesNotExistError
	}
	return room, nil
}

func (m *Manager) Rooms() []string {
	rooms := make([]string, 0, len(m.rooms))
	for room := range m.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

func (m *Manager) RemoveRoom(name string) error {
	if _, ok := m.rooms[name]; !ok {
		return RoomDoesNotExistError
	}
	delete(m.rooms, name)
	return nil
}
