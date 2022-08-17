//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package auth

import (
	"time"
)

var _ StoreTokener = &MemoryStore{}

// MemoryStore satisfies Store interface
type MemoryStore struct {
	store    map[string][]byte
	tokens   []string
	key      []byte
	duration time.Duration
}

// NewMemoryStore is default constructor for MemoryStore
func NewMemoryStore(key []byte, duration time.Duration) *MemoryStore {
	return &MemoryStore{
		store:    map[string][]byte{},
		tokens:   []string{},
		key:      key,
		duration: duration,
	}
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

func (m *MemoryStore) AddToken(token string) error {
	m.tokens = append(m.tokens, token)
	return nil
}

func (m *MemoryStore) TokenExists(token string) (bool, error) {
	for _, elem := range m.tokens {
		if elem == token {
			return true, nil
		}
	}
	return false, nil
}

func (m *MemoryStore) RemoveToken(token string) error {
	for i, elem := range m.tokens {
		if elem == token {
			m.tokens[i] = m.tokens[len(m.tokens)-1]
			m.tokens = m.tokens[:len(m.tokens)-1]
			return nil
		}
	}
	return nil
}
