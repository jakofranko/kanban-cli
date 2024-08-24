package main

import (
	"database/sql"
	"log"
)

type ProjectDB struct {
	db *sql.DB
}

type projectStatus int

const (
	open projectStatus = iota
	archived
)

func (p *ProjectDB) CreateTable() error {
	// Create our table if it doesn't exist.
	// A project should have the following data:
	// id
	// name
	// sort_order
	// status
	createStatement := `
    CREATE TABLE IF NOT EXISTS projects (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        sort_order INTEGER,
        status INTEGER DEFAULT 0
    )
    `

	_, err := p.db.Exec(createStatement)

	return err
}

func (p *ProjectDB) GetAll() ([]Project, error) {
	rows, err := p.db.Query("SELECT * FROM projects")
	if err != nil {
		return nil, err
	}

	var projects []Project
	for rows.Next() {
		var project Project
		rows.Scan(
			&project.id,
			&project.name,
			&project.order,
			&project.status,
		)

		projects = append(projects, project)
	}

	return projects, nil
}

func (p *ProjectDB) GetByStatus(s projectStatus) ([]Project, error) {
	rows, err := p.db.Query("SELECT * FROM projects WHERE status = ?", s)
	if err != nil {
		return nil, err
	}

	var projects []Project
	for rows.Next() {
		var project Project
		rows.Scan(
			&project.id,
			&project.name,
			&project.order,
			&project.status,
		)

		projects = append(projects, project)
	}

	return projects, nil
}

func (p *ProjectDB) GetHighestOrder() (int, error) {
	row := p.db.QueryRow("SELECT sort_order FROM projects ORDER BY sort_order DESC")

	err := row.Err()
	if err != nil {
		return 0, err
	}

	var highestOrder int
	row.Scan(&highestOrder)
	return highestOrder, nil
}

func (p *ProjectDB) Insert(projectName string) (sql.Result, error) {
	newOrder, err := p.GetHighestOrder()
	if err != nil {
		log.Fatal(err)
	}

	result, err := p.db.Exec("INSERT INTO projects (name, sort_order) VALUES(?, ?)", projectName, newOrder)
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}

func (p *ProjectDB) ArchiveProject(id int) error {
	_, err := p.db.Exec("UPDATE projects SET status = ? WHERE id = ?", archived, id)
	if err != nil {
		return err
	}

	return nil
}
