package lock

import "sync"

var RoomLocks sync.Map

func GetRoomLock(
	roomID uint,
) *sync.Mutex {

	lock, _ :=
		RoomLocks.LoadOrStore(
			roomID,
			&sync.Mutex{},
		)

	return lock.(*sync.Mutex)
}