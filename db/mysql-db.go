package db

import (
	"database/sql"
	db "example/db/sqlc"
	db2 "example/db/sqlc2"

	"fmt"
	"log"
)

var Query *db.Queries
var Query2 *db2.Queries
var DB *sql.DB
var DB2 *sql.DB

func init() {
	var err error
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "root", "localhost", 3308, "service_pro")
	DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	Query = db.New(DB)
	Query2 = db2.New(DB2)

}
