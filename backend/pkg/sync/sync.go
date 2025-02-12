package sync

import (
	"sync"
)

var (
	userMutexes = make(map[int]*sync.Mutex)
	userMutexesLock = sync.Mutex{}
	GlobalDictLock = sync.Mutex{}
)

func GetUserMutex(userID int) *sync.Mutex {
	userMutexesLock.Lock()
	defer userMutexesLock.Unlock()

	if _, ok := userMutexes[userID]; !ok {
		userMutexes[userID] = &sync.Mutex{}
	}

	return userMutexes[userID]
}