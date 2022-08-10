//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package users_test

import (
	"errors"
	"fmt"
	"github.com/a-clap/dictionary/internal/users"
	"testing"
)

type LoadSaverMock struct {
	users     map[string]string
	returnErr bool
}

func (l *LoadSaverMock) Load(name string) (password string, err error) {
	if l.returnErr {
		return "", fmt.Errorf("internal error")
	}

	password, _ = l.users[name]
	return
}

func (l *LoadSaverMock) Save(name, password string) error {
	if l.returnErr {
		return fmt.Errorf("internal error")
	}

	l.users[name] = password
	return nil
}

func (l *LoadSaverMock) NameExists(name string) (bool, error) {
	if l.returnErr {
		return false, fmt.Errorf("internal error")
	}
	_, ok := l.users[name]
	return ok, nil
}

func (l *LoadSaverMock) Remove(name string) error {
	if l.returnErr {
		return fmt.Errorf("internal error")
	}
	delete(l.users, name)
	return nil
}

func TestUsers_Add(t *testing.T) {
	// Table driven tests
	type fields struct {
		LoadSaver users.LoadSaver
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
			fields: fields{&LoadSaverMock{users: map[string]string{}}},
			io: []argsErr{{
				name:    "adam",
				pass:    "password",
				err:     false,
				errType: nil,
			}},
		},
		{
			name: "add already existing user twice",
			fields: fields{&LoadSaverMock{users: map[string]string{
				"adam": "password",
			}}},
			io: []argsErr{
				{
					name:    "adam",
					pass:    "password",
					err:     true,
					errType: users.ErrExist,
				},
			},
		},
		{
			name:   "invalid argument: password",
			fields: fields{&LoadSaverMock{users: map[string]string{}}},
			io: []argsErr{
				{
					name:    "adam",
					pass:    "",
					err:     true,
					errType: users.ErrInvalid,
				},
			},
		},
		{
			name:   "invalid argument: name",
			fields: fields{&LoadSaverMock{users: map[string]string{}}},
			io: []argsErr{
				{
					name:    "",
					pass:    "1",
					err:     true,
					errType: users.ErrInvalid,
				},
			},
		},
		{
			name:   "invalid argument: pass and name",
			fields: fields{&LoadSaverMock{users: map[string]string{}}},
			io: []argsErr{
				{
					name:    "",
					pass:    "",
					err:     true,
					errType: users.ErrInvalid,
				},
			},
		},
		{
			name:   "handle internal IO error",
			fields: fields{&LoadSaverMock{users: map[string]string{}, returnErr: true}},
			io: []argsErr{
				{
					name:    "adam",
					pass:    "password",
					err:     true,
					errType: users.ErrIO,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := users.New(tt.fields.LoadSaver)
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
		mock := &LoadSaverMock{
			users:     map[string]string{},
			returnErr: false,
		}
		u := users.New(mock)

		name := "adam"
		password := "some crazy password"

		err := u.Add(name, password)
		if err != nil {
			t.Errorf("%s: Add() error %v unexpected", t.Name(), err)
		}
		// Naive compare
		if mock.users[name] == password {
			t.Errorf("%s: Add() saves plain password", t.Name())
		}

	})
}

func TestUsers_Remove(t *testing.T) {
	type fields struct {
		LoadSaver users.LoadSaver
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
			fields: fields{&LoadSaverMock{
				users:     make(map[string]string),
				returnErr: true,
			}},
			args: []io{
				{
					name:    "adam",
					err:     true,
					errType: users.ErrIO,
				},
			},
		},
		{
			name: "can't remove not existing user",
			fields: fields{LoadSaver: &LoadSaverMock{
				users:     make(map[string]string),
				returnErr: false,
			}},
			args: []io{
				{
					name:    "not_exists",
					err:     true,
					errType: users.ErrNotExist,
				},
			},
		},
		{
			name: "remove existing user",
			fields: fields{LoadSaver: &LoadSaverMock{
				users: map[string]string{
					"adam": "pwd",
				},
				returnErr: false,
			}},
			args: []io{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := users.New(tt.fields.LoadSaver)
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
		LoadSaver users.LoadSaver
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
			fields: fields{LoadSaver: &LoadSaverMock{
				users:     map[string]string{},
				returnErr: true,
			}},
			args: []io{
				{
					name:     "dont matter",
					password: "also",
					auth:     false,
					err:      true,
					errType:  users.ErrIO,
				},
			},
		},
		{
			name: "unauthorized access",
			fields: fields{LoadSaver: &LoadSaverMock{
				users: map[string]string{
					"adam": "correct_pwd_but_not_hashed",
					"beta": "wrong_pwd",
				},
				returnErr: false,
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
			u := users.New(tt.fields.LoadSaver)

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
		m := &LoadSaverMock{
			users:     map[string]string{},
			returnErr: false,
		}
		u := users.New(m)
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
