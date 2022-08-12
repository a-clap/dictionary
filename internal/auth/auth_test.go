//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package auth

import (
	"errors"
	"fmt"
	"testing"
)

type MemoryStoreError struct {
	store     MemoryStore
	returnErr bool
}

func (m *MemoryStoreError) Load(name string) (password string, err error) {
	if m.returnErr {
		return "", fmt.Errorf("internal error")
	}
	return m.store.Load(name)
}

func (m *MemoryStoreError) Save(name, password string) error {
	if m.returnErr {
		return fmt.Errorf("internal error")
	}
	return m.store.Save(name, password)
}

func (m *MemoryStoreError) NameExists(name string) (bool, error) {
	if m.returnErr {
		return false, fmt.Errorf("internal error")
	}
	return m.store.NameExists(name)
}

func (m *MemoryStoreError) Remove(name string) error {
	if m.returnErr {
		return fmt.Errorf("internal error")
	}
	return m.store.Remove(name)
}

func TestUsers_Add(t *testing.T) {
	// Table driven tests
	type fields struct {
		Store Store
	}
	type argsErr struct {
		name    string
		pass    string
		err     bool
		errType error
	}
	tests := []struct {
		name   string
		fields fields
		io     []argsErr
	}{
		{
			name:   "add single user",
			fields: fields{Store: &MemoryStore{store: map[string]string{}}},
			io: []argsErr{{
				name:    "adam",
				pass:    "password",
				err:     false,
				errType: nil,
			}},
		},
		{
			name: "add already existing user twice",
			fields: fields{Store: &MemoryStore{store: map[string]string{
				"adam": "password",
			}}},
			io: []argsErr{
				{
					name:    "adam",
					pass:    "password",
					err:     true,
					errType: ErrExist,
				},
			},
		},
		{
			name:   "invalid argument: password",
			fields: fields{Store: &MemoryStore{store: map[string]string{}}},
			io: []argsErr{
				{
					name:    "adam",
					pass:    "",
					err:     true,
					errType: ErrInvalid,
				},
			},
		},
		{
			name:   "invalid argument: name",
			fields: fields{Store: &MemoryStore{store: map[string]string{}}},
			io: []argsErr{
				{
					name:    "",
					pass:    "1",
					err:     true,
					errType: ErrInvalid,
				},
			},
		},
		{
			name:   "invalid argument: pass and name",
			fields: fields{Store: &MemoryStore{store: map[string]string{}}},
			io: []argsErr{
				{
					name:    "",
					pass:    "",
					err:     true,
					errType: ErrInvalid,
				},
			},
		},
		{
			name:   "handle internal IO error",
			fields: fields{&MemoryStoreError{store: MemoryStore{}, returnErr: true}},
			io: []argsErr{
				{
					name:    "adam",
					pass:    "password",
					err:     true,
					errType: ErrIO,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(tt.fields.Store)
			for _, v := range tt.io {
				err := u.Add(v.name, v.pass)
				if (err != nil) != v.err {
					t.Fatalf("%s: Add() error = %v, wantErr %v", tt.name, err, v.err)
				}
				// Check error type, if error needed
				if v.err {
					if !errors.Is(err, v.errType) {
						t.Errorf("%s: Add() error = %v, errType %v", tt.name, err, v.errType)
					}
				}
			}
		})
	}

	// Custom tests
	t.Run("add doesn't store passwords directly", func(t *testing.T) {
		mock := &MemoryStore{store: map[string]string{}}
		u := New(mock)

		name := "adam"
		password := "some crazy password"

		err := u.Add(name, password)
		if err != nil {
			t.Errorf("%s: Add() error %v unexpected", t.Name(), err)
		}
		// Naive compare
		if mock.store[name] == password {
			t.Errorf("%s: Add() saves plain password", t.Name())
		}
	})
}

func TestUsers_Remove(t *testing.T) {
	type fields struct {
		Store Store
	}
	type io struct {
		name    string
		err     bool
		errType error
	}
	tests := []struct {
		name   string
		fields fields
		args   []io
	}{
		{
			name: "handle io error",
			fields: fields{&MemoryStoreError{
				store:     MemoryStore{map[string]string{}},
				returnErr: true,
			}},
			args: []io{
				{
					name:    "adam",
					err:     true,
					errType: ErrIO,
				},
			},
		},
		{
			name:   "can't remove not existing user",
			fields: fields{Store: &MemoryStore{store: map[string]string{}}},
			args: []io{
				{
					name:    "not_exists",
					err:     true,
					errType: ErrNotExist,
				},
			},
		},
		{
			name: "remove existing user",
			fields: fields{Store: &MemoryStore{store: map[string]string{
				"adam": "pwd",
			}}},
			args: []io{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(tt.fields.Store)
			for _, v := range tt.args {
				err := u.Remove(v.name)
				if (err != nil) != v.err {
					t.Fatalf("%s: Remove() error = %v, wantErr %v", tt.name, err, v.err)
				}
				// Check error type, if error needed
				if v.err {
					if !errors.Is(err, v.errType) {
						t.Errorf("%s: Add() error = %v, errType %v", tt.name, err, v.errType)
					}
				}

			}
		})
	}
}

func TestUsers_Auth(t *testing.T) {
	type fields struct {
		LoadSaver Store
	}
	type io struct {
		name     string
		password string
		auth     bool
		err      bool
		errType  error
	}
	tests := []struct {
		name   string
		fields fields
		args   []io
	}{
		{
			name: "handle io error",
			fields: fields{LoadSaver: &MemoryStoreError{
				returnErr: true,
			}},
			args: []io{
				{
					name:     "dont matter",
					password: "also",
					auth:     false,
					err:      true,
					errType:  ErrIO,
				},
			},
		},
		{
			name: "unauthorized access",
			fields: fields{LoadSaver: &MemoryStore{
				store: map[string]string{
					"adam": "correct_pwd_but_not_hashed",
					"beta": "wrong_pwd",
				},
			}},
			args: []io{
				{
					name:     "adam",
					password: "correct_pwd_but_not_hashed",
					auth:     false,
					err:      false,
				},
				{
					name:     "beta",
					password: "other",
					auth:     false,
					err:      false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(tt.fields.LoadSaver)

			for _, args := range tt.args {
				auth, err := u.Auth(args.name, args.password)
				if (err != nil) != args.err {
					t.Errorf("%s: Auth() error = %#v, wantErr %v", tt.name, err, args.err)
				}
				// Check error type, if error needed
				if args.err {
					if !errors.Is(err, args.errType) {
						t.Errorf("%s: Add() error = %v, errType %v", tt.name, err, args.errType)
					}
				}
				if auth != args.auth {
					t.Errorf("%s: Auth() got %v, want %v", tt.name, auth, args.auth)
				}
			}

		})
	}

	t.Run("authorized access", func(t *testing.T) {
		//	Custom test - add user and then check authorized access
		m := &MemoryStore{
			store: map[string]string{},
		}
		u := New(m)
		name := "testing"
		password := "awesome password"
		if err := u.Add(name, password); err != nil {
			t.Errorf("%s: Add() unexpected error %#v", t.Name(), err)
		}

		if auth, err := u.Auth(name, password); err != nil {
			t.Errorf("%s: Auth() unexpected error %#v", t.Name(), err)
		} else if !auth {
			t.Errorf("%s: Auth() expected to authorize user %v", t.Name(), name)
		}
	})

}
