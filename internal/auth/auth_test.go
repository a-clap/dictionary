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

func (m *MemoryStoreError) Load(name string) (data []byte, err error) {
	if m.returnErr {
		return nil, fmt.Errorf("internal error")
	}
	return m.store.Load(name)
}

func (m *MemoryStoreError) Save(name string, data []byte) error {
	if m.returnErr {
		return fmt.Errorf("internal error")
	}
	return m.store.Save(name, data)
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
	type args struct {
		user    User
		err     bool
		errType error
	}
	tests := []struct {
		name   string
		fields fields
		io     []args
	}{
		{
			name:   "add single user",
			fields: fields{Store: &MemoryStore{store: map[string][]byte{}}},
			io: []args{{
				user: User{
					Name:     "adam",
					Password: "password",
				},
				err:     false,
				errType: nil,
			}},
		},
		{
			name: "add already existing user twice",
			fields: fields{Store: &MemoryStore{store: map[string][]byte{
				"adam": []byte("password"),
			}}},
			io: []args{
				{
					user: User{
						Name:     "adam",
						Password: "password",
					},
					err:     true,
					errType: ErrExist,
				},
			},
		},
		{
			name:   "invalid argument: password",
			fields: fields{Store: &MemoryStore{store: map[string][]byte{}}},
			io: []args{
				{
					user: User{
						Name:     "adam",
						Password: "",
					},
					err:     true,
					errType: ErrInvalid,
				},
			},
		},
		{
			name:   "invalid argument: name",
			fields: fields{Store: &MemoryStore{store: map[string][]byte{}}},
			io: []args{
				{
					user: User{
						Name:     "",
						Password: "1",
					},
					err:     true,
					errType: ErrInvalid,
				},
			},
		},
		{
			name:   "invalid argument: pass and name",
			fields: fields{Store: &MemoryStore{store: map[string][]byte{}}},
			io: []args{
				{
					user: User{
						Name:     "",
						Password: "",
					},
					err:     true,
					errType: ErrInvalid,
				},
			},
		},
		{
			name:   "handle internal IO error",
			fields: fields{&MemoryStoreError{store: MemoryStore{}, returnErr: true}},
			io: []args{
				{
					user: User{
						Name:     "adam",
						Password: "password",
					},
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
				err := u.Add(v.user)
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
		mock := &MemoryStore{store: map[string][]byte{}}
		u := New(mock)

		user := User{
			Name:     "adam",
			Password: "some crazy password",
		}

		err := u.Add(user)
		if err != nil {
			t.Errorf("%s: Add() error %v unexpected", t.Name(), err)
		}
		// Naive compare
		if string(mock.store[user.Name]) == user.Password {
			t.Errorf("%s: Add() saves plain password", t.Name())
		}
	})
}

func TestUsers_Remove(t *testing.T) {
	type fields struct {
		Store Store
	}
	type io struct {
		user    User
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
				store:     MemoryStore{map[string][]byte{}},
				returnErr: true,
			}},
			args: []io{
				{
					user: User{
						Name: "adam",
					},
					err:     true,
					errType: ErrIO,
				},
			},
		},
		{
			name:   "can't remove not existing user",
			fields: fields{Store: &MemoryStore{store: map[string][]byte{}}},
			args: []io{
				{
					user: User{
						Name: "not exists",
					},
					err:     true,
					errType: ErrNotExist,
				},
			},
		},
		{
			name: "remove existing user",
			fields: fields{Store: &MemoryStore{store: map[string][]byte{
				"adam": []byte("pwd"),
			}}},
			args: []io{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(tt.fields.Store)
			for _, v := range tt.args {
				err := u.Remove(v.user)
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
		user    User
		auth    bool
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
			fields: fields{LoadSaver: &MemoryStoreError{
				returnErr: true,
			}},
			args: []io{
				{
					user: User{
						Name:     "dont matter",
						Password: "also",
					},
					auth:    false,
					err:     true,
					errType: ErrIO,
				},
			},
		},
		{
			name: "unauthorized access",
			fields: fields{LoadSaver: &MemoryStore{
				store: map[string][]byte{
					"adam": []byte("correct_pwd_but_not_hashed"),
					"beta": []byte("wrong_pwd"),
				},
			}},
			args: []io{
				{
					user: User{
						Name:     "adam",
						Password: "correct_pwd_but_not_hashed",
					},
					auth: false,
					err:  false,
				},
				{
					user: User{
						Name:     "beta",
						Password: "other",
					},
					auth: false,
					err:  false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(tt.fields.LoadSaver)

			for _, args := range tt.args {
				auth, err := u.Auth(args.user)
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
			store: map[string][]byte{},
		}
		u := New(m)

		user := User{
			Name:     "testing",
			Password: "awesome password",
		}

		if err := u.Add(user); err != nil {
			t.Errorf("%s: Add() unexpected error %#v", t.Name(), err)
		}

		if auth, err := u.Auth(user); err != nil {
			t.Errorf("%s: Auth() unexpected error %#v", t.Name(), err)
		} else if !auth {
			t.Errorf("%s: Auth() expected to authorize user %v", t.Name(), user.Name)
		}
	})

}
