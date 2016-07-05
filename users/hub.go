package users

import (
	"log"
	"fmt"
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
			key := user.Token
			h.connectedUsers[key] = user
			log.Println("User registered! key: ", key)
			fmt.Println("Number of users connected : ", len(h.connectedUsers))
		case user := <-h.unregister:
			key := user.Token
			delete(h.connectedUsers, key)
		}
	}
}
