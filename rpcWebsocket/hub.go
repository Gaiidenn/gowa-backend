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
	broadcast:   	make(chan *RpcCall),
	register:    	make(chan *connection),
	unregister:  	make(chan *connection),
	registerUser:  	make(chan *users.User),
	unregisterUser: make(chan *users.User),
	connections: 	make(map[string]*connection),
}

func RunHub() {
	h.run();
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
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
		case c := <-h.unregister:
			for key, cTmp := range h.connections {
				if cTmp == c {
					log.Println("Unregistring connection : ", key)
					delete(h.connections, key)
					close(c.call)
				}
			}
		case m := <- h.broadcast:
			log.Println("trying to broadcast")
			log.Println(m)
			for key, c := range h.connections {
				select {
				case c.call <- m:
				default:
					log.Println("Broadcast failed => deleting connection : ", key)
					delete(h.connections, key)
					close(c.call)
				}
			}
		case user := <- h.registerUser:
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
			h.connections[key].user = user
			log.Println("User registered! key: ", key, "; user:", *h.connections[key].user)
		case user := <- h.unregisterUser:
			key := user.Token
			if (len(key) > 0) {
				if _, ok := h.connections[key]; ok {
					*h.connections[key].user = users.User{}
				}
				log.Println("User unregistered! key: ", key)
			}
		}
	}
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