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
	ErrExist        = errors.New("user already exists")
	ErrNotExist     = errors.New("user doesn't exist")
	ErrInvalid      = errors.New("invalid argument")
	ErrIO           = errors.New("io error")
	ErrHash         = errors.New("hash error") // tried to generate this error during tests, didn't happen
	ErrExpired      = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
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
	}

	// Tokener realizes access to:
	// key via Key() used to generate token
	// duration via Duration() of generated key
	Tokener interface {
		Key() []byte
		Duration() time.Duration
	}

	StoreTokener interface {
		Store
		Tokener
	}
)

func New(storeTokener StoreTokener) *Manager {
	return &Manager{i: storeTokener}
}

func (u *Manager) Add(user User) error {
	if exists, err := u.Exists(user); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("%s %w", user.Name, ErrExist)
	}

	if len(user.Name) == 0 || len(user.Password) == 0 {
		return fmt.Errorf("%w: name and password must be provided", ErrInvalid)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%w %v", ErrHash, err)
	}

	if err := u.save(user.Name, hashedPassword); err != nil {
		return err
	}

	return nil
}

func (u *Manager) Remove(user User) error {
	if exists, err := u.Exists(user); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("%w %s", ErrNotExist, user.Name)
	}

	return u.i.Remove(user.Name)
}

func (u *Manager) Exists(user User) (bool, error) {
	return u.nameExists(user.Name)
}

func (u *Manager) Auth(user User) (bool, error) {
	if exists, err := u.Exists(user); err != nil {
		return false, err
	} else if !exists {
		return false, fmt.Errorf("%w %s", ErrNotExist, user.Name)
	}

	if hashPass, err := u.load(user.Name); err != nil {
		return false, err
	} else {
		return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(user.Password)) == nil, nil
	}
}

func (u *Manager) Token(user User) (string, error) {

	expires := time.Now().Add(u.i.Duration())
	user.claims.Name = user.Name
	user.claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: expires},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user.claims)
	return token.SignedString(u.i.Key())
}

func (u *Manager) ValidateToken(token string) (User, error) {
	var user User
	tkn, err := jwt.ParseWithClaims(token, &user.claims, func(token *jwt.Token) (interface{}, error) {
		return u.i.Key(), nil
	})

	if err != nil {
		if validationError, ok := err.(*jwt.ValidationError); ok {
			if (validationError.Errors & jwt.ValidationErrorExpired) == jwt.ValidationErrorExpired {
				return user, ErrExpired
			}
		}
		return user, err
	}

	if tkn.Valid {
		return user, nil
	}

	return user, ErrInvalidToken
}

// load is wrapper for interface call Load, returns appropriate wrapped error
func (u *Manager) load(name string) (data []byte, err error) {
	data, err = u.i.Load(name)
	if err != nil {
		return nil, fmt.Errorf("%w: Load: name %s, error: %v", ErrIO, name, err)
	}
	return data, err
}

// save is wrapper for interface call Save, returns appropriate wrapped error
func (u *Manager) save(name string, data []byte) error {
	if err := u.i.Save(name, data); err != nil {
		return fmt.Errorf("%w: Save: name %s, error: %v", ErrIO, name, err)
	}
	return nil
}

// nameExists is wrapper for interface call NameExists, returns appropriate wrapped error
func (u *Manager) nameExists(name string) (bool, error) {
	exist, err := u.i.NameExists(name)
	if err != nil {
		return false, fmt.Errorf("%w: NameExists: %s, error: %v", ErrIO, name, err)
	}
	return exist, err
}

// remove is wrapper for interface call Remove, returns appropriate wrapped error
func (u *Manager) remove(name string) error {
	err := u.i.Remove(name)
	if err != nil {
		return fmt.Errorf("%w: Remove: %s, error: %v", ErrIO, name, err)
	}
	return err
}
