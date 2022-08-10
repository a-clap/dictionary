//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package users

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	i LoadSaver
}

var (
	ErrExist    = errors.New("user already exists")
	ErrNotExist = errors.New("user doesn't not exist")
	ErrInvalid  = errors.New("invalid argument")
	ErrIO       = errors.New("io error")
	ErrHash     = errors.New("hash error") // tried to generate this error during tests, didn't happen
)

type (
	// LoadSaver realizes access to some kind of user database (maybe even just map[string]string)
	// Errors returned by interface should be ONLY related to internal IO errors
	LoadSaver interface {
		// Load user password, if user doesn't exist, return ""
		Load(name string) (password string, err error)
		// Save new user with password, overwrites user password, if already exists
		Save(name, password string) error
		// NameExists allows to check whether user with name already exists
		NameExists(name string) (bool, error)
		// Remove user with provided name, if user doesn't exist, don't do anything
		Remove(name string) error
	}
)

func New(loadSaver LoadSaver) *Users {
	return &Users{i: loadSaver}
}

func (u *Users) Add(name, password string) error {
	if exists, err := u.Exists(name); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("%s %w", name, ErrExist)
	}

	if len(name) == 0 || len(password) == 0 {
		return fmt.Errorf("%w: name and password must be provided", ErrInvalid)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%w %v", ErrHash, err)
	}

	if err := u.save(name, string(hashedPassword)); err != nil {
		return err
	}

	return nil
}

func (u *Users) Remove(name string) error {
	if exists, err := u.Exists(name); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("%w %s", ErrNotExist, name)
	}

	return u.Remove(name)
}

func (u *Users) Exists(name string) (bool, error) {
	return u.nameExists(name)
}

func (u *Users) Auth(name, password string) (bool, error) {
	if exists, err := u.Exists(name); err != nil {
		return false, err
	} else if !exists {
		return false, fmt.Errorf("%w %s", ErrNotExist, name)
	}

	if hashPass, err := u.load(name); err != nil {
		return false, err
	} else {
		return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password)) == nil, nil
	}
}

// load is wrapper for interface call Load, returns appropriate wrapped error
func (u *Users) load(name string) (password string, err error) {
	password, err = u.i.Load(name)
	if err != nil {
		return "", fmt.Errorf("%w: Load: name %s, error: %v", ErrIO, name, err)
	}
	return password, err
}

// save is wrapper for interface call Save, returns appropriate wrapped error
func (u *Users) save(name, password string) error {
	if err := u.i.Save(name, password); err != nil {
		return fmt.Errorf("%w: Save: name %s, error: %v", ErrIO, name, err)
	}
	return nil
}

// nameExists is wrapper for interface call NameExists, returns appropriate wrapped error
func (u *Users) nameExists(name string) (bool, error) {
	exist, err := u.i.NameExists(name)
	if err != nil {
		return false, fmt.Errorf("%w: NameExists: %s, error: %v", ErrIO, name, err)
	}
	return exist, err
}

// remove is wrapper for interface call Remove, returns appropriate wrapped error
func (u *Users) remove(name string) error {
	err := u.i.Remove(name)
	if err != nil {
		return fmt.Errorf("%w: Remove: %s, error: %v", ErrIO, name, err)
	}
	return err
}
