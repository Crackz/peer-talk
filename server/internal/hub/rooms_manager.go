package hub

import "peer-talk/internal/room"

type roomsManager struct {
	rooms map[string]*room.Room
}

func newRoomsManager() *roomsManager {
	return &roomsManager{
		rooms: make(map[string]*room.Room),
	}
}

func (rm *roomsManager) getOne(id string) *room.Room {
	room, isRoomExist := rm.rooms[id]
	if !isRoomExist {
		return nil
	}

	return room
}

func (rm *roomsManager) createOne(id string) *room.Room {
	newRoom := room.New(id)
	rm.rooms[id] = newRoom
	return newRoom

}

func (rm *roomsManager) getOrCreateOne(id string) *room.Room {
	room, isRoomExist := rm.rooms[id]
	if !isRoomExist {
		return rm.createOne(id)
	}

	return room
}

func (rm *roomsManager) deleteOne(id string) {
	delete(rm.rooms, id)
}
