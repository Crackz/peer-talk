package hub

import "github.com/olahol/melody"

type CreateOrJoinRoomPayload struct {
	RoomId string `json:"roomId" validate:"required,notblank,min=1,max=1000"`
}

type MessagePayload struct {
	Text string `json:"text" validate:"required,notblank,min=1"`
}

type SignalPayload struct {
	Signal     any    `json:"signal" validate:"required"`
	ReceiverId string `json:"receiverId" validate:"required"`
}

type PeerInitPayload struct {
	ReceiverId string `json:"receiverId" validate:"required"`
}

type SocketRegister struct {
	Socket          *melody.Session
	JoinRoomPayload *CreateOrJoinRoomPayload
}

type SocketSignal struct {
	Socket        *melody.Session
	SignalPayload *SignalPayload
}

type SocketMessage struct {
	Socket *melody.Session
	Text   string
}

type SocketPeerInit struct {
	Socket     *melody.Session
	ReceiverId string
}
