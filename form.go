package main

import (
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var titleStyle = lipgloss.NewStyle().
	Padding(1).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(secondaryColor).
	Align(lipgloss.Center)

var descStyle = titleStyle

// Form Model
type Form struct {
	focused     status
	editing     bool
	index       int // Index within current list
	title       textinput.Model
	description textarea.Model
	project     int
	id          int // DB id of task
	keys        formKeyMap
	help        help.Model
	width       int
	height      int
}

func NewTitle() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "What is the task's title?"
	return ti
}

func TitleView(title textinput.Model) string {
	return titleStyle.Render(title.View())
}

func DescView(title textarea.Model) string {
	return descStyle.Render(title.View())
}

func NewDescription() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Brief description"
	return ta
}

func NewForm(width, height int, focused status, project int) *Form {
	form := &Form{width: width, height: height, focused: focused, keys: formKeys, help: help.New()}
	form.title = NewTitle()
	form.description = NewDescription()
	form.editing = false
	form.project = project

	form.title.Focus()
	return form
}

func UpdateForm(task Task, index int) *Form {
	form := &Form{focused: task.Status, keys: formKeys, help: help.New()}
	form.title = NewTitle()
	form.description = NewDescription()
	form.editing = true
	form.index = index

	form.title.SetValue(task.Name)
	form.description.SetValue(task.Info)
	form.id = task.Id
	form.project = task.ProjectId

	form.title.Focus()
	return form
}

func (m Form) Init() tea.Cmd {
	return nil
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Next):
			if m.title.Focused() {
				m.title.Blur()
				m.description.Focus()
				return m, textarea.Blink
			} else {
				// Insert new task into db
				models[form] = m

				if m.editing {
					return models[board], m.UpdateTask
				}

				return models[board], m.CreateTask
			}
		case key.Matches(msg, m.keys.Back):
			if m.description.Focused() {
				m.description.Blur()
				m.title.Focus()
				return m, textinput.Blink
			} else {
				return models[board], nil
			}
		}
	}

	// Pass all other key presses to the inputs
	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
		return m, cmd
	} else {
		m.description, cmd = m.description.Update(msg)
		return m, cmd
	}
}

func (m Form) View() string {
	render := lipgloss.JoinVertical(
		lipgloss.Center,
		TitleView(m.title),
		DescView(m.description),
		m.help.View(formKeys),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, render)
}

func (m Form) CreateTask() tea.Msg {
	task := NewTask(m.focused, m.title.Value(), m.description.Value(), 0, m.project)

	// Insert task into db
	taskDB := GetDB()
	defer taskDB.db.Close()
	newTask, err := taskDB.Insert(task.Name, task.Info, task.ProjectId, task.Status)
	if err != nil {
		log.Fatal(err)
	}

	// Assign the new ID to the new task, so it can be actioned
	// in the board without taking a large poopoo.
	var newId int64
	newId, err = newTask.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	task.Id = int(newId)

	// Return create task message
	return CreateTaskMsg{task: task}
}

func (m Form) UpdateTask() tea.Msg {
	task := NewTask(m.focused, m.title.Value(), m.description.Value(), m.id, m.project)

	// Update task in db
	taskDB := GetDB()
	defer taskDB.db.Close()
	taskDB.Update(task)

	return EditTaskMsg{task: task, index: m.index}
}
