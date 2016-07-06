package users

import (
	"log"
	"fmt"
	"github.com/Gaiidenn/gowa-backend/rpcWebsocket"
)

type hub struct {
	// Register/Unregister a user as connected.
	register chan *User
	unregister chan *User
	connectedUsers map[string]*User
}

var h = hub{
	register:    make(chan *User),
	unregister:  make(chan *User),
	connectedUsers: make(map[string]*User),
}

func RunHub() {
	h.run();
}

func (h *hub) run() {
	for {
		select {
		case user := <-h.register:
			h.registerUser(user)
		case user := <-h.unregister:
			h.unregisterUser(user)
		}
	}
}

func (h *hub) registerUser(user *User) {
	key := user.Token
	user.Connected = true
	if h.keyExists(key) && h.connectedUsers[key].Username != user.Username {
		h.unregisterUser(user)
	}
	if (len(key) > 0) {
		h.connectedUsers[key] = user
		log.Println("User registered! key: ", key)
		fmt.Println("Number of users connected : ", len(h.connectedUsers))
		var reply *bool
		call := rpcWebsocket.RpcCall{
			Method: "UsersService.updateList",
			Args: user,
			Reply: reply,
		}
		rpcWebsocket.Broadcast(&call)
	}
}

func (h *hub) unregisterUser(user *User) {
	user.Connected = false
	key := user.Token
	if (len(key) > 0) {
		delete(h.connectedUsers, key)
		log.Println("User unregistered! key: ", key)
		user.Connected = false

		var reply *bool
		var call rpcWebsocket.RpcCall
		if user.Document.Key != nil && len(*user.Document.Key) > 0 {
			user.Token = ""
			call = rpcWebsocket.RpcCall{
				Method: "UsersService.updateList",
				Args: user,
				Reply: reply,
			}
		} else {
			call = rpcWebsocket.RpcCall{
				Method: "UsersService.removeFromList",
				Args: key,
				Reply: reply,
			}
		}
		rpcWebsocket.Broadcast(&call)
	}
}

func (h *hub) keyExists(key string) bool {
	_, ok := h.connectedUsers[key];
	return ok
}