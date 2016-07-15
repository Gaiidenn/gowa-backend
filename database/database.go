package database

import (

	"database/sql"

	_ "gopkg.in/cq.v1"
	"log"
)

var db *sql.DB

func InitConnection(dbName string, dbUsername string, dbPassword string) {
	/*db = ara.New().
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
		//log.Println(err)
	} else {
		//log.Println("Database successfully created")
	}

	db.SwitchDatabase(dbName)*/
	var baseURL string
	baseURL = "http://" + dbUsername + ":" + dbPassword + "@localhost:7474"
	dbTmp, err := sql.Open("neo4j-cypher", baseURL)
	if err != nil {
		log.Fatal(err)
	}
	db = dbTmp
	//log.Println("Connection to database successfull")
	//initCollections()
}

func GetDB() *sql.DB {
	return db
}
/*
func initCollections() {
	cols := []string{
		"users",
		"chats",
	}

	for _, col := range cols {
		createCollection(col)
	}
}

func createCollection(colName string) {
	 _, err := db.Run(&ara.CreateCollection{Name: colName})
	 if err != nil {
		 //log.Println(err)
		 return
	 }
	 //log.Println("Collection", colName, "successfully created")
}
*/