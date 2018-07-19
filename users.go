package main

import (
	"errors"

	"github.com/sjsafranek/goutils/hashers"
	"github.com/sjsafranek/goutils/jsonhelpers"
)

// User: user for wiki engine
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SetPassword sets password
func (self *User) SetPassword(password string) {
	// self.Password = hashers.SHA512FromBytes([]byte(password))
	self.Password = hashers.Sha512HashString(password)
}

// IsPassword checks if password is the set password
func (self *User) IsPassword(password string) bool {
	// return self.Password == hashers.SHA512FromBytes([]byte(password))
	return self.Password == hashers.Sha512HashString(password)
}

// Users: collection of users
type Users struct {
	// Filename string
	Users []*User `json:"users"`
}

// Fetch: fetches json file containing users array.
// @args file{string}	users file
func (self *Users) Fetch(file string) error {
	return jsonhelpers.Fetch(file, self)
}

// Save: saves users to json file
func (self *Users) Save(file string) error {
	return jsonhelpers.Save(file, self)
}

// Unmarshal: json unmarshals string to struct
// @args string
// @return error
func (self *Users) Unmarshal(data string) error {
	return jsonhelpers.Unmarshal(data, self)
}

// Marshal: json marshals struct
// @return string
// @return error
func (self Users) Marshal() (string, error) {
	return jsonhelpers.Marshal(self)
}

// Get user by username
func (self *Users) Get(username string) (*User, error) {
	for _, user := range self.Users {
		if username == user.Username {
			return user, nil
		}
	}
	return &User{}, errors.New("User not found")
}

// Has has user with username
func (self *Users) Has(username string) bool {
	_, err := self.Get(username)
	return err == nil
}

// Add user to users
func (self *Users) Add(user *User) error {
	if !self.Has(user.Username) {
		self.Users = append(self.Users, user)
		return nil
	}

	return errors.New("User already exists")
}
