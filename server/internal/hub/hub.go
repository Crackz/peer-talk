package hub

import (
	"fmt"
	"peer-talk/internal/auth"
	"peer-talk/internal/common"
	"peer-talk/internal/member"
	"peer-talk/internal/socket_message"

	"github.com/olahol/melody"
)

type Hub struct {
	roomsManager   *roomsManager
	membersManager *membersManager
	Register       chan *SocketRegister
	Unregister     chan *melody.Session
	Message        chan *SocketMessage
	Signal         chan *SocketSignal
	PeerInit       chan *SocketPeerInit
}

func New() *Hub {
	return &Hub{
		roomsManager:   newRoomsManager(),
		membersManager: newMembersManager(),
		Register:       make(chan *SocketRegister),
		Unregister:     make(chan *melody.Session),
		Message:        make(chan *SocketMessage),
		Signal:         make(chan *SocketSignal),
		PeerInit:       make(chan *SocketPeerInit),
	}
}

func (h *Hub) getMemberFromSocket(socket *melody.Session) *member.Member {
	user := socket.MustGet(string(common.AuthenticatedUserContextKey)).(*auth.AuthenticatedUser)
	member := h.membersManager.getOne(user.Id)

	return member
}

func (h *Hub) createOrJoinRoomHandler(socket *melody.Session, joinRoomPayload *CreateOrJoinRoomPayload) error {
	user := socket.MustGet(string(common.AuthenticatedUserContextKey)).(*auth.AuthenticatedUser)
	member := h.membersManager.getOne(user.Id)

	if member != nil {
		if member.GetRoomId() == joinRoomPayload.RoomId {
			return nil
		}

		err := h.leaveRoomHandler(socket)
		if err != nil {
			return err
		}
	}

	room := h.roomsManager.getOrCreateOne(joinRoomPayload.RoomId)
	member = h.membersManager.createOne(socket, user, room.GetId())
	room.AddMember(member)
	err := room.BroadcastPeerInitMessage(member, socket_message.PeerReceiveInitType)
	if err != nil {
		return err
	}

	fmt.Println("JOINED ROOM : ", room.GetId())
	return nil

}

func (h *Hub) leaveRoomHandler(socket *melody.Session) error {
	member := h.getMemberFromSocket(socket)
	if member == nil {
		return nil
	}

	room := h.roomsManager.getOne(member.GetRoomId())
	if room == nil {
		return nil
	}

	err := room.BroadcastPeerInitMessage(member, socket_message.PeerCloseInitType)
	if err != nil {
		return err
	}

	room.RemoveMember(member)
	h.membersManager.deleteOne(member.Id)

	if room.MembersCount() > 0 {
		leftMsg := socket_message.NewLeftUserSocketMessage(member.Username)
		err := room.BroadcastTextMessage(leftMsg)
		if err != nil {
			return err
		}
	}

	if room.MembersCount() == 0 {
		h.roomsManager.deleteOne(room.GetId())
	}

	return nil
}

func (h *Hub) messageHandler(socket *melody.Session, msgText string) error {
	member := h.getMemberFromSocket(socket)
	if member == nil {
		return nil
	}

	room := h.roomsManager.getOne(member.GetRoomId())
	if room == nil {
		return nil
	}

	message := socket_message.NewSocketTextMessage(msgText, member.Username, socket_message.UserMessageType)

	return room.BroadcastTextMessage(message)
}

func (h *Hub) signalHandler(socket *melody.Session, signalPayload *SignalPayload) error {
	member := h.getMemberFromSocket(socket)
	if member == nil {
		return nil
	}

	room := h.roomsManager.getOne(member.GetRoomId())
	if room == nil {
		return nil
	}

	return room.SendPeerSignalMessage(member.Id, signalPayload.ReceiverId, signalPayload.Signal)

}

func (h *Hub) peerInitHandler(socket *melody.Session, receiverId string) error {
	member := h.getMemberFromSocket(socket)
	if member == nil {
		return nil
	}

	room := h.roomsManager.getOne(member.GetRoomId())
	if room == nil {
		return nil
	}

	return room.SendPeerInitMessage(member.Id, receiverId)

}

func (h *Hub) Run() {
	for {
		var err error
		select {
		case registerSocketMessage := <-h.Register:
			err = h.createOrJoinRoomHandler(registerSocketMessage.Socket, registerSocketMessage.JoinRoomPayload)
		case socket := <-h.Unregister:
			err = h.leaveRoomHandler(socket)
		case message := <-h.Message:
			err = h.messageHandler(message.Socket, message.Text)
		case signal := <-h.Signal:
			err = h.signalHandler(signal.Socket, signal.SignalPayload)
		case peerInit := <-h.PeerInit:
			err = h.peerInitHandler(peerInit.Socket, peerInit.ReceiverId)
		}
		if err != nil {
			fmt.Printf("unhandled error: %v\n", err)
		}
	}
}
