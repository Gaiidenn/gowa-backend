package users

import (
	"log"
	"errors"
)

// Log the user in app
func (user *User) Login() error {
	token := user.Token
	userTmp, err := user.getByUsername(user.Username)
	if err != nil {
		return err
	}
	if userTmp == nil {
		return errors.New("unknown username")
	}
	if (userTmp.Password != user.Password) {
		return errors.New("wrong password")
	}
	*user = *userTmp
	if len(token) > 0 {
		user.Token = token
	}
	user.Connected = true
	return nil
}

// Get all Users from collection
func (user *User) GetAllWithConnected() (*[]User, error) {
	users, err := user.GetAll()
	if err != nil {
		return nil, err
	}
	// Adding connected users
	for _, cUser := range h.connectedUsers {
		exists := false
		for i, user := range users {
			if cUser.Username == user.Username {
				exists = true
				users[i].Connected = true
				users[i].Token = cUser.Token
				continue
			}
		}
		if !exists {
			users = append(users, *cUser)
		}
	}

	return &users, nil
}

func (user *User) availableUsername() (bool, error) {
	u, err := user.getByUsername(user.Username)
	if err != nil {
		return false, err
	}
	var key string
	if user.Document.Key != nil {
		key = *user.Document.Key
	} else {
		key = ""
	}
	if u != nil && u.Username == user.Username && *u.Document.Key != key {
		return false, nil
	}
	for _, u = range h.connectedUsers {
		if u.Username == user.Username && u.Token != user.Token {
			return false, nil
		}
	}
	return true, nil
}

func (user *User) readyForSave() bool {
	log.Println(len(user.Username), len(user.Password), len(user.Email))
	if len(user.Username) < 4 {
		return false
	}
	if len(user.Password) < 4 {
		return false
	}
	if len(user.Email) < 4 {
		return false
	}
	return true
}