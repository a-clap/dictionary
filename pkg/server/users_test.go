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
				u: users.New(&LoadSaverMock{
					users:     map[string]string{},
					returnErr: false,
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
						code: http.StatusCreated,
						body: `{"name":"adam"}`,
					},
				},
			},
		},
		{
			name: "handle errors",
			fields: fields{
				u: users.New(&LoadSaverMock{
					users:     map[string]string{},
					returnErr: false,
				}),
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
				u: users.New(&LoadSaverMock{
					users:     map[string]string{},
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
