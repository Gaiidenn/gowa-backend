package users

import (
	"github.com/Gaiidenn/gowa-backend/rpcWebsocket"
	"log"
	"errors"
)

// UserRPCService for jsonRPC requests
type UserRPCService struct {
}

// Save the user in database
func (us *UserRPCService) Save(
		params *struct{
			Token string
			User User
		}, reply *User) error {
	log.Println(params)
	log.Println(string(params.User.Username))
	user := params.User

	free, err := user.availableUsername()
	if err != nil {
		return err
	}
	if !free {
		return errors.New("username already exists")
	}
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
	*reply = user
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
