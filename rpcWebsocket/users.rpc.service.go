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
		err := user.Save()
		if err != nil {
			return err
		}
	}

	h.registerUserChan <- user

	var resp *bool
	call := RpcCall{
		Method: "UsersService.updateList",
		Args: user,
		Reply: resp,
	}
	log.Println("Trying to broadcast")
	h.broadcast <- &call

	*reply = *user
	return nil
}
/*
var reply *bool
var call RpcCall
if user.Document.Key != nil && len(*user.Document.Key) > 0 {
user.Token = ""
call = RpcCall{
Method: "UsersService.updateList",
Args: user,
Reply: reply,
}
} else {
call = RpcCall{
Method: "UsersService.removeFromList",
Args: key,
Reply: reply,
}
}
h.broadcast <- &call
*/
// Log the user in app
func (us *UserRPCService) Login(userLogin *users.User, user *users.User) error {
	err := userLogin.Login()
	if err != nil {
		return err
	}
	h.registerUserChan <- userLogin

	var reply *bool
	var call RpcCall
	user.Token = ""
	call = RpcCall{
		Method: "UsersService.updateList",
		Args: user,
		Reply: reply,
	}
	h.broadcast <- &call

	*user = *userLogin
	return nil
}

// Log the user in app
func (us *UserRPCService) Logout(user *users.User, ok *bool) error {
	h.unregisterUserChan <- user

	key := user.Token
	var reply *bool
	var call RpcCall
	call = RpcCall{
		Method: "UsersService.removeFromList",
		Args: key,
		Reply: reply,
	}
	h.broadcast <- &call

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
			if !exists {
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
