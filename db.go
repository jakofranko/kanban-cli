package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	gap "github.com/muesli/go-app-paths"
)

const (
	dbDriver = "sqlite3"
	dbName   = "./kanbandb"
)

// Get or Setup XDG-compliant path for SQLite DB
func getDbPath() string {
	// Get XDG paths
	scope := gap.NewScope(gap.User, "kanban")
	dirs, err := scope.DataDirs()
	if err != nil {
		log.Fatal(err)
	}

	// Create directory if it doesn't exist
	var kanbanDir string
	if len(dirs) > 0 {
		kanbanDir = dirs[0]
	} else {
		kanbanDir, _ = os.UserHomeDir()
	}

	if err := initKanbanDir(kanbanDir); err != nil {
		log.Fatal(err)
	}

	return kanbanDir
}

func initKanbanDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o770)
		}

		return err
	}

	return nil
}

func GetDB() TaskDB {
    // Uncomment for local dev
	// db, err := sql.Open(dbDriver, dbName)

    // Comment for local dev
	db, err := sql.Open(dbDriver, filepath.Join(getDbPath(), dbName))
	if err != nil {
		log.Fatal(err)
	}

	// Initializse our TaskDB
	t := TaskDB{db}

	// This will create the table if it does not exist
	err = t.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	return t
}

func GetProjectDB() ProjectDB {
	db, err := sql.Open(dbDriver, dbName)
	// db, err := sql.Open(dbDriver, filepath.Join(getDbPath(), dbName))
	if err != nil {
		log.Fatal(err)
	}

	p := ProjectDB{db}

	// This will create the table if it does not exist
	err = p.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	return p
}
