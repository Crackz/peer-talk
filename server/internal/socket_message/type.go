package socket_message

import "time"

type MessageType string
type PeerInitType string

const (
	BotMessageType        MessageType  = "BOT_MESSAGE"
	UserMessageType       MessageType  = "USER_MESSAGE"
	PeerInitMessageType   MessageType  = "PEER_INIT_MESSAGE"
	PeerSignalMessageType MessageType  = "PEER_SIGNAL_MESSAGE"
	PeerReceiveInitType   PeerInitType = "PEER_RECEIVE_INIT_TYPE"
	PeerSendInitType      PeerInitType = "PEER_SEND_INIT_TYPE"
	PeerCloseInitType     PeerInitType = "PEER_CLOSE_INIT_TYPE"
	serverUsername                     = "SERVER"
)

type SocketMessage[T any] struct {
	ID      string      `json:"id"`
	MsgType MessageType `json:"type"`
	Payload T           `json:"payload"`
}

type SocketTextMessagePayload struct {
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

type SocketPeerInitMessagePayload struct {
	UserId   string       `json:"userId"`
	InitType PeerInitType `json:"initType"`
}

type SocketPeerSignalMessagePayload struct {
	Signal any    `json:"signal"`
	UserId string `json:"userId"`
}
