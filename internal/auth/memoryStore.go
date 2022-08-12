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
	key      []byte
	duration time.Duration
}

func (m *MemoryStore) Key() []byte {
	return m.key
}

func (m *MemoryStore) Duration() time.Duration {
	return m.duration
}

func (m *MemoryStore) Load(name string) ([]byte, error) {
	data, _ := m.store[name]
	return data, nil
}

func (m *MemoryStore) Save(name string, data []byte) error {
	m.store[name] = data
	return nil
}

func (m *MemoryStore) NameExists(name string) (bool, error) {
	_, ok := m.store[name]
	return ok, nil
}

func (m *MemoryStore) Remove(name string) error {
	delete(m.store, name)
	return nil
}

func NewMemoryStore(key []byte, duration time.Duration) *MemoryStore {
	return &MemoryStore{
		store:    map[string][]byte{},
		key:      key,
		duration: duration,
	}
}
