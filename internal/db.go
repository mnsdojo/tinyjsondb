package internal

import (
	"sync"
)

type Database struct {
	data     map[string]interface{}
	mu       sync.RWMutex
	lm       *LockManager
	filePath string
}
