package room

import (
	"encoding/json"
	"fmt"
	"peer-talk/internal/member"
	"peer-talk/internal/socket_message"

	"github.com/google/uuid"
)

type Room struct {
	id      string
	members map[string]*member.Member
}

func New(id string) *Room {
	return &Room{
		id:      id,
		members: make(map[string]*member.Member),
	}
}

func (r *Room) GetId() string {
	return r.id
}

func (r *Room) IsMember(memberId string) bool {
	_, isMemberExist := r.members[memberId]
	return isMemberExist
}

func (r *Room) AddMember(member *member.Member) {
	r.members[member.Id] = member

}

func (r *Room) RemoveMember(member *member.Member) {
	delete(r.members, member.Id)
}

func (r *Room) MembersCount() int {
	return len(r.members)
}

func (r *Room) broadcast(message any, exceptMember *member.Member) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// fmt.Printf("broadcasting message: %v\n", message)

	for _, member := range r.members {
		if exceptMember != nil && exceptMember.Id == member.Id {
			continue
		}
		member.WriteMessage(messageBytes)
	}

	return nil
}

func (r *Room) sendToMember(memberId string, message any) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// fmt.Printf("sending message: %v to member: %v\n", message, memberId)

	for _, member := range r.members {
		if member.Id == memberId {
			member.WriteMessage(messageBytes)
		}
	}

	return nil
}

func (r *Room) BroadcastTextMessage(textMessage *socket_message.SocketMessage[socket_message.SocketTextMessagePayload]) error {
	return r.broadcast(textMessage, nil)
}

func (r *Room) BroadcastPeerInitMessage(member *member.Member, initType socket_message.PeerInitType) error {
	initPeerMessage := &socket_message.SocketMessage[socket_message.SocketPeerInitMessagePayload]{
		ID:      uuid.New().String(),
		MsgType: socket_message.PeerInitMessageType,
		Payload: socket_message.SocketPeerInitMessagePayload{
			UserId:   member.Id,
			InitType: initType,
		},
	}

	return r.broadcast(initPeerMessage, member)
}

func (r *Room) SendPeerInitMessage(senderId string, receiverId string) error {
	signalMessage := &socket_message.SocketMessage[socket_message.SocketPeerInitMessagePayload]{
		ID:      uuid.New().String(),
		MsgType: socket_message.PeerInitMessageType,
		Payload: socket_message.SocketPeerInitMessagePayload{
			UserId:   senderId,
			InitType: socket_message.PeerSendInitType,
		},
	}

	return r.sendToMember(receiverId, signalMessage)
}

func (r *Room) SendPeerSignalMessage(senderId string, receiverId string, signal any) error {
	fmt.Printf("sending signal message from sender %v to receiver %v\n", senderId, receiverId)

	signalMessage := &socket_message.SocketMessage[socket_message.SocketPeerSignalMessagePayload]{
		ID:      uuid.New().String(),
		MsgType: socket_message.PeerSignalMessageType,
		Payload: socket_message.SocketPeerSignalMessagePayload{
			Signal: signal,
			UserId: senderId,
		},
	}

	return r.sendToMember(receiverId, signalMessage)
}
