//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package server_test

import (
	"bytes"
	"fmt"
	"github.com/a-clap/dictionary/internal/users"
	"github.com/a-clap/dictionary/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MemoryStoreError struct {
	store     *users.MemoryStore
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

func TestServer_addUser(t *testing.T) {
	type fields struct {
		u *users.Users
	}
	type in struct {
		url    string
		method string
		body   string
	}
	type out struct {
		code int
		body string
	}
	type params struct {
		in  in
		out out
	}
	tests := []struct {
		name   string
		fields fields
		params []params
	}{
		{
			name: "add user",
			fields: fields{
				u: users.New(users.NewMemoryStore()),
			},
			params: []params{
				{
					in: in{
						url:    "/api/user/add",
						method: http.MethodPost,
						body:   `{"name": "adam", "password": "pwd"}`,
					},
					out: out{
						code: http.StatusCreated,
						body: `{"name":"adam"}`,
					},
				},
			},
		},
		{
			name: "handle errors",
			fields: fields{
				u: users.New(users.NewMemoryStore()),
			},
			params: []params{
				{
					in: in{
						url:    "/api/user/add",
						method: http.MethodPost,
						body:   `hello world`,
					},
					out: out{
						code: http.StatusBadRequest,
						body: "error",
					},
				},
				{
					in: in{
						url:    "/api/user/add",
						method: http.MethodPost,
						body:   `{"name": "adam", "password": ""}`,
					},
					out: out{
						code: http.StatusInternalServerError,
						body: "error",
					},
				},
			},
		},
		{
			name: "handle IO error",
			fields: fields{
				u: users.New(&MemoryStoreError{
					store:     users.NewMemoryStore(),
					returnErr: true,
				}),
			},
			params: []params{
				{
					in: in{
						url:    "/api/user/add",
						method: http.MethodPost,
						body:   `{"name": "adam", "password": "pwd"}`,
					},
					out: out{
						code: http.StatusInternalServerError,
						body: `"error":`,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := require.New(t)
			s := server.New(tt.fields.u)

			for i, param := range tt.params {
				reader := bytes.NewBuffer([]byte(param.in.body))

				request, err := http.NewRequest(param.in.method, param.in.url, reader)

				req.Nil(err, "unexpected err on %v", i)

				response := httptest.NewRecorder()
				s.ServeHTTP(response, request)

				assert.Equal(t, param.out.code, response.Code)
				assert.Contains(t, response.Body.String(), param.out.body)
			}
		})
	}
}

func TestUserToken_Validate(t *testing.T) {
	type args struct {
		u        server.User
		duration time.Duration
	}
	type token struct {
		err     bool
		errType error
	}
	type validate struct {
		err       bool
		errType   error
		validated bool
	}
	tests := []struct {
		name     string
		args     args
		token    token
		validate validate
	}{
		{
			name: "validation test",
			args: args{
				u: server.User{
					Name:     "adam",
					Password: "",
				},
				duration: 3 * time.Second,
			},
			token: token{
				err: false,
			},
			validate: validate{
				err:       false,
				validated: true,
			},
		},
		{
			name: "should expire",
			args: args{
				u: server.User{
					Name:     "adam",
					Password: "",
				},
				duration: 1 * time.Microsecond,
			},
			token: token{
				err: false,
			},
			validate: validate{
				err:       true,
				errType:   server.ErrExpired,
				validated: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.args.u.Token(tt.args.duration)

			if tt.token.err {
				require.NotNil(t, err, tt.name)
			} else {
				require.Nil(t, err, tt.name)
			}

			validated, err := tt.args.u.Validate(got)

			if tt.validate.err {
				require.NotNil(t, err, tt.name)
				require.Equal(t, tt.validate.errType, err)
			} else {
				require.Nil(t, err, tt.name)
			}

			require.Equal(t, tt.validate.validated, validated)
		})
	}
}
