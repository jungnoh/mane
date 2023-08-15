package mane

import "sync"

type ThreadSafeMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewThreadSafeMap[K comparable, V any]() *ThreadSafeMap[K, V] {
	return &ThreadSafeMap[K, V]{
		data: make(map[K]V),
	}
}

func (m *ThreadSafeMap[K, V]) Get(key K) (value V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok = m.data[key]
	return
}

func (m *ThreadSafeMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *ThreadSafeMap[K, V]) SetIfAbsent(key K, value func() V) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, loaded := m.data[key]; !loaded {
		m.data[key] = value()
		return m.data[key], true
	}
	return m.data[key], false
}

func (m *ThreadSafeMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

type LockMap struct {
	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

func NewLockMap() LockMap {
	return LockMap{
		locks: make(map[string]*sync.Mutex),
	}
}

func (m *LockMap) Lock(key string) {
	m.mu.Lock()
	if _, ok := m.locks[key]; !ok {
		m.locks[key] = &sync.Mutex{}
	}
	m.mu.Unlock()
	m.locks[key].Lock()
}

func (m *LockMap) Unlock(key string) {
	m.mu.Lock()
	if _, ok := m.locks[key]; !ok {
		m.mu.Unlock()
		return
	}
	m.mu.Unlock()
	m.locks[key].Unlock()
}
