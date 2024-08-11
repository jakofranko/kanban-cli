package main

import (
	"log"
	"math"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	keys     projectListKeyMap
	help     help.Model

	// Store these if this is the first view, and pass to subsequent models
	width  int
	height int
}

func (p *ProjectsTable) Init() tea.Cmd {
	return nil
}

func (p *ProjectsTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.height = msg.Height
		p.width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, p.keys.Quit):
			return p, tea.Quit
		case key.Matches(msg, p.keys.Up):
			if !p.table.Focused() {
				p.table.Focus()
			}
			p.table.MoveUp(1)
		case key.Matches(msg, p.keys.Down):
			if !p.table.Focused() {
				p.table.Focus()
			}
			p.table.MoveDown(1)
		case key.Matches(msg, p.keys.Select):
			row := p.table.SelectedRow()

			// Get a new kanban board for this project
			b := NewBoard(row[0], p.width, p.height)
			models[projects] = p
			models[board] = b
			return models[board], nil
		}
	}
	return p, nil
}

func (p *ProjectsTable) View() string {
	return lipgloss.JoinVertical(lipgloss.Center, p.table.View(), p.help.View(p.keys))
}

func NewProjectsTable() *ProjectsTable {
	// Fetch unique project names from the TaskDB
	taskDB := GetDB()
	defer taskDB.db.Close()
	projectNames, err := taskDB.GetUniqueProjectNames()
	if err != nil {
		log.Fatal(err)
	}

	// The column should be at least as wide as the column title
	const projectTitle = "Project Name"
	var longestProjectName int
	longestProjectName = len(projectTitle)

	var rows []table.Row

	for _, pn := range projectNames {
		flpn := float64(longestProjectName)
		fpn := float64(len(pn))
		lpn := math.Max(flpn, fpn)
		longestProjectName = int(lpn)

		tasks, err := taskDB.GetProjectTasksByStatus(pn)
		if err != nil {
			log.Fatal(err)
		}

		// build the row
		var row table.Row
		row = append(row, pn)

		// Add the tasks to the appropriate columns
		var statuses [3]string
		for _, task := range tasks {
			statuses[task.status] = strconv.Itoa(task.count)
		}

		row = append(row, statuses[0:3]...)

		rows = append(rows, row)
	}

	columns := []table.Column{
		{Title: projectTitle, Width: longestProjectName},
		{Title: "Todo", Width: 4},
		{Title: "In Progress", Width: 11},
		{Title: "Done", Width: 4},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	return &ProjectsTable{
		table: t,
		keys:  projectListKeys,
	}
}
