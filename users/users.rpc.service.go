package users

import (
	"errors"
)

// UserRPCService for jsonRPC requests
type UserRPCService struct {
}

// Save the user in database
func (us *UserRPCService) Save(user *User, reply *User) error {
	free, err := user.availableUsername()
	if err != nil {
		return err
	}
	if !free {
		return errors.New("username already exists")
	}
	h.register <- user
	if user.readyForSave() {
		err := user.Save()
		if err != nil {
			return err
		}
	}
	*reply = *user
	return nil
}

// Log the user in app
func (us *UserRPCService) Login(userLogin *User, user *User) error {
	err := userLogin.Login()
	if err != nil {
		return err
	}
	h.register <- userLogin

	*user = *userLogin
	return nil
}

// Log the user in app
func (us *UserRPCService) Logout(user *User, ok *bool) error {
	h.unregister <- user
	*ok = true
	return nil
}



// Get all UserRPCService from collection
func (us *UserRPCService) GetAll(_ *string, reply *[]User) error {
	var user User;
	users, err := user.GetAllWithConnected()
	if err != nil {
		return err
	}
	*reply = *users
	return nil
}
