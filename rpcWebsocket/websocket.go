package rpcWebsocket

import (
	//"log"
	//"net/http"
	"golang.org/x/net/websocket"
	"net/rpc"
	"net/rpc/jsonrpc"
	"log"
	//"time"
	"time"
	"fmt"
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
	/*
	var reply string

	err := rc.Call("App.log", "My test", &reply)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(reply)*/
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
				fmt.Println("no response from client -> closing connection!")
				return
			}
			fmt.Println("client still connected -> ok!")
		}
	}
}
/*
// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	//c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	//c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Println("Message received from client : ", string(message))
		h.broadcast <- message
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	//ticker := time.NewTicker(pingPeriod)
	defer func() {
		//ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				//c.write(websocket.CloseMessage, []byte{})
				return
			}
			//if err := c.write(websocket.TextMessage, message); err != nil {
			if err := websocket.Message.Send(c.ws, message); err != nil {
				return
			}
		/*case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
//func serveWs(w http.ResponseWriter, r *http.Request) {
func serveWs(ws *websocket.Conn) {

	rc := jsonrpc.NewClient(ws)
	c := &connection{
		send: make(chan []byte, 256),
		ws: ws,
		rc: rc,
	}
	h.register <- c

	go c.writePump()
	jsonrpc.ServeConn(c.ws)
	//c.readPump()
}
*/
