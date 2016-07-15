package rpcWebsocket

import (
	"log"
	"math/rand"
	"github.com/Gaiidenn/gowa-backend/users"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Inbound messages from the connections.
	broadcast      chan *RpcCall

	// Register requests from the connections.
	register       chan *connection

	// Unregister requests from connections.
	unregister     chan *connection

	// Register a user to its connection.
	registerUser   chan *users.User

	// Unregister a user from its connection.
	unregisterUser chan *users.User

	// Registered connections.
	connections    map[string]*connection
}

const KEY_LENGTH = 8

var h = hub{
	broadcast:   	make(chan *RpcCall, 100),
	register:    	make(chan *connection, 100),
	unregister:  	make(chan *connection, 100),
	registerUser:  	make(chan *users.User, 100),
	unregisterUser: make(chan *users.User, 100),
	connections: 	make(map[string]*connection),
}

func RunHub() {
	h.run();
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.Register(c)
		case c := <-h.unregister:
			h.Unregister(c)
		case m := <- h.broadcast:
			h.Broadcast(m)
		case user := <- h.registerUser:
			h.RegisterUser(user)
		case user := <- h.unregisterUser:
			h.UnregisterUser(user)
		}
	}
}

func (h *hub) Register(c *connection) {
	log.Println("trying to register")
	key := c.user.Token
	h.connections[key] = c
	log.Println("Connection registered! key: ", key)
	log.Println("Number of connections : ", len(h.connections))

	var reply *bool
	call := RpcCall{
		Method: "App.setToken",
		Args: key,
		Reply: reply,
	}
	c.call <- &call

	h.broadcastConnectedCount()
}

func (h *hub) Unregister(c *connection) {
	for key, cTmp := range h.connections {
		if cTmp == c {
			h.UnregisterUser(c.user)
			log.Println("Unregistring connection : ", key)
			delete(h.connections, key)
			close(c.call)
		}
	}
}

func (h *hub) Broadcast(m *RpcCall) {
	log.Println("trying to broadcast")
	log.Println(m)
	for key, c := range h.connections {
		select {
		case c.call <- m:
		default:
			log.Println("Broadcast failed => deleting connection : ", key)
			h.Unregister(c)
		}
	}
}

func (h *hub) RegisterUser(user *users.User) {
	log.Println("trying to register")
	key := user.Token
	if len(key) == 0 {
		log.Println("empty key")
		return
	}
	if _, ok := h.connections[key]; !ok {
		log.Println("unknown connection (key : \"", key, "\")")
		return
	}
	if h.connections[key].user != nil && h.connections[key].user.Token == user.Token && h.connections[key].user.Username != user.Username {
		 h.UnregisterUser(h.connections[key].user)
	}
	h.connections[key].user = user
	log.Println("User registered! key: ", key, "; user:", *h.connections[key].user)
	var reply *bool
	call := RpcCall{
		Method: "UsersService.updateList",
		Args: user,
		Reply: reply,
	}
	h.Broadcast(&call)

	h.broadcastConnectedCount()
}

func (h *hub) UnregisterUser(user *users.User) {
	key := user.Token
	if (len(key) > 0) {
		if _, ok := h.connections[key]; ok {
			if h.connections[key].user.ID != "" {
				log.Println("offline : ", *h.connections[key].user)
				h.connections[key].user.Connected = false;
				var reply *bool
				call := RpcCall{
					Method: "UsersService.updateList",
					Args: user,
					Reply: reply,
				}
				h.Broadcast(&call)
			} else {
				log.Println("User unregistered! key: ", key)
				var reply *bool
				call := RpcCall{
					Method: "UsersService.removeFromList",
					Args: user.Username,
					Reply: reply,
				}
				h.Broadcast(&call)
			}
			*h.connections[key].user = users.User{}
		}
	}
	h.broadcastConnectedCount()
}

type ConnectedUsersCount struct {
	anonymous     int `json:"anonymous"`
	notRegistered int `json:"notRegistered"`
	registered    int `json:"registered"`
}

func (h *hub) broadcastConnectedCount() {
	count := h.connectedUsersCount()
	var reply *bool
	call := RpcCall{
		Method: "UsersService.setConnectedCount",
		Args: count,
		Reply: reply,
	}
	h.Broadcast(&call)
}

func (h *hub) connectedUsersCount() *ConnectedUsersCount {
	var count ConnectedUsersCount
	count.anonymous = 0;
	count.notRegistered = 0;
	count.registered = 0;

	for _, c := range h.connections {
		u := c.user
		if u == nil {
			count.anonymous++
			continue
		}
		if len(u.Username) == 0 {
			count.anonymous++
			continue
		}
		if len(u.ID) == 0 {
			count.notRegistered++
			continue
		}
		count.registered++
	}

	return &count
}

const lettersBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (h *hub) generateKey() string {
	key := randStringBytes(KEY_LENGTH)
	log.Println("First key : ", key)
	_, ok := h.connections[key]
	for ok == true {
		key = randStringBytes(KEY_LENGTH)
		_, ok = h.connections[key];
		log.Println("Other key : ", key)
	}
	return key
}

func cleanKey(key []byte) string {
	var skey = ""
	for _, v := range key {
		valid := false
		for _, lchar := range lettersBytes {
			if string(v) == string(lchar) {
				valid = true
			}
		}
		if valid {
			skey = skey + string(v)
		}
	}
	return skey
}

func (h *hub) keyExists(key string) bool {
	_, ok := h.connections[key];
	return ok
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = lettersBytes[rand.Intn(len(lettersBytes))]
	}
	return string(b)
}