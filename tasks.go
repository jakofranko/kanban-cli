package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
)

// The number of swim lanes should be dynamic, as well as find their
// appropriate items in the task db according to this list only.
type status int

const (
	todo status = iota
	inProgress
	done
)

var status_strings = [...]string{"todo", "in progress", "done"}

func (s status) String() string {
	return status_strings[s]
}

func GetStatusFromString(s string) (int, error) {
	for i, v := range status_strings {
		if v == s {
			return i, nil
		}
	}

	return -1, errors.New(fmt.Sprintf("Invalid status %s", s))
}

type Task struct {
	Id      int
	Name    string
	Info    string
	Status  status
	Project string
}

type CreateTaskMsg struct {
	task Task
}

type EditTaskMsg struct {
	task  Task
	index int
}

type DeleteTaskMsg struct {
	index int
}

func (t *Task) Next() {
	if t.Status == done {
		t.Status = todo
	} else {
		t.Status++
	}
}

// Implement the bubbles/list.Item interface
func (t Task) FilterValue() string {
	return t.Name
}

func (t Task) Title() string {
	return t.Name
}

func (t Task) Description() string {
	return t.Info
}

func (t *Task) Merge(newT Task) {
	newValues := reflect.ValueOf(&newT).Elem()
	oldValues := reflect.ValueOf(t).Elem()

	// Loop through new values and assign them to a new task
	for i := 0; i < newValues.NumField(); i++ {
		fieldName := newValues.Type().Field(i).Name
		newField := newValues.Field(i).Interface()

		// Ignore ID fields
		if fieldName == "Id" {
			continue
		}

		if oldValues.CanSet() {
			if v, ok := newField.(int64); ok && newField != 0 {
				oldValues.Field(i).SetInt(v)
				continue
			}

			if v, ok := newField.(string); ok && newField != "" {
				oldValues.Field(i).SetString(v)
				continue
			}

			if v, ok := newField.(status); ok {
				oldValues.Field(i).SetInt(int64(v))
				continue
			}

			log.Printf("Unsupported value for %s : %T", newField, newField)
		}
	}
}

func NewTask(status status, name string, info string, id int, project string) Task {
	return Task{Status: status, Name: name, Info: info, Id: id, Project: project}
}

type TaskDB struct {
	db *sql.DB
}

func (t *TaskDB) CreateTable() error {
	// Create our table if it doesn't exist.
	// A task should have the following data:
	// id
	// name
	// info
	// status
	// project --> this is for a future project
	createStatement := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        info TEXT,
        status INTEGER,
        project TEXT
    )
    `

	_, err := t.db.Exec(createStatement)

	return err
}

func (t *TaskDB) Insert(name, info, project string, status status) (sql.Result, error) {
	result, err := t.db.Exec(
		"INSERT INTO tasks (name, info, status, project) VALUES(?, ?, ?, ?)",
		name,
		info,
		status,
		project,
	)

	return result, err
}

func (t *TaskDB) Delete(id int) error {
	_, err := t.db.Exec(
		"DELETE FROM tasks WHERE id = ?",
		id,
	)

	return err
}

func (t *TaskDB) Get(id int) (Task, error) {
	var task Task
	err := t.db.QueryRow("SELECT * FROM tasks WHERE id = ?", id).
		Scan(
			&task.Id,
			&task.Name,
			&task.Info,
			&task.Status,
			&task.Project,
		)

	return task, err
}

func (t *TaskDB) Update(task Task) error {
	// Get current task
	curr, err := t.Get(task.Id)
	if err != nil {
		return err
	}

	// Mutate current task with updated task values
	curr.Merge(task)

	// Perform the update
	_, mErr := t.db.Exec(
		"UPDATE tasks SET name = ?, info = ?, status = ?, project = ? WHERE id = ?",
		curr.Name,
		curr.Info,
		curr.Status,
		curr.Project,
		curr.Id,
	)

	return mErr
}

func (t *TaskDB) NextStatus(task Task) (Task, error) {
	// First, increment the task itself
	task.Next()

	// Then, update this task in the DB
	err := t.Update(task)
	if err != nil {
		return Task{}, err
	}

	// Finally, return the task with the incremented status
	return task, nil
}

func (t *TaskDB) GetAll() ([]Task, error) {
	var tasks []Task

	rows, err := t.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var task Task
		err = rows.Scan(
			&task.Id,
			&task.Name,
			&task.Info,
			&task.Status,
			&task.Project,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, err
}

func (t *TaskDB) GetByStatus(status status, project string) ([]Task, error) {
	var tasks []Task

	rows, err := t.db.Query("SELECT * FROM tasks WHERE status = ? AND project = ?", status, project)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var task Task
		err = rows.Scan(
			&task.Id,
			&task.Name,
			&task.Info,
			&task.Status,
			&task.Project,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, err
}

func (t *TaskDB) GetUniqueProjectNames() ([]string, error) {
	var projects []string

	rows, err := t.db.Query("SELECT DISTINCT project FROM tasks")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var project string
		err = rows.Scan(&project)
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, err
}

type ProjectTasksByStatusRow struct {
	status status
	id     int
	count  int
}

func (t *TaskDB) GetProjectTasksByStatus(projectName string) ([]ProjectTasksByStatusRow, error) {
	rows, err := t.db.Query("SELECT status, id, COUNT(id) FROM tasks WHERE project = ? GROUP BY status", projectName)
	if err != nil {
		return nil, err
	}

	var tasks []ProjectTasksByStatusRow
	for rows.Next() {
		var task ProjectTasksByStatusRow
		err = rows.Scan(
			&task.status,
			&task.id,
			&task.count,
		)

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t *TaskDB) AddNewProject(projectName string) error {
	_, err := t.db.Exec("INSERT INTO tasks (name, info, status, project) VALUES('A new beginning', '', 0, ?)", projectName)

	return err
}
