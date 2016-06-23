package main

import (
	"flag"
	"log"
	"net/http"
	"net/rpc"
	"io/ioutil"
	"strings"
	"text/template"
	"encoding/json"
	"golang.org/x/net/websocket"
)

type Configuration struct {
	ClientPath string `json: "clientPath"`
	DBName string `json: "dbName"`
	DBUsername string `json: "dbUsername"`
	DBPassword string `json: "dbPassword"`
}
var config = Configuration{}

var addr *string
var clientDir *string
var dbName *string
var dbUsername *string
var dbPassword *string
var homeTempl *template.Template

func init() {
	loadConfig()
	addr = flag.String("addr", ":8080", "http service address")
	clientDir = flag.String("clientDir", config.ClientPath, "client app directory")
	homeTempl = template.Must(template.ParseFiles(config.ClientPath + "/index.html"))
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

	initDB()

	// Initialize websocket hub
	go h.run()

	// Register rpc methods
	initRPCRegistration()

	// Define requests handlers
	http.Handle("/jsonrpc", websocket.Handler(jsonrpcHandler))
	http.Handle("/push", websocket.Handler(pushHandler))
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
	homeTempl.Execute(w, r.Host)
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
	userRPCService := new(UserRPCService)
	rpc.Register(userRPCService)
}

/*
// Msg type
type Msg string

// Echo just response with the same msg as received
func (msg *Msg) Echo(str string, reply *string) error {
	log.Println("Msg.Echo(", str, ")")
	*reply = str
	return nil
}
*/
