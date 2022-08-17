//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Manager struct {
	i StoreTokener
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	claims   struct {
		Name string `json:"name"`
		jwt.RegisteredClaims
	}
}

var (
	ErrExist              = errors.New("user already exists")
	ErrNotExist           = errors.New("user doesn't exist")
	ErrInvalid            = errors.New("invalid argument")
	ErrIO                 = errors.New("io error")
	ErrHash               = errors.New("hash error") // tried to generate this error during tests, didn't happen
	ErrExpired            = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrBlacklisted        = errors.New("user blacklisted - logged out")
)

type (
	// Store realizes access to some kind of user database (maybe even just map[string][]data)
	// Errors returned by interface should be ONLY related to internal IO errors
	Store interface {
		// Load user data, if user doesn't exist, return ""
		Load(name string) (data []byte, err error)
		// Save new user with custom data. Overwrites, if already exists
		Save(name string, data []byte) error
		// NameExists allows to check whether user with name already exists
		NameExists(name string) (bool, error)
		// Remove user with provided name, if user doesn't exist, don't do anything
		Remove(name string) error
		// AddToken adds token to token blacklists
		AddToken(token string) error
		// TokenExists checks whether token exists in blacklist
		TokenExists(token string) (bool, error)
		// RemoveToken removes token from blacklist
		RemoveToken(token string) error
	}

	// Tokener realizes access to:
	// key via Key() used to generate token
	// duration via Duration() for generated jwt token
	Tokener interface {
		Key() []byte
		Duration() time.Duration
	}

	StoreTokener interface {
		Store
		Tokener
	}
)

// New is default constructor for Manager
func New(storeTokener StoreTokener) *Manager {
	return &Manager{i: storeTokener}
}

// Add adds user to base
func (m *Manager) Add(user User) error {
	if exists, err := m.Exists(user); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("%w %s", ErrExist, user.Name)
	}

	if len(user.Name) == 0 || len(user.Password) == 0 {
		return fmt.Errorf("%w: name and password must be provided", ErrInvalid)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%w %v", ErrHash, err)
	}

	if err := m.save(user.Name, hashedPassword); err != nil {
		return err
	}

	return nil
}

// Remove user from Manager
func (m *Manager) Remove(user User) error {
	if exists, err := m.Exists(user); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("%w %s", ErrNotExist, user.Name)
	}

	return m.i.Remove(user.Name)
}

// Exists checks, whether particular user exists
func (m *Manager) Exists(user User) (bool, error) {
	return m.nameExists(user.Name)
}

// Auth serves as user authentication (login)
func (m *Manager) Auth(user User) (bool, error) {
	if exists, err := m.Exists(user); err != nil {
		return false, err
	} else if !exists {
		return false, fmt.Errorf("%w %s", ErrNotExist, user.Name)
	}

	if hashPass, err := m.load(user.Name); err != nil {
		return false, err
	} else {
		return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(user.Password)) == nil, nil
	}
}

func (m *Manager) Logout(token string) (*User, error) {
	user, err := m.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	err = m.addToken(token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Token returns jwtToken, if user provides correct credentials
func (m *Manager) Token(user User) (string, error) {
	if auth, err := m.Auth(user); err != nil {
		return "", err
	} else if !auth {
		return "", ErrInvalidCredentials
	}

	expires := time.Now().Add(m.i.Duration())
	user.claims.Name = user.Name
	user.claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: expires},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user.claims)
	return token.SignedString(m.i.Key())
}

// ValidateToken returns associated User to token, error on invalid token
func (m *Manager) ValidateToken(token string) (*User, error) {
	var user User
	tkn, err := jwt.ParseWithClaims(token, &user.claims, func(token *jwt.Token) (interface{}, error) {
		return m.i.Key(), nil
	})

	if err != nil {
		if validationError, ok := err.(*jwt.ValidationError); ok {
			if (validationError.Errors & jwt.ValidationErrorExpired) == jwt.ValidationErrorExpired {
				return nil, ErrExpired
			}
		}
		return nil, err
	}

	if !tkn.Valid {
		return nil, ErrInvalidToken
	}

	if blacklisted, err := m.tokenExists(token); err != nil {
		return nil, err
	} else if blacklisted {
		return nil, ErrBlacklisted
	}

	user.Name = user.claims.Name

	return &user, nil
}

// load is wrapper for interface call Load, returns appropriate wrapped error
func (m *Manager) load(name string) (data []byte, err error) {
	data, err = m.i.Load(name)
	if err != nil {
		return nil, fmt.Errorf("%w: Load: name %s, error: %v", ErrIO, name, err)
	}
	return data, err
}

// save is wrapper for interface call Save, returns appropriate wrapped error
func (m *Manager) save(name string, data []byte) error {
	if err := m.i.Save(name, data); err != nil {
		return fmt.Errorf("%w: Save: name %s, error: %v", ErrIO, name, err)
	}
	return nil
}

// nameExists is wrapper for interface call NameExists, returns appropriate wrapped error
func (m *Manager) nameExists(name string) (bool, error) {
	exist, err := m.i.NameExists(name)
	if err != nil {
		return false, fmt.Errorf("%w: NameExists: %s, error: %v", ErrIO, name, err)
	}
	return exist, err
}

// remove is wrapper for interface call Remove, returns appropriate wrapped error
func (m *Manager) remove(name string) error {
	err := m.i.Remove(name)
	if err != nil {
		return fmt.Errorf("%w: Remove: %s, error: %v", ErrIO, name, err)
	}
	return err
}

// AddToken adds token to token blacklists
func (m *Manager) addToken(token string) error {
	err := m.i.AddToken(token)
	if err != nil {
		return fmt.Errorf("%w: AddToken: %s, error: %v", ErrIO, token, err)
	}
	return nil
}
func (m *Manager) tokenExists(token string) (bool, error) {
	exists, err := m.i.TokenExists(token)
	if err != nil {
		return false, fmt.Errorf("%w: TokenExists: %s, error: %v", ErrIO, token, err)
	}
	return exists, nil
}
func (m *Manager) removeToken(token string) error {
	err := m.i.RemoveToken(token)
	if err != nil {
		return fmt.Errorf("%w: RemoveToken: %s, error: %v", ErrIO, token, err)
	}
	return nil
}
