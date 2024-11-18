package internal

import (
	"fmt"
	"os"
	"sync"
)

type Database struct {
	data     map[string]interface{}
	mu       sync.RWMutex
	lm       *LockManager
	filePath string
}

func NewDatabase(filePath string) (*Database, error) {
	db := &Database{
		data:     make(map[string]interface{}),
		lm:       NewLockManager(),
		filePath: filePath,
	}

	// Load existing data if file exists
	if err := Load(filePath, &db.data); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading database: %w", err)
	}
	return db, nil
}

// Create adds a new key-value pair to the database

func (db *Database) Create(key string, value interface{}) error {
	db.lm.AcquireLock(key)
	defer db.lm.ReleaseLock(key)
	db.mu.Lock()
	defer db.mu.Unlock()

	// check if key exists
	if _, exists := db.data[key]; exists {
		return fmt.Errorf("key '%s' already exists", key)
	}
	db.data[key] = value
	return db.saveToFile()

}

func (db *Database) Read(key string) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if value, exists := db.data[key]; exists {
		return value, nil
	}
	return nil, fmt.Errorf("key '%s' not found", key)

}

func (db *Database) saveToFile() error {
	if err := Save(db.filePath, db.data); err != nil {
		return fmt.Errorf("error saving database: %w", err)
	}
	return nil
}

func (db *Database) ReadAll() map[string]interface{} {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Create a copy of the data to prevent external modifications
	dataCpy := make(map[string]interface{}, len(db.data))
	for k, v := range db.data {
		dataCpy[k] = v
	}
	return dataCpy
}

func (db *Database) Update(key string, value interface{}) error {
	// Acquire locks
	db.lm.AcquireLock(key)
	defer db.lm.ReleaseLock(key)

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.data[key]; exists {
		return fmt.Errorf("key '%s' not found", key)
	}

	db.data[key] = value
	return db.saveToFile()
}

func (db *Database) Delete(key string) error {
	db.lm.AcquireLock(key)
	defer func() {
		db.lm.ReleaseLock(key)
		db.lm.DeleteLock(key) // Clean up the lock after deletion
	}()

	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if key exists
	if _, exists := db.data[key]; !exists {
		return fmt.Errorf("key '%s' not found", key)
	}
	delete(db.data, key)
	return db.saveToFile()
}

func (db *Database) Clear() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data = make(map[string]interface{})
	return db.saveToFile()
}

func (db *Database) Close() error {
	// Save final state
	if err := db.saveToFile(); err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}
	return nil
}
