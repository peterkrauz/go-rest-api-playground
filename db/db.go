package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	HOST = "database"
	PORT = 5432
)

var ErrorNoItemFound = fmt.Errorf("no item found")

type Database struct {
	Connection *sql.DB
}

func Initialize(username, password, dbName string) (Database, error) {
	database := Database{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, dbName)
	connection, error := sql.Open("postgres", dsn)
	if error != nil {
		return database, error
	}
	database.Connection = connection
	log.Println("Successfully connected to Database")
	return database, nil
}
