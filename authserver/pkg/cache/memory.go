package cache

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type cacheItem struct {
	value      string
	expiration time.Time
}

type memoryCacheService struct {
	data   map[string]cacheItem
	mutex  sync.RWMutex
	logger *slog.Logger
}

func (m *memoryCacheService) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	expiration := time.Now().Add(ttl)
	m.data[key] = cacheItem{
		value:      value,
		expiration: expiration,
	}

	// Clean up expired items periodically
	go m.cleanup()

	return nil
}

func (m *memoryCacheService) Get(ctx context.Context, key string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	item, exists := m.data[key]
	if !exists {
		return "", nil
	}

	if time.Now().After(item.expiration) {
		// Item expired, clean it up
		delete(m.data, key)
		return "", nil
	}

	return item.value, nil
}

func (m *memoryCacheService) Delete(ctx context.Context, key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.data, key)
	return nil
}

func (m *memoryCacheService) cleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	for key, item := range m.data {
		if now.After(item.expiration) {
			delete(m.data, key)
		}
	}
}
