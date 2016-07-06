package rpcWebsocket

import (
	"log"
	"fmt"
	"math/rand"
	"github.com/Gaiidenn/gowa-backend/users"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Inbound messages from the connections.
	broadcast          chan *RpcCall

	// Register requests from the connections.
	register           chan *connection

	// Unregister requests from connections.
	unregister         chan *connection

	// Register a user to its connection.
	registerUserChan   chan *users.User

	// Unregister a user from its connection.
	unregisterUserChan chan *users.User

	// Registered connections.
	connections        map[string]*connection
}

const KEY_LENGTH = 8

var h = hub{
	broadcast:   make(chan *RpcCall),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	registerUserChan:  make(chan *users.User),
	unregisterUserChan:  make(chan *users.User),
	connections: make(map[string]*connection),
}

func RunHub() {
	h.run();
}

func Broadcast(call *RpcCall) {
	h.broadcast <- call
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			log.Println("trying to register")
			key := h.generateKey(c)
			h.connections[key] = c
			log.Println("Connection registered! key: ", key)
			var reply bool
			call := &RpcCall{
				Method: "App.setToken",
				Args: key,
				Reply: &reply,
			}
			fmt.Println("Number of connections : ")
			fmt.Println(len(h.connections))
			c.call <- call
		case c := <-h.unregister:
			for key, cTmp := range h.connections {
				if cTmp == c {
					fmt.Println("Unregistring connection : ", key)
					delete(h.connections, key)
					close(c.call)
				}
			}
		case m := <-h.broadcast:
			log.Println("trying to broadcast")
			log.Println(m)
			for key, c := range h.connections {
				select {
				case c.call <- m:
				default:
					close(c.call)
					delete(h.connections, key)
				}
			}
		case user := <- h.registerUserChan:
			key := user.Token
			if len(key) == 0 {
				log.Println("empty key")
				return
			}
			cUser := h.connections[key].user
			user.Connected = true
			if h.keyExists(key) && cUser != nil && cUser.Username != user.Username {
				log.Println("unregistring before registring again")
				h.unregisterUserChan <- user
			}
			h.connections[key].user = user
			log.Println("User registered! key: ", key)
		case user := <- h.unregisterUserChan:
			user.Connected = false
			key := user.Token
			var u users.User
			*h.connections[key].user = u
			log.Println("User unregistered! key: ", key)
		}
	}
}

const lettersBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (h *hub) generateKey(c *connection) string {
	key := randStringBytes(KEY_LENGTH)
	log.Println("First key : ", key)
	v, ok := h.connections[key]
	for ok == true || (v != nil && v != c) {
		key = randStringBytes(KEY_LENGTH)
		v, ok = h.connections[key];
		log.Println("Other key : ", key)
	}
	log.Println("v, ok : ", v, ok)
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