package chat

import (
	"gossip/internal/domains/user"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
)

type chatUser struct {
	service *service
	ingress chan event
	user    *user.User
	roomIds map[uuid.UUID]bool
	conn    *websocket.Conn
}

func (u *chatUser) disconnect() error {
	close(u.ingress)
	return u.conn.Close()
}
