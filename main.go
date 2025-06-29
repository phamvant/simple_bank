package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"TestProj/api"
	db "TestProj/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://admin:password@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, _ := sql.Open(dbDriver, dbSource)

	err := conn.Ping()

	if err != nil {
		log.Fatal("DB connect err")
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("Failed to start server")
	}
}
