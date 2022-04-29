package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/miaocansky/go-tool/cache"
	"time"
)

type MemcachedStore struct {
	cache.BaseCache
	client *memcache.Client
}

func NewMemcachedStore() *MemcachedStore {
	adress := "127.0.0.1:32768"
	client := memcache.New(adress)
	return &MemcachedStore{
		client: client,
	}

}

func (m *MemcachedStore) Close() error {
	m.client = nil
	return nil
}

func (m *MemcachedStore) Has(key string) bool {
	_, err := m.client.Get(key)
	return err == nil
}

func (m *MemcachedStore) Del(key string) error {
	return m.client.Delete(key)
}

func (m *MemcachedStore) Get(key string) interface{} {
	item, err := m.client.Get(key)
	return m.Unmarshal(item.Value, err)
}

func (m *MemcachedStore) Set(key string, val interface{}, ttl time.Duration) error {
	bytes, err := m.Marshal(val)
	if err != nil {
		return err
	}
	item := &memcache.Item{Key: key, Value: bytes, Expiration: int32(ttl / time.Second)}
	return m.client.Set(item)
}

func (m *MemcachedStore) SetHash(key string, val interface{}) error {
	panic("implement me")
}

func (m *MemcachedStore) GetHash(key, field string) interface{} {
	panic("implement me")
}
