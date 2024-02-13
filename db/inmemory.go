package db

import (
	"sync"
)

// InMemoryDB is a simple in-memory key-value store
var (
	inMemoryInstance *InMemoryDB
	onceI            sync.Once
)

// InMemoryDB is a simple in-memory key-value store
type InMemoryDB struct {
	data map[string]string
	mu   sync.RWMutex
}

// GetInMemoryDBInstance returns a singleton instance of InMemoryDB
func GetInMemoryDBInstance() *InMemoryDB {
	onceI.Do(func() {
		inMemoryInstance = NewInMemoryDB()
	})
	return inMemoryInstance
}

// NewInMemoryDB creates a new instance of InMemoryDB
func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		data: make(map[string]string),
	}
}

// Set sets the value for the given key in the in-memory database
func (db *InMemoryDB) Set(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}

// Get returns the value associated with the given key from the in-memory database
func (db *InMemoryDB) Get(key string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, ok := db.data[key]
	return val, ok
}
