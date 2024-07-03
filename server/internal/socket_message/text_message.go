package socket_message

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func NewSocketTextMessage(text string, username string, msgType MessageType) *SocketMessage[SocketTextMessagePayload] {
	return &SocketMessage[SocketTextMessagePayload]{
		ID:      uuid.New().String(),
		MsgType: msgType,
		Payload: SocketTextMessagePayload{
			Text:      text,
			Username:  username,
			CreatedAt: time.Now(),
		},
	}
}

func NewWelcomeSocketMessage(username string) *SocketMessage[SocketTextMessagePayload] {
	return NewSocketTextMessage(
		fmt.Sprintf("A New User %v Joined The Room", username),
		serverUsername,
		BotMessageType,
	)
}

func NewLeftUserSocketMessage(username string) *SocketMessage[SocketTextMessagePayload] {
	return NewSocketTextMessage(
		fmt.Sprintf("User: %v Left Room", username),
		serverUsername,
		BotMessageType,
	)
}
