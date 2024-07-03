package hub

import (
	"peer-talk/internal/auth"
	"peer-talk/internal/member"

	"github.com/olahol/melody"
)

type membersManager struct {
	members map[string]*member.Member
}

func newMembersManager() *membersManager {
	return &membersManager{
		members: make(map[string]*member.Member),
	}
}

func (mm *membersManager) getOne(id string) *member.Member {
	member, isMemberExist := mm.members[id]
	if !isMemberExist {
		return nil
	}

	return member
}

func (mm *membersManager) createOne(socket *melody.Session, user *auth.AuthenticatedUser, roomId string) *member.Member {
	newMember := member.New(socket, user, roomId)
	mm.members[user.Id] = newMember
	return newMember
}

func (mm *membersManager) deleteOne(id string) {
	delete(mm.members, id)
}
