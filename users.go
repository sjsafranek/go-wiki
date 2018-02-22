package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"compress/gzip"
	"os"

	"./utils"
)

// User: user for wiki engine
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SetPassword sets password
func (self *User) SetPassword(password string) {
	self.Password = utils.SHA512FromBytes([]byte(password))
}

// IsPassword checks if password is the set password
func (self *User) IsPassword(password string) bool {
	return self.Password == utils.SHA512FromBytes([]byte(password))
}

// Users: collection of users
type Users struct {
	// Filename string
	Users []*User `json:"users"`
}

// Fetch: fetches json file containing users array.
// @args file{string}	users file
func (self *Users) Fetch(file string) error {
	b, err := ioutil.ReadFile(file)
	if nil != err {
		return err
	}
	return self.Unmarshal(string(b))
}

// Save: saves users to json file
func (self *Users) Save(file string) error {
	contents, err := self.Marshal()
	if nil != err {
		return err
	}
	return ioutil.WriteFile(file, []byte(contents), 0644)
}

// Unmarshal: json unmarshals string to struct
// @args string
// @return error
func (self *Users) Unmarshal(data string) error {
	return json.Unmarshal([]byte(data), self)
}

// Marshal: json marshals struct
// @return string
// @return error
func (self Users) Marshal() (string, error) {
	b, err := json.Marshal(self)
	if nil != err {
		return "", err
	}
	return string(b), nil
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

func (self *Users) SaveGzip() error {
	f, err := os.Create("users.json.gz")
	if err != nil {
		return err
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	defer w.Close()
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	return enc.Encode(self.Users)
}

func (self *Users) LoadGzip() error {
	var err error
	f, err := os.Open("users.json.gz")
	defer f.Close()
	if err != nil {
		return err
	}

	reader, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer reader.Close()

	err = json.NewDecoder(reader).Decode(self.Users)
	if err != nil {
		return err
	}

	return nil
}
