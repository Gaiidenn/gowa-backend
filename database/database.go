package database

import (

	"database/sql"

	_ "gopkg.in/cq.v1"
	"log"
)

var db *sql.DB

func InitConnection(dbName string, dbUsername string, dbPassword string) {
	var baseURL string
	baseURL = "http://" + dbUsername + ":" + dbPassword + "@localhost:7474"
	dbTmp, err := sql.Open("neo4j-cypher", baseURL)
	if err != nil {
		log.Fatal(err)
	}
	db = dbTmp
	log.Println("Connection to database successfull")
	//initCollections()
}

func GetDB() *sql.DB {
	return db
}
