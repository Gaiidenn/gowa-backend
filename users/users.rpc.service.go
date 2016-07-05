package users

import (
	"github.com/Gaiidenn/gowa-backend/rpcWebsocket"
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
		user.Username = ""
		h.register <- user
		return errors.New("username already exists")
	}
	h.register <- user
	if user.readyForSave() {
		err := user.Save()
		if err != nil {
			return err
		}
	}
	var s string
	call := rpcWebsocket.RpcCall{
		Method: "UsersService.updateList",
		Args: user,
		Reply: &s,
	}
	rpcWebsocket.Broadcast(&call)
	*reply = *user
	return nil
}

// Log the user in app
func (us *UserRPCService) Login(userLogin *User, user *User) error {
	err := userLogin.Login()
	if err != nil {
		return err
	}
	*user = *userLogin
	return nil
}

// Get all UserRPCService from collection
func (us *UserRPCService) GetAll(_ *string, reply *[]User) error {
	var user User;
	users, err := user.GetAll()
	if err != nil {
		return err
	}
	*reply = *users
	return nil
}
