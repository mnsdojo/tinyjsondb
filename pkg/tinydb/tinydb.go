package tinydb

import (
	"fmt"
	"sync"
	"time"

	"github.com/mnsdojo/tinyjsondb/internal"
)

type TinyDB struct {
	db       *internal.Database
	cache    *internal.Cache
	mu       sync.RWMutex
	closed   bool
	cacheTTL time.Duration
}

// Option is a function type for configuring TinyDB
type Option func(*config)

// config holds congiguration parameters
type config struct {
	cacheTTL time.Duration
}

func WithCacheTTL(ttl time.Duration) Option {
	return func(c *config) {
		c.cacheTTL = ttl
	}
}

func NewTinyDB(filePath string, options ...Option) (*TinyDB, error) {
	// Default configuration
	config := &config{
		cacheTTL: 10 * time.Minute,
	}

	// Apply options
	for _, opt := range options {
		opt(config)
	}

	// Create internal database
	database, err := internal.NewDatabase(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize TinyDB: %w", err)
	}

	cache := internal.NewCache(config.cacheTTL)

	return &TinyDB{
		db:       database,
		cache:    cache,
		cacheTTL: config.cacheTTL,
	}, nil

}

func (tdb *TinyDB) Create(key string, value interface{}) error {
	tdb.mu.Lock()
	defer tdb.mu.Unlock()

	if tdb.closed {
		return fmt.Errorf("database is closed")
	}

	if err := tdb.db.Create(key, value); err != nil {
		return err

	}

	tdb.cache.Set(key, value)
	return nil
}

func (tdb *TinyDB) Read(key string) (interface{}, error) {
	tdb.mu.RLock()
	defer tdb.mu.RUnlock()
	if tdb.closed {
		return nil, fmt.Errorf("database is closed")
	}

	if cachedValue, found := tdb.cache.Get(key); found {
		return cachedValue, nil
	}

	//  get from the db

	value, err := tdb.db.Read(key)
	if err != nil {
		return nil, err
	}

	tdb.cache.Set(key, value)
	return value, nil

}

func (tdb *TinyDB) Update(key string, value interface{}) error {
	tdb.mu.Lock()
	defer tdb.mu.Unlock()

	if tdb.closed {
		return fmt.Errorf("database is closed")
	}

	// Check if key exists first
	_, err := tdb.db.Read(key)
	if err != nil {
		return fmt.Errorf("cannot update: %v", err)
	}

	// Update in database
	if err := tdb.db.Update(key, value); err != nil {
		return err
	}

	// Update cache
	tdb.cache.Set(key, value)
	return nil
}

func (tdb *TinyDB) Delete(key string) error {

	tdb.mu.Lock()
	defer tdb.mu.Unlock()
	if tdb.closed {
		return fmt.Errorf("database is closed")
	}

	// Delete from database
	if err := tdb.db.Delete(key); err != nil {
		return err
	}

	// Invalidate cache
	tdb.cache.Invalidate(key)
	return nil

}

func (tdb *TinyDB) ReadAll() map[string]interface{} {
	tdb.mu.RLock()
	defer tdb.mu.RUnlock()

	if tdb.closed {
		return nil
	}
	return tdb.db.ReadAll()
}

// Close database and cache
func (tdb *TinyDB) Close() error {
	tdb.mu.Lock()
	defer tdb.mu.Unlock()

	if tdb.closed {
		return nil
	}

	// Close database
	err := tdb.db.Close()
	tdb.closed = true
	return err
}
