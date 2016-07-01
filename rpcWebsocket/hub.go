package rpcWebsocket

import (
	"log"
	"fmt"
	"math/rand"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Inbound messages from the connections.
	broadcast chan *RpcCall

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection

	// Registered connections.
	connections map[string]*connection
}

const KEY_LENGTH = 8

var h = hub{
	broadcast:   make(chan *RpcCall),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
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