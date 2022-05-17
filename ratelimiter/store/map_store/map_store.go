package map_store

import (
	"sync"
	"time"
)

type Value struct {
	val            int64
	lastUpdateTime time.Time
}

type MapStore struct {
	data       map[string]Value
	sync       sync.RWMutex
	expireTime time.Duration
}

func NewMapStore(expireTime time.Duration) *MapStore {
	store := &MapStore{
		data:       make(map[string]Value),
		expireTime: expireTime,
	}
	return store
}

func (m *MapStore) Inc(key string, d time.Time) error {
	m.sync.Lock()
	defer m.sync.Unlock()

	data := m.data[getKey(key)]
	data.val++
	data.lastUpdateTime = time.Now().UTC()
	m.data[getKey(key)] = data
	return nil

	//panic("implement me")
}

func (m *MapStore) Get(key string, previousWindow, currentWindow time.Time) (prevValue int64, currValue int64, err error) {
	panic("implement me")
}
func (m *MapStore) UpdateToken() {
	panic("implement me")
}

func getKey(key string) string {
	return key

}
