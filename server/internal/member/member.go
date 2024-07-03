package member

import (
	"peer-talk/internal/auth"

	"github.com/olahol/melody"
)

type Member struct {
	Id       string
	Username string
	roomId   string
	Socket   *melody.Session
}

func New(socket *melody.Session, user *auth.AuthenticatedUser, roomId string) *Member {
	return &Member{
		Id:       user.Id,
		Username: user.Name,
		Socket:   socket,
		roomId:   roomId,
	}
}

func (m *Member) GetRoomId() string {
	return m.roomId
}

func (m *Member) WriteMessage(messageBytes []byte) {
	m.Socket.Write(messageBytes)
}
