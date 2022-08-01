package users_test

import (
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

func (l *LoadSaverMock) Exists(name string) (bool, error) {
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
		name string
		pass string
		err  bool
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
				name: "adam",
				pass: "password",
				err:  false,
			}},
		},
		{
			name:   "add same user twice",
			fields: fields{&LoadSaverMock{users: map[string]string{}}},
			io: []argsErr{
				{
					name: "adam",
					pass: "password",
					err:  false,
				},
				{
					name: "adam",
					pass: "password2",
					err:  false,
				},
			},
		},
		{
			name:   "handle internal IO error",
			fields: fields{&LoadSaverMock{users: map[string]string{}, returnErr: true}},
			io: []argsErr{
				{
					name: "adam",
					pass: "password",
					err:  true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := users.New(tt.fields.LoadSaver)
			for _, v := range tt.io {
				if err := u.Add(v.name, v.pass); (err != nil) != v.err {
					t.Errorf("%s: Add() error = %v, wantErr %v", tt.name, err, v.err)
				}
			}
		})
	}

	// Custom tests
	t.Run("add doesn't store passwords directly", func(t *testing.T) {
		mock := &LoadSaverMock{
			users:     make(map[string]string),
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
