package main

import (
	"sync"
)

//The second implementation.
type MemoryDataStore struct {
	*sync.RWMutex
	Users map[int64]string
}

func (mds *MemoryDataStore) Name() string {
	return "MemoryDataStore"
}

func (mds *MemoryDataStore) FindUserNameById(id int64) (string, error) {
	mds.RWMutex.RLock()
	defer mds.RWMutex.RUnlock()
	username, ok := mds.Users[id]
	if !ok {
		return "", UserNotFoundError
	}
	return username, nil
}

func NewMemoryDataStore(conf map[string]string) (DataStore, error) {
	return &MemoryDataStore{
		Users: map[int64]string{
			1: "mnbbrown",
			0: "root",
		},
		RWMutex: &sync.RWMutex{},
	}, nil
}
