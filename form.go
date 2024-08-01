package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var titleStyle = lipgloss.NewStyle().
	Padding(1).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("33"))

var descStyle = titleStyle

// Form Model
type Form struct {
	focused     status
	title       textinput.Model
	description textarea.Model
	editing     bool
	index       int
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

func NewForm(focused status) *Form {
	form := &Form{focused: focused}
	form.title = NewTitle()
	form.description = NewDescription()
	form.editing = false

	form.title.Focus()
	return form
}

func UpdateForm(focused status, title string, description string, index int) *Form {
	form := &Form{focused: focused}
	form.title = NewTitle()
	form.description = NewDescription()
	form.editing = true
	form.index = index

	form.title.SetValue(title)
	form.description.SetValue(title)

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
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+y":
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
	return lipgloss.JoinVertical(lipgloss.Left, TitleView(m.title), DescView(m.description))
}

func (m Form) CreateTask() tea.Msg {
	task := NewTask(m.focused, m.title.Value(), m.description.Value())

	// Insert task into db
	// TODO

	// Return create task message
	return CreateTaskMsg{task: task}
}

func (m Form) UpdateTask() tea.Msg {
	task := NewTask(m.focused, m.title.Value(), m.description.Value())
	return EditTaskMsg{task: task, index: m.index}
}
