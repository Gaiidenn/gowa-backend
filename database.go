package main

import (
	"log"
	ara "github.com/solher/arangolite"
)

var db *ara.DB

func initDB() {
	db = ara.New().
	    LoggerOptions(false, false, false).
	    Connect("http://localhost:8529", "_system", "root", "")

	_, err := db.Run(&ara.CreateDatabase{
	    Name: "test",
	    Users: []map[string]interface{}{
	        {"username": "root", "passwd": ""},
	    },
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Database successfully created")
	}

	db.SwitchDatabase("test")
	initCollections()
}

func initCollections() {
	cols := []string{
		"users",
		"docs",
	}

	for _, col := range cols {
		createCollection(col)
	}
}

func createCollection(colName string) {
	 _, err := db.Run(&ara.CreateCollection{Name: colName})
	 if err != nil {
		 log.Println(err)
		 return
	 }
	 log.Println("Collection", colName, "successfully created")
}
