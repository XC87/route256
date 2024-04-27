package keymutex

import (
	"sync"
)

// корявенькая реализация, todo получше
type KeyRWMutex struct {
	mu    sync.Mutex
	mutex map[string]*sync.RWMutex
}

func NewKeyRWMutex() *KeyRWMutex {
	return &KeyRWMutex{
		mutex: map[string]*sync.RWMutex{},
	}
}

func (km *KeyRWMutex) RLock(key string) {
	km.mu.Lock()
	defer km.mu.Unlock()

	if km.mutex[key] == nil {
		km.mutex[key] = &sync.RWMutex{}
	}
	km.mutex[key].RLock()
}

func (km *KeyRWMutex) RUnlock(key string) {
	km.mu.Lock()
	defer km.mu.Unlock()

	if km.mutex[key] != nil {
		km.mutex[key].RUnlock()
	}
}

func (km *KeyRWMutex) Lock(key string) {
	km.mu.Lock()
	defer km.mu.Unlock()

	if km.mutex[key] == nil {
		km.mutex[key] = &sync.RWMutex{}
	}
	km.mutex[key].Lock()
}

func (km *KeyRWMutex) Unlock(key string) {
	km.mu.Lock()
	defer km.mu.Unlock()

	if km.mutex[key] != nil {
		km.mutex[key].Unlock()
	}
}
