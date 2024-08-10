package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Project struct {
	name      string
	todoTasks []Task
	ipTasks   []Task
	doneTasks []Task
}

type ProjectsTable struct {
	projects []Project
	table    table.Model
}

func (p *ProjectsTable) Init() tea.Cmd {
	return nil
}

func (p *ProjectsTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return p, nil
}

func (p *ProjectsTable) View() string {
	return p.table.View()
}

func NewProjectsTable() *ProjectsTable {
	columns := []table.Column{
		{Title: "Project", Width: 10},
		{Title: "Todo", Width: 4},
		{Title: "In Progress", Width: 11},
		{Title: "Done", Width: 4},
	}

	rows := []table.Row{
		{"Test Project", "4", "5", "6"},
		{"Test Project 2", "6", "7", "8"},
		{"Test Project 3", "5", "6", "7"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	return &ProjectsTable{table: t}
}
