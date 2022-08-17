//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var _ StoreTokener = &MemoryStoreError{}

type MemoryStoreError struct {
	store     *MemoryStore
	returnErr bool
}

func (m *MemoryStoreError) AddToken(token string) error {
	if m.returnErr {
		return fmt.Errorf("io err")
	}
	return m.store.AddToken(token)
}

func (m *MemoryStoreError) TokenExists(token string) (bool, error) {
	if m.returnErr {
		return false, fmt.Errorf("io err")
	}
	return m.store.TokenExists(token)
}

func (m *MemoryStoreError) RemoveToken(token string) error {
	if m.returnErr {
		return fmt.Errorf("io err")
	}
	return m.store.RemoveToken(token)
}

func (m *MemoryStoreError) Key() []byte {
	return []byte("super key")
}

func (m *MemoryStoreError) Duration() time.Duration {
	return 1 * time.Minute
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
		storeTokener StoreTokener
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
			fields: fields{storeTokener: &MemoryStore{store: map[string][]byte{}}},
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
			fields: fields{storeTokener: &MemoryStore{store: map[string][]byte{
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
			fields: fields{storeTokener: &MemoryStore{store: map[string][]byte{}}},
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
			fields: fields{storeTokener: &MemoryStore{store: map[string][]byte{}}},
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
			fields: fields{storeTokener: &MemoryStore{store: map[string][]byte{}}},
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
			fields: fields{&MemoryStoreError{store: NewMemoryStore([]byte("key"), time.Hour), returnErr: true}},
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
			u := New(tt.fields.storeTokener)
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
		Store StoreTokener
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
				store:     NewMemoryStore([]byte("key"), time.Hour),
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
		storeTokener StoreTokener
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
			fields: fields{storeTokener: &MemoryStoreError{
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
			fields: fields{storeTokener: &MemoryStore{
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
			u := New(tt.fields.storeTokener)

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

func TestManager_TokenValidateToken(t *testing.T) {
	type fields struct {
		i StoreTokener
	}
	type args struct {
		add           User
		token         User
		messWithToken bool
		newToken      string
	}
	type wants struct {
		addErr          bool
		addErrType      error
		tokenErr        bool
		tokenErrType    error
		validateErr     bool
		validateErrType error
		validate        User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wants  wants
	}{
		{
			name: "the right path",
			fields: fields{
				i: NewMemoryStore([]byte("key"), time.Hour),
			},
			args: args{
				add: User{
					Name:     "adam",
					Password: "pwd",
				},
				token: User{
					Name:     "adam",
					Password: "pwd",
				},
			},
			wants: wants{
				addErr:          false,
				addErrType:      nil,
				tokenErr:        false,
				tokenErrType:    nil,
				validateErr:     false,
				validateErrType: nil,
				validate: User{
					Name: "adam",
				},
			},
		},
		{
			name: "token expired",
			fields: fields{
				i: NewMemoryStore([]byte("key"), time.Microsecond),
			},
			args: args{
				add: User{
					Name:     "adam",
					Password: "pwd",
				},
				token: User{
					Name:     "adam",
					Password: "pwd",
				},
			},
			wants: wants{
				addErr:          false,
				addErrType:      nil,
				tokenErr:        false,
				tokenErrType:    nil,
				validateErr:     true,
				validateErrType: ErrExpired,
				validate:        User{},
			},
		},
		{
			name: "not existing user",
			fields: fields{
				i: NewMemoryStore([]byte("key"), time.Hour),
			},
			args: args{
				add: User{
					Name:     "adam",
					Password: "pwd",
				},
				token: User{
					Name: "hehe",
				},
			},
			wants: wants{
				addErr:          false,
				addErrType:      nil,
				tokenErr:        true,
				tokenErrType:    ErrNotExist,
				validateErr:     false,
				validateErrType: nil,
				validate:        User{},
			},
		},
		{
			name: "invalid credentials user",
			fields: fields{
				i: NewMemoryStore([]byte("key"), time.Hour),
			},
			args: args{
				add: User{
					Name:     "adam",
					Password: "pwd",
				},
				token: User{
					Name:     "adam",
					Password: "pwd2",
				},
			},
			wants: wants{
				addErr:          false,
				addErrType:      nil,
				tokenErr:        true,
				tokenErrType:    ErrInvalidCredentials,
				validateErr:     false,
				validateErrType: nil,
				validate:        User{},
			},
		},
		{
			name: "mess with token",
			fields: fields{
				i: NewMemoryStore([]byte("key"), time.Hour),
			},
			args: args{
				add: User{
					Name:     "adam",
					Password: "pwd",
				},
				token: User{
					Name:     "adam",
					Password: "pwd",
				},
				messWithToken: true,
				newToken:      "blabla",
			},
			wants: wants{
				addErr:          false,
				addErrType:      nil,
				tokenErr:        false,
				tokenErrType:    nil,
				validateErr:     true,
				validateErrType: jwt.NewValidationError("token contains an invalid number of segments", jwt.ValidationErrorMalformed),
				validate:        User{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(tt.fields.i)

			// First, need to add user
			err := u.Add(tt.args.add)
			if tt.wants.addErr {
				require.NotNil(t, err)
				require.Equal(t, err, tt.wants.addErrType)
				return
			}
			require.Nil(t, err)

			// Then generate token
			got, err := u.Token(tt.args.token)
			if tt.wants.tokenErr {
				require.NotNil(t, err)
				require.True(t, errors.Is(err, tt.wants.tokenErrType))
				return
			}
			require.Nil(t, err)
			require.NotEmpty(t, got)

			if tt.args.messWithToken {
				got = tt.args.newToken
			}

			// Try to validate token
			user, err := u.ValidateToken(got)
			if tt.wants.validateErr {
				require.NotNil(t, err)
				require.Equal(t, err, tt.wants.validateErrType)
				return
			}

			require.Nil(t, err)
			require.Equal(t, tt.wants.validate.Name, user.Name)

		})
	}
}

func TestManager_Logout(t *testing.T) {
	type fields struct {
		i StoreTokener
	}
	type args struct {
		loginUser *User
		token     string
	}
	type wants struct {
		logout  User
		err     bool
		errType error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wants  wants
	}{
		{
			name: "the right path",
			fields: fields{
				i: NewMemoryStore([]byte("super secret key"), time.Hour),
			},
			args: args{
				loginUser: &User{
					Name:     "adam",
					Password: "pwd",
				},
				token: "",
			},
			wants: wants{
				logout: User{
					Name:     "adam",
					Password: "",
				},
				err:     false,
				errType: nil,
			},
		},
		{
			name: "wrong token",
			fields: fields{
				i: NewMemoryStore([]byte("super secret key"), time.Hour),
			},
			args: args{
				loginUser: nil,
				token:     "",
			},
			wants: wants{
				logout:  User{},
				err:     true,
				errType: jwt.NewValidationError("token contains an invalid number of segments", jwt.ValidationErrorMalformed),
			},
		},
		{
			name: "wrong token #2",
			fields: fields{
				i: NewMemoryStore([]byte("super secret key"), time.Hour),
			},
			args: args{
				loginUser: nil,
				token:     "123asfasb543rqsa",
			},
			wants: wants{
				logout:  User{},
				err:     true,
				errType: jwt.NewValidationError("token contains an invalid number of segments", jwt.ValidationErrorMalformed),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.fields.i)

			if tt.args.loginUser != nil {
				if err := m.Add(*tt.args.loginUser); err != nil {
					t.Fatalf("%s: err not expected %#v", t.Name(), err)
				}
				auth, err := m.Auth(*tt.args.loginUser)
				if err != nil {
					t.Fatalf("%s: err not expected %#v", t.Name(), err)
				}
				if !auth {
					t.Fatalf("%s: auth expected %#v", t.Name(), auth)
				}

				tt.args.token, err = m.Token(*tt.args.loginUser)
				if err != nil {
					t.Fatalf("%s: err not expected %#v", t.Name(), err)
				}
			}

			got, err := m.Logout(tt.args.token)
			if (err != nil) != tt.wants.err {
				t.Errorf("%s: Logout() error = %#v, tt.wants.err %#v", t.Name(), err, tt.wants.err)
				return
			}

			if tt.wants.err {
				require.Equal(t, tt.wants.errType, err)
				return
			}

			require.NotNil(t, got)
			require.Equal(t, tt.wants.logout.Name, got.Name)
			require.Empty(t, tt.wants.logout.Password)
		})
	}
}

func TestManager_LoginLogout(t *testing.T) {
	t.Run("logout twice", func(t *testing.T) {
		m := New(NewMemoryStore([]byte("key"), time.Hour))

		addUser := User{
			Name:     "adam",
			Password: "pwd",
		}

		// First, need to add user
		require.Nil(t, m.Add(addUser))

		// Get token for user
		token, err := m.Token(addUser)
		require.Nil(t, err)
		require.NotEmpty(t, token)

		// Check whether token is right
		validateToken, err := m.ValidateToken(token)
		require.Nil(t, err)
		require.Equal(t, validateToken.Name, addUser.Name)

		// Logout user
		logout, err := m.Logout(token)
		require.Nil(t, err)
		require.Equal(t, logout.Name, addUser.Name)

		// Logout for second time
		logout, err = m.Logout(token)
		require.NotNil(t, err)
		require.Nil(t, logout)

	})

}
