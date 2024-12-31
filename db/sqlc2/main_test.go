package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const ()

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "root", "localhost", 3308, "simple_bank")
	var err error
	testDB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
