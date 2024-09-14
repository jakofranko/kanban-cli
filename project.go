package main

import (
	"log"
	"math"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var tableStyle = lipgloss.NewStyle().
	Margin(2, 1)

var newProjectStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Border(lipgloss.RoundedBorder(), true).
	BorderForeground(lipgloss.Color("#C0FF3E"))

var centerCtyle = lipgloss.NewStyle().
	Align(lipgloss.Center)

type Project struct {
	id        int
	name      string
	order     int
	status    projectStatus
	todoTasks []Task
	ipTasks   []Task
	doneTasks []Task
}

type ProjectsTable struct {
	projects []Project
	table    table.Model
	keys     projectListKeyMap
	help     help.Model
	view     projectStatus

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
		p.setViewSize(msg.Height)
	case RefreshProjectsMsg:
		columns, rows := buildTable(p.view)
		p.table.SetColumns(columns)
		p.table.SetRows(rows)
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
			pId, err := strconv.Atoi(row[0])
			if err != nil {
				log.Fatal(err)
				return nil, nil
			}

			b := NewBoard(pId, p.width, p.height)
			models[projects] = p
			models[board] = b
			return models[board], nil
		case key.Matches(msg, p.keys.Help):
			p.help.ShowAll = !p.help.ShowAll
			p.help.Update(nil)
			p.setViewSize(p.height)
		case key.Matches(msg, p.keys.New):
			f := NewProjectForm(p.width, p.height)
			return f, nil
		case key.Matches(msg, p.keys.Archive):
			i := p.table.Cursor()
			if i >= 0 {
				projectDB := GetProjectDB()
				defer projectDB.db.Close()

				row := p.table.SelectedRow()
				pId, err := strconv.Atoi(row[0])
				if err != nil {
					log.Fatal(err)
				}

				projectDB.ArchiveProject(pId)

				// Remove row from this view
				rows := p.table.Rows()
				newRows := append(rows[:i], rows[i+1:]...)
				p.table.SetRows(newRows)
				p.table.SetCursor(i - 1)
			}
		case key.Matches(msg, p.keys.ViewArchived):
			if p.view == open {
				p.view = archived
			} else {
				p.view = open
			}

			columns, rows := buildTable(p.view)
			p.table.SetColumns(columns)
			p.table.SetRows(rows)
		}
	}

	return p, nil
}

func (p *ProjectsTable) View() string {
	var heading string
	if p.view == archived {
		heading = "(Archived) Projects"
	} else {
		heading = "Projects"
	}

	he := centerCtyle.Width(p.width).Render(heading)
	t := tableStyle.Width(p.width).Align(lipgloss.Center).Render(p.table.View())
	h := helpStyle.Width(p.width).Align(lipgloss.Center).Render(p.help.View(p.keys))
	return lipgloss.JoinVertical(lipgloss.Left, he, t, h)
}

func (p *ProjectsTable) setViewSize(height int) {
	// Get all UI elements in view
	h := lipgloss.Height(p.help.View(p.keys))
	heading := 1

	// There is a magic number of height added
	// that is equal to the padding, margin, and border.
	// I don't know a better way to pull this out progromatically.
	p.table.SetHeight(height - h - heading - 6)
}

func buildTable(s projectStatus) ([]table.Column, []table.Row) {
	projectDB := GetProjectDB()
	defer projectDB.db.Close()

	// Get all projects from project db
	projects, err := projectDB.GetByStatus(s)
	if err != nil {
		log.Fatal(err)
	}

	// The column should be at least as wide as the column title
	const projectTitle = "Project Name"
	var longestProjectName int
	longestProjectName = len(projectTitle)

	taskDB := GetDB()
	defer taskDB.db.Close()

	var rows []table.Row

	for _, p := range projects {
		flpn := float64(longestProjectName)
		fpn := float64(len(p.name))
		lpn := math.Max(flpn, fpn)
		longestProjectName = int(lpn)

		tasks, err := taskDB.GetProjectTasksByStatus(p.id)
		if err != nil {
			log.Fatal(err)
		}

		// build the row
		var row table.Row
		row = append(row, strconv.Itoa(p.id), p.name)

		// Add the tasks to the appropriate columns
		statuses := [3]string{"0", "0", "0"}
		for _, task := range tasks {
			statuses[task.status] = strconv.Itoa(task.count)
		}

		row = append(row, statuses[0:3]...)

		rows = append(rows, row)
	}

	columns := []table.Column{
		{Title: "ID", Width: 2},
		{Title: projectTitle, Width: longestProjectName},
		{Title: "Todo", Width: 4},
		{Title: "In Progress", Width: 11},
		{Title: "Done", Width: 4},
	}

	return columns, rows
}

func NewProjectsTable() *ProjectsTable {
	columns, rows := buildTable(open)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	// Set table styles by extracting defaults, and the resetting them
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("241")).
		BorderRight(true).
		BorderBottom(true)

	s.Cell = s.Cell.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("241")).
		BorderRight(true)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#C0FF3E")).
		Bold(true)

	t.SetStyles(s)

	return &ProjectsTable{
		table: t,
		keys:  projectListKeys,
		help:  help.New(),
	}
}

type NewProject struct {
	model  textinput.Model
	name   string
	width  int
	height int
}

func NewProjectForm(width, height int) *NewProject {
	t := textinput.New()
	t.Placeholder = "Project Name"
	t.Focus()
	f := &NewProject{model: t, width: width, height: height}

	return f
}

func (f *NewProject) Init() tea.Cmd {
	return textinput.Blink
}

func (f *NewProject) View() string {
	render := newProjectStyle.Render(f.model.View())
	return lipgloss.Place(f.width, f.height, lipgloss.Center, lipgloss.Center, render)
}

func (f *NewProject) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return f, tea.Quit
		case "ctrl+b":
			return models[projects], nil
		case "enter":
			projectDB := GetProjectDB()
			defer projectDB.db.Close()

			_, err := projectDB.Insert(f.model.Value())
			if err != nil {
				log.Fatal(err)
			}

			return models[projects], f.RefreshProjects
		default:
			// Pass all keystrokes to textinput
			f.model, cmd = f.model.Update(msg)

			return f, cmd
		}
	}

	return f, nil
}

// Simple message to tell the project model to build the rows again
type RefreshProjectsMsg struct{}

func (f *NewProject) RefreshProjects() tea.Msg {
	return RefreshProjectsMsg{}
}
