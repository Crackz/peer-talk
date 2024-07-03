package socket_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"peer-talk/internal/auth"
	"peer-talk/internal/common"
	"peer-talk/internal/hub"

	"github.com/labstack/echo"
	"github.com/olahol/melody"
)

type validator interface {
	Validate(i any) error
	UnmarshalAndValidate(bytes []byte, payload any) (any, error)
}

type socketMessage struct {
	EventName string          `json:"eventName"`
	Payload   json.RawMessage `json:"payload"`
}

type socketHandler struct {
	wsHandler *melody.Melody
	hub       *hub.Hub
	validator validator
}

func NewSocketHandler(hub *hub.Hub, validator validator) *socketHandler {
	wsHandler := melody.New()
	wsHandler.Config.MaxMessageSize = 1024 * 1024
	return &socketHandler{
		wsHandler: wsHandler,
		hub:       hub,
		validator: validator,
	}
}

func (sh *socketHandler) RegisterRoute(r *echo.Group) {
	r.GET("/ws", func(c echo.Context) error {

		request := c.Request()
		requestCtx := request.Context()

		user := c.Get(string(common.AuthenticatedUserContextKey)).(*auth.AuthenticatedUser)
		ctx := context.WithValue(requestCtx, common.AuthenticatedUserContextKey, user)

		requestWithCtx := request.WithContext(ctx)

		return sh.wsHandler.HandleRequest(c.Response().Writer, requestWithCtx)
	})

	sh.handleEvents()
}

func (sh *socketHandler) handleEvents() {
	sh.wsHandler.HandleConnect(sh.connectHandler)
	sh.wsHandler.HandleDisconnect(sh.disconnectHandler)
	sh.wsHandler.HandleMessage(sh.messageHandler)
	sh.wsHandler.HandleError(func(s *melody.Session, err error) {
		fmt.Printf("Server Socket ERROR %v\n", err)

	})
}

func (sh *socketHandler) connectHandler(s *melody.Session) {
	ctx := s.Request.Context()
	authenticatedUser := ctx.Value(common.AuthenticatedUserContextKey).(*auth.AuthenticatedUser)

	s.Set(string(common.AuthenticatedUserContextKey), authenticatedUser)
}

func (sh *socketHandler) disconnectHandler(socket *melody.Session) {
	sh.hub.Unregister <- socket
}

func (sh *socketHandler) messageHandler(socket *melody.Session, msg []byte) {
	socketMsg := &socketMessage{}

	if err := json.Unmarshal(msg, socketMsg); err != nil {
		fmt.Printf("here %v\n", err)
		return
	}

	switch socketMsg.EventName {
	case "createOrJoinRoom":
		{
			sh.createOrJoinRoomHandler(socket, socketMsg.Payload)
		}
	case "message":
		{
			sh.textMessageHandler(socket, socketMsg.Payload)
		}
	case "signal":
		{
			sh.signalHandler(socket, socketMsg.Payload)
		}
	case "peerInit":
		{
			sh.peerInitHandler(socket, socketMsg.Payload)
		}
	default:
		{
			fmt.Println("Unknown Event Name: ", socketMsg.EventName)
		}
	}

}

func (sh *socketHandler) createOrJoinRoomHandler(socket *melody.Session, payload []byte) {
	createOrJoinRoomPayload, err := sh.validator.UnmarshalAndValidate(payload, &hub.CreateOrJoinRoomPayload{})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	socketRegister := &hub.SocketRegister{
		Socket:          socket,
		JoinRoomPayload: createOrJoinRoomPayload.(*hub.CreateOrJoinRoomPayload),
	}

	sh.hub.Register <- socketRegister
}

func (sh *socketHandler) textMessageHandler(socket *melody.Session, payload []byte) {
	messagePayload, err := sh.validator.UnmarshalAndValidate(payload, &hub.MessagePayload{})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	socketMessage := &hub.SocketMessage{
		Socket: socket,
		Text:   messagePayload.(*hub.MessagePayload).Text,
	}
	sh.hub.Message <- socketMessage
}

func (sh *socketHandler) signalHandler(socket *melody.Session, payload []byte) {
	signalPayload, err := sh.validator.UnmarshalAndValidate(payload, &hub.SignalPayload{})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	socketSignal := &hub.SocketSignal{
		Socket:        socket,
		SignalPayload: signalPayload.(*hub.SignalPayload),
	}

	sh.hub.Signal <- socketSignal
}

func (sh *socketHandler) peerInitHandler(socket *melody.Session, payload []byte) {
	peerInitPayload, err := sh.validator.UnmarshalAndValidate(payload, &hub.PeerInitPayload{})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	socketPeerInit := &hub.SocketPeerInit{
		Socket:     socket,
		ReceiverId: peerInitPayload.(*hub.PeerInitPayload).ReceiverId,
	}

	sh.hub.PeerInit <- socketPeerInit
}
