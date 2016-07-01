package rpcWebsocket

import (
	"golang.org/x/net/websocket"
	"net/rpc"
	"net/rpc/jsonrpc"
	"log"
	"time"
)

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection
	ws *websocket.Conn

	// The rpc client
	rc *rpc.Client

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
	var key = make([]byte, KEY_LENGTH + 4) // KEY_LENGTH + 4 because what we receive is ["..key."]
	_, err := ws.Read(key)
	if err != nil {
		log.Println(err)
		return
	}
	skey := cleanKey(key)
	if !h.keyExists(skey) {
		log.Println("invalid key")
		return
	}
	jsonrpc.ServeConn(ws)
}

func PushHandler(ws *websocket.Conn) {
	log.Println("connection websocket on pushHandler")
	rc := jsonrpc.NewClient(ws)

	call := make(chan *RpcCall)
	c := &connection{
		ws: ws,
		rc: rc,
		call: call,
	}
	h.register <- c
	c.callPump()
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
			log.Println("call ok -> reply : ")
			log.Println(call.Reply)
		case <- ticker.C :
			if _, err := c.ws.Write([]byte{}); err != nil {
				return
			}
		}
	}
}