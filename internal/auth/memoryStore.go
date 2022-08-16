//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package auth

import (
	"time"
)

var _ StoreTokener = &MemoryStore{}

// NewMemoryStore is default constructor for MemoryStore
func NewMemoryStore(key []byte, duration time.Duration) *MemoryStore {
	return &MemoryStore{
		store:    map[string][]byte{},
		key:      key,
		duration: duration,
	}
}

// MemoryStore satisfies Store interface
type MemoryStore struct {
	store    map[string][]byte
	key      []byte
	duration time.Duration
}

// Key is responsible for returning key to generate jwtToken
func (m *MemoryStore) Key() []byte {
	return m.key
}

// Duration returns token validation time
func (m *MemoryStore) Duration() time.Duration {
	return m.duration
}

// Load loads user data from store
func (m *MemoryStore) Load(name string) ([]byte, error) {
	data, _ := m.store[name]
	return data, nil
}

// Save users data into store
func (m *MemoryStore) Save(name string, data []byte) error {
	m.store[name] = data
	return nil
}

// NameExists returns true, whether user with provided name exists
func (m *MemoryStore) NameExists(name string) (bool, error) {
	_, ok := m.store[name]
	return ok, nil
}

// Remove user from store
func (m *MemoryStore) Remove(name string) error {
	delete(m.store, name)
	return nil
}
