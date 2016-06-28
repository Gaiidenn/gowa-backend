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
			key := h.generateKey(c)
			h.connections[key] = c
			fmt.Println("Number of connections : ")
			fmt.Println(len(h.connections))
		case c := <-h.unregister:
			key := string(c)
			if _, ok := h.connections[key]; ok {
				delete(h.connections, key)
				close(c.call)
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
	key := randStringBytes(8)
	for h.connections[key] != nil || h.connections[key] != c {
		key = randStringBytes(8)
	}
	return key
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = lettersBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}