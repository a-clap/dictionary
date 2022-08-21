//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/auth"
	"github.com/a-clap/dictionary/pkg/server"
	"github.com/a-clap/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	logger.Init(logger.NewDefaultZap(zapcore.DebugLevel))
}

var _ auth.StoreTokener = &memoryStoreError{}

type memoryStoreError struct {
	store *auth.MemoryStore
	err   bool
}

func (m *memoryStoreError) AddToken(token string) error {
	if m.err {
		return fmt.Errorf("io error")
	}
	return m.store.AddToken(token)
}

func (m *memoryStoreError) TokenExists(token string) (bool, error) {
	if m.err {
		return false, fmt.Errorf("io error")
	}
	return m.store.TokenExists(token)
}

func (m *memoryStoreError) RemoveToken(token string) error {
	if m.err {
		return fmt.Errorf("io error")
	}
	return m.store.RemoveToken(token)
}

func (m *memoryStoreError) Key() []byte {
	return m.store.Key()
}

func (m *memoryStoreError) Duration() time.Duration {
	return m.store.Duration()
}

// Load loads user data from store
func (m *memoryStoreError) Load(name string) ([]byte, error) {
	if m.err {
		return nil, fmt.Errorf("io error")
	}
	return m.store.Load(name)
}

// Save users data into store
func (m *memoryStoreError) Save(name string, data []byte) error {
	if m.err {
		return fmt.Errorf("io error")
	}
	return m.store.Save(name, data)
}

// NameExists returns true, whether user with provided name exists
func (m *memoryStoreError) NameExists(name string) (bool, error) {
	if m.err {
		return false, fmt.Errorf("io error")
	}
	return m.store.NameExists(name)
}

// Remove user from store
func (m *memoryStoreError) Remove(name string) error {
	if m.err {
		return fmt.Errorf("io error")
	}
	return m.store.Remove(name)
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
				h: auth.NewMemoryStore([]byte("extra private key"), 1*time.Minute),
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
				h: auth.NewMemoryStore([]byte("extra private key"), 1*time.Minute),
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
				h: &memoryStoreError{
					store: auth.NewMemoryStore([]byte("key"), 1*time.Hour),
					err:   true,
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
				h: auth.NewMemoryStore([]byte("key"), 1*time.Minute),
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

func TestServer_authUser(t *testing.T) {
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
			name: "add user, login, incorrect auth",
			fields: fields{
				h: auth.NewMemoryStore([]byte("key"), 1*time.Minute),
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
					name: "login user",
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
				{
					name: "use auth - wrong body",
					in: in{
						url:    "/api/translate/ping",
						method: http.MethodGet,
						body:   `{"name": "adam", "password": "pwd"}`,
					},
					out: out{
						code: http.StatusUnauthorized,
						body: "token",
					},
				},
				{
					name: "use auth - wrong token",
					in: in{
						url:    "/api/translate/ping",
						method: http.MethodGet,
						body:   `{"token": "1234"}`,
					},
					out: out{
						code: http.StatusUnauthorized,
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

				assert.Equal(t, param.out.code, response.Code, "test: %s", param.name)
				resp := response.Body.String()
				assert.Contains(t, resp, param.out.body, "test: %s", param.name)
			}
		})
	}

}

func TestServer_authUserToken(t *testing.T) {
	type addUsers struct {
		body string
	}

	users := []addUsers{{
		body: `{"name": "adam", "password": "pwd"}`,
	}}
	// Prepare server
	s := server.New(auth.NewMemoryStore([]byte("key"), time.Hour))

	t.Run("add users", func(t *testing.T) {
		for _, user := range users {
			reader := bytes.NewBuffer([]byte(user.body))
			request, err := http.NewRequest(http.MethodPost, "/api/user/add", reader)
			require.Nil(t, err)

			response := httptest.NewRecorder()
			s.ServeHTTP(response, request)
			require.Equal(t, http.StatusCreated, response.Code)
		}
	})

	t.Run("authorize with correct token", func(t *testing.T) {
		reader := bytes.NewBuffer([]byte(users[0].body))
		request, err := http.NewRequest(http.MethodPost, "/api/user/login", reader)
		require.Nil(t, err)

		response := httptest.NewRecorder()
		s.ServeHTTP(response, request)
		require.Equal(t, http.StatusOK, response.Code)

		resp := make(map[string]interface{})

		err = json.NewDecoder(response.Body).Decode(&resp)
		require.Nil(t, err)
		v, ok := resp["token"]
		require.True(t, ok)

		// Use v as token
		request, err = http.NewRequest(http.MethodGet, "/api/translate/ping", nil)
		request.Header.Set("Authorization", v.(string))
		require.Nil(t, err)
		s.ServeHTTP(response, request)
		require.Equal(t, http.StatusOK, response.Code)
		expected := `{"message":"pong"}`
		require.Equal(t, expected, response.Body.String())

	})
}
