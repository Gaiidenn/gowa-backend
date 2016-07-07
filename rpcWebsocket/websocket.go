package rpcWebsocket

import (
	"golang.org/x/net/websocket"
	"net/rpc"
	"net/rpc/jsonrpc"
	"log"
	"time"
	"github.com/Gaiidenn/gowa-backend/users"
)

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection
	ws *websocket.Conn

	// The rpc client
	rc *rpc.Client

	// The user attached to connection
	user *users.User

	// Buffered channel of outbound messages.
	call chan *RpcCall
}

type RpcCall struct {
	Method string
	Args interface{}
	Reply interface{}
}

func JsonrpcHandler(ws *websocket.Conn) {
	log.Println("connection websocket on jsonrpcHandler")
	key := readKey(ws)
	if !h.keyExists(key) {
		log.Println("invalid key")
		return
	}
	jsonrpc.ServeConn(ws)
}

func PushHandler(ws *websocket.Conn) {
	log.Println("connection websocket on pushHandler")

	var c *connection

	// Looking for an existing key
	key := readKey(ws)
	if len(key) > 0 && h.keyExists(key) {
		// Key exists => get connection
		c = h.connections[key]
	} else {
		// No key or key doesn't exists => creating a new connection
		rc := jsonrpc.NewClient(ws)
		key = h.generateKey()
		user := &users.User{
			Token: key,
		}
		call := make(chan *RpcCall)
		c = &connection{
			ws: ws,
			rc: rc,
			user: user,
			call: call,
		}
	}

	// Registering the connection
	h.register <- c

	c.callPump()
}

func readKey(ws *websocket.Conn) string {
	var key = make([]byte, KEY_LENGTH + 4) // KEY_LENGTH + 4 because what we receive is ["..key."]
	_, err := ws.Read(key)
	if err != nil {
		log.Println(err)
		return ""
	}
	return cleanKey(key)
}

// callPump pumps calls from the hub to the rpc connection.
func (c *connection) callPump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		c.rc.Close()
		h.unregister <- c
	}()
	for {
		select {
		case call, ok := <- c.call :
			log.Println("trying to call :", call)
			if !ok {
				h.unregister <- c
				return
			}
			if err := c.rc.Call(call.Method, call.Args, &call.Reply); err != nil {
				log.Println("error calling : ")
				log.Println(err)
				h.unregister <- c
				return
			}
			log.Println("call ok -> reply : ", call.Reply)
		case <- ticker.C :
			if _, err := c.ws.Write([]byte{}); err != nil {
				return
			}
		}
	}
}