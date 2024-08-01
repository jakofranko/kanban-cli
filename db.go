package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbDriver = "sqlite3"
	dbName   = "./kanbandb"
)

func GetDB() TaskDB {
	db, err := sql.Open(dbDriver, dbName)
	if err != nil {
		log.Fatal(err)
	}

	// Initializse our TaskDB
	t := TaskDB{db}

	// This will create the table if it does not exist
	t.CreateTable()

	return t
}
