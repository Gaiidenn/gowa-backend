package main

import (
	"flag"
	"log"
	"net/http"
	"net/rpc"
	"io/ioutil"
	"strings"
	"encoding/json"
	"golang.org/x/net/websocket"
	"github.com/Gaiidenn/gowa-backend/database"
	"github.com/Gaiidenn/gowa-backend/rpcWebsocket"
)

type Configuration struct {
	ClientPath string `json: "clientPath"`
	DBName     string `json: "dbName"`
	DBUsername string `json: "dbUsername"`
	DBPassword string `json: "dbPassword"`
}
var config = Configuration{}

var addr *string
var clientDir *string
var dbName *string
var dbUsername *string
var dbPassword *string

func init() {
	loadConfig()
	addr = flag.String("addr", ":8080", "http service address")
	clientDir = flag.String("clientDir", config.ClientPath, "client app directory")
	dbName = flag.String("dbName", config.DBName, "database name")
	dbUsername = flag.String("dbUsername", config.DBUsername, "database username")
	dbPassword = flag.String("dbPassword", config.DBPassword, "database password")
}

func loadConfig() {
	if config == (Configuration{}) {
		file, err := ioutil.ReadFile("config.json")
		if err != nil {
			log.Fatal("OpenConfigFile: ", err)
		}
		err = json.Unmarshal(file, &config)
		if err != nil {
			log.Fatal("ParseConfigFile: ", err)
		}
	}
}

func main() {
	flag.Parse()

	database.InitConnection(*dbName, *dbUsername, *dbPassword)

	// Initialize websocket hub
	go rpcWebsocket.RunHub()

	// Register rpc methods
	initRPCRegistration()

	// Define requests handlers
	http.Handle("/jsonrpc", websocket.Handler(rpcWebsocket.JsonrpcHandler))
	http.Handle("/push", websocket.Handler(rpcWebsocket.PushHandler))
	http.HandleFunc("/", serveIndex)

	// Start server
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if validFileRequest(r.URL.Path) {
		filePath := *clientDir + r.URL.Path
		http.ServeFile(w, r, filePath)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, *clientDir + "/index.html")
}

func validFileRequest(path string) bool {
	if path == "/" {
		return false
	}

	s := strings.Split(path, "/")
	file := strings.Split(s[len(s)-1], ".")

	if len(file) < 2 {
		return false
	}

	if fileExt := file[len(file)-1]; fileExt == "js" || fileExt == "html" || fileExt == "css" {
		return true
	}

	return false
}

func initRPCRegistration() {
	userRPCService := new(rpcWebsocket.UserRPCService)
	rpc.Register(userRPCService)
}