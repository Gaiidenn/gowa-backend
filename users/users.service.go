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
/*
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
*/
func (user *User) AvailableUsername() (bool, error) {
	u, err := user.getByUsername(user.Username)
	if err != nil {
		return false, err
	}

	if u != nil && u.Username == user.Username && u.ID != user.ID {
		return false, nil
	}

	return true, nil
}

func (user *User) ReadyForSave() bool {
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