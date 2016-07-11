package rpcWebsocket

import (
	"errors"
	"github.com/Gaiidenn/gowa-backend/users"
	"log"
)

// UserRPCService for jsonRPC requests
type UserRPCService struct {
}

// Save the user in database
func (us *UserRPCService) Save(user *users.User, reply *users.User) error {
	free, err := user.AvailableUsername()
	if err != nil {
		return err
	}
	if !free {
		return errors.New("username already exists")
	}

	for _, c := range h.connections {
		if c.user != nil && c.user.Username == user.Username && c.user.Token != user.Token {
			return nil
		}
	}

	if user.ReadyForSave() {
		// Saving token
		key := user.Token
		err := user.Save()
		if err != nil {
			return err
		}
		user.Token = key
	}

	h.registerUser <- user

	*reply = *user
	return nil
}

// Log the user in app
func (us *UserRPCService) Login(userLogin *users.User, user *users.User) error {
	log.Println("Login : ", userLogin)
	key := userLogin.Token

	err := userLogin.Login()
	if err != nil {
		return err
	}
	// userLogin now is full filled
	userLogin.Connected = true
	userLogin.Token = key

	h.registerUser <- userLogin

	*user = *userLogin
	return nil
}

// Log the user in app
func (us *UserRPCService) Logout(user *users.User, ok *bool) error {
	user.Connected = false
	h.unregisterUser <- user

	*ok = true
	return nil
}



// Get all UserRPCService from collection
func (us *UserRPCService) GetAll(_ *string, reply *[]users.User) error {
	var user users.User
	users, err := user.GetAll()
	if err != nil {
		return err
	}
	// Adding connected users
	for _, c := range h.connections {
		if c.user != nil {
			exists := false
			for i, u := range users {
				if c.user.Username == u.Username {
					exists = true
					users[i].Connected = true
					users[i].Token = c.user.Token
					continue
				}
			}
			if !exists && len(c.user.Username) > 0 {
				users = append(users, *c.user)
			}
		}
	}

	if err != nil {
		return err
	}
	*reply = users
	return nil
}