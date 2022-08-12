//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package auth

// MemoryStore satisfies Store interface
type MemoryStore struct {
	store map[string]string
}

func (m *MemoryStore) Load(name string) (string, error) {
	password, _ := m.store[name]
	return password, nil
}

func (m *MemoryStore) Save(name, password string) error {
	m.store[name] = password
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

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{store: map[string]string{}}
}
