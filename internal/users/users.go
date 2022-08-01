package users

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	LoadSaver
}

type user struct {
	name     string
	password string
}
type (
	// LoadSaver realizes access to some kind of user database (may be even just map[string]string)
	// Errors returned by interface should be ONLY related to internal IO errors
	LoadSaver interface {
		// Load user password, if user doesn't exist, return ""
		Load(name string) (password string, err error)
		// Save new user with password, overwrites user password, if already exists
		Save(name, password string) error
		// Exists allows to check whether user with name already exists
		Exists(name string) (bool, error)
		// Remove user with provided name, if user doesn't exist, don't do anything
		Remove(name string) error
	}
)

func New(loadSaver LoadSaver) *Users {
	return &Users{loadSaver}
}

func (u *Users) Add(name, password string) error {
	if len(name) == 0 || len(password) == 0 {
		return fmt.Errorf("name and password must be provided")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error on bcrypt %#v", err)
	}

	if err := u.Save(name, string(hashedPassword)); err != nil {
		return fmt.Errorf("interface error %#v", err)
	}

	return nil
}

func (u *Users) Remove(name string) error {
	return nil
}

func (u *Users) Login(name, password string) {

}
