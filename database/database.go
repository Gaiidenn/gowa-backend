package database

import (
	"log"
	ara "github.com/solher/arangolite"
	"fmt"
)

var db *ara.DB

func InitConnection(dbName string, dbUsername string, dbPassword string) {
	db = ara.New().
		LoggerOptions(false, false, false).
		Connect("http://localhost:8529", "_system", "root", "")
	fmt.Println(dbUsername)
	_, err := db.Run(&ara.CreateDatabase{
		Name: dbName,
		Users: []map[string]interface{}{
			{"username": dbUsername, "passwd": dbPassword},
		},
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Database successfully created")
	}

	db.SwitchDatabase(dbName)
	initCollections()
}

func GetDB() *ara.DB {
	return db
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
