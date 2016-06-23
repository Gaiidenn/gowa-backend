package rpcWebsocket

import "log"

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
	connections map[*connection]bool
}

var h = hub{
	broadcast:   make(chan *RpcCall),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
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
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.call)
			}
		case m := <-h.broadcast:
			log.Println("trying to broadcast")
			log.Println(m)
			for c := range h.connections {
				select {
				case c.call <- m:
				default:
					close(c.call)
					delete(h.connections, c)
				}
			}
		}
	}
}
