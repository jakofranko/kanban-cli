package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type status int

const (
    todo status = iota
    inProgress
    done
)

type Task struct {
	title       string
	description string
	status      status
}

func (t *Task) Next() {
    if t.status == done {
        t.status = todo
    } else {
        t.status++
    }
}

// Implement the bubbles/list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

func NewTask(status status, title string, description string) Task {
    return Task{status: status, title: title, description: description}
}

func (m Form) CreateTask() tea.Msg {
    task := NewTask(m.focused, m.title.Value(), m.description.Value())
    return task
}
