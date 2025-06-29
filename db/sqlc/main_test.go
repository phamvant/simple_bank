package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var conn *sql.DB

func TestMain(m *testing.M) {
	var err error
	conn, err = sql.Open("postgres", "postgresql://admin:password@localhost:5432/simple_bank?sslmode=disable")

	print(conn)

	if err != nil {
		log.Fatal("can't connect to db")
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
