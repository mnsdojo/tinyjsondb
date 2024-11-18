package internal

import "sync"

type LockManager struct {
	locks map[string]*sync.Mutex // master lock
	mu    sync.RWMutex
}

func NewLockManager() *LockManager {
	return &LockManager{
		locks: make(map[string]*sync.Mutex),
	}
}

func (lm *LockManager) AcquireLock(key string) {
	// lock the list of locks
	lm.mu.Lock()

	lock, exists := lm.locks[key]
	if !exists {
		lock = &sync.Mutex{}
		lm.locks[key] = lock
	}

	lm.mu.Unlock()
	lock.Lock()
}

// Releasing lock for specific key
func (lm *LockManager) ReleaseLock(key string) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	if lock, exists := lm.locks[key]; exists {
		lock.Unlock()
	}
}

func (lm *LockManager) DeleteLock(key string) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	delete(lm.locks, key)
}

func (lm *LockManager) GetLock(key string) *sync.Mutex {
	// Read lock because we're just looking, not touching
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	return lm.locks[key]
}
