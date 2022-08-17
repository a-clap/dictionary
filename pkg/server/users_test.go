//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package server_test

import (
	"bytes"
	"fmt"
	"github.com/a-clap/dictionary/internal/auth"
	"github.com/a-clap/dictionary/internal/logger"
	"github.com/a-clap/dictionary/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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
		h      server.Handler
		logger logger.Logger
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
				h:      auth.NewMemoryStore([]byte("extra private key"), 1*time.Minute),
				logger: logger.NewDummy(),
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
				h:      auth.NewMemoryStore([]byte("extra private key"), 1*time.Minute),
				logger: logger.NewDevelopment(),
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
			s := server.New(tt.fields.h, tt.fields.logger)

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
		h      server.Handler
		logger logger.Logger
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
				h:      auth.NewMemoryStore([]byte("key"), 1*time.Minute),
				logger: logger.NewDummy(),
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
			s := server.New(tt.fields.h, tt.fields.logger)

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
