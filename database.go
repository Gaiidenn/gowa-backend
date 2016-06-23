package main

import (
	"log"
	ara "github.com/solher/arangolite"
	"flag"
	"fmt"
)

var db *ara.DB

func init() {
	dbName = flag.String("dbName", config.DBName, "database name")
	dbUsername = flag.String("dbUsername", config.DBUsername, "database username")
	dbPassword = flag.String("dbPassword", config.DBPassword, "database password")
}

func initDB() {
	db = ara.New().
		LoggerOptions(false, false, false).
		Connect("http://localhost:8529", "_system", "root", "")
	fmt.Println(*dbUsername)
	_, err := db.Run(&ara.CreateDatabase{
		Name: *dbName,
		Users: []map[string]interface{}{
			{"username": *dbUsername, "passwd": *dbPassword},
		},
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Database successfully created")
	}

	db.SwitchDatabase(*dbName)
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
