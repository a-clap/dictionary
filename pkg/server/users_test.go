//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package server_test

import (
	"bytes"
	"fmt"
	"github.com/a-clap/dictionary/internal/auth"
	"github.com/a-clap/dictionary/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type UsersInterfaceTest struct {
	duration  time.Duration
	key       []byte
	store     *auth.MemoryStore
	returnErr bool
}

func (u *UsersInterfaceTest) TokenExpireTime() time.Duration {
	return u.duration
}
func (u *UsersInterfaceTest) TokenKey() []byte {
	return u.key
}

func (u *UsersInterfaceTest) Load(name string) (password string, err error) {
	if u.returnErr {
		return "", fmt.Errorf("internal error")
	}
	return u.store.Load(name)
}

func (u *UsersInterfaceTest) Save(name, password string) error {
	if u.returnErr {
		return fmt.Errorf("internal error")
	}
	return u.store.Save(name, password)
}

func (u *UsersInterfaceTest) NameExists(name string) (bool, error) {
	if u.returnErr {
		return false, fmt.Errorf("internal error")
	}
	return u.store.NameExists(name)
}

func (u *UsersInterfaceTest) Remove(name string) error {
	if u.returnErr {
		return fmt.Errorf("internal error")
	}
	return u.store.Remove(name)
}

func NewUserInterfaceTest(duration time.Duration, key string, err bool) *UsersInterfaceTest {
	return &UsersInterfaceTest{
		duration:  duration,
		key:       []byte(key),
		store:     auth.NewMemoryStore(),
		returnErr: err,
	}
}

func TestServer_addUser(t *testing.T) {
	type fields struct {
		h server.Handler
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
				h: NewUserInterfaceTest(1*time.Minute, "key", false),
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
				h: NewUserInterfaceTest(1*time.Minute, "key", false),
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
				h: &UsersInterfaceTest{
					store:     auth.NewMemoryStore(),
					returnErr: true,
				},
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
			s := server.New(tt.fields.h)

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
		key      []byte
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
				key:      []byte("key"),
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
				key:      []byte("key"),
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

			got, err := tt.args.u.Token(tt.args.duration, tt.args.key)

			if tt.token.err {
				require.NotNil(t, err, tt.name)
			} else {
				require.Nil(t, err, tt.name)
			}

			validated, err := tt.args.u.Validate(got, tt.args.key)

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

func TestServer_loginUser(t *testing.T) {
	type fields struct {
		h server.Handler
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
		name string
		in   in
		out  out
	}
	tests := []struct {
		name   string
		fields fields
		params []params
	}{
		{
			name: "add user, then login",
			fields: fields{
				h: NewUserInterfaceTest(1*time.Minute, "key", false),
			},
			params: []params{
				{
					name: "add user",
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
				{
					name: "login user - invalid credentials",
					in: in{
						url:    "/api/user/login",
						method: http.MethodPost,
						body:   `{"name": "adam", "password": "pwd1"}`,
					},
					out: out{
						code: http.StatusUnauthorized,
						body: "error",
					},
				},
				{
					name: "login user - invalid user",
					in: in{
						url:    "/api/user/login",
						method: http.MethodPost,
						body:   `{"name": "adam2", "password": "pwd"}`,
					},
					out: out{
						code: http.StatusInternalServerError,
						body: "error",
					},
				},
				{
					name: "login user - success",
					in: in{
						url:    "/api/user/login",
						method: http.MethodPost,
						body:   `{"name": "adam", "password": "pwd"}`,
					},
					out: out{
						code: http.StatusOK,
						body: "token",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := require.New(t)
			s := server.New(tt.fields.h)

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
