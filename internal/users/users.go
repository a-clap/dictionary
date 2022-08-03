//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package users

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	LoadSaver
}

func errorInterface(f string, err error) error {
	return fmt.Errorf("error on interface %s: %#v", f, err)
}

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
	return &Users{loadSaver}
}

func (u *Users) Add(name, password string) error {
	if exists, err := u.Exists(name); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("user %s already exist", name)
	}

	if len(name) == 0 || len(password) == 0 {
		return fmt.Errorf("name and password must be provided")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error on bcrypt %#v", err)
	}

	if err := u.Save(name, string(hashedPassword)); err != nil {
		return errorInterface("Save", err)
	}

	return nil
}

func (u *Users) Remove(name string) error {
	if exists, err := u.Exists(name); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("user %s doesn't exists", name)
	}

	return u.Remove(name)
}

func (u *Users) Exists(name string) (bool, error) {
	exists, err := u.NameExists(name)
	if err != nil {
		return false, errorInterface("NameExists", err)
	}
	return exists, nil
}

func (u *Users) Auth(name, password string) (bool, error) {
	if exists, err := u.Exists(name); err != nil {
		return false, err
	} else if !exists {
		return false, fmt.Errorf("user %s doesn't exists", name)
	}

	if hashPass, err := u.Load(name); err != nil {
		return false, errorInterface("Load", err)
	} else {
		return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password)) == nil, nil
	}

}
