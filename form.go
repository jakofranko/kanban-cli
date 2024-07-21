package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Form Model
type Form struct {
    focused status
    title textinput.Model
    description textarea.Model
}


func NewTitle() textinput.Model {
    ti := textinput.New()
    ti.Placeholder = "What is the task's title?"
    return ti
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

    form.title.Focus()
    return form
}

func (m Form) Init() tea.Cmd {
    return nil
}

func NewTask(status status, title string, description string) Task {
    return Task{status: status, title: title, description: description}
}

func (m Form) CreateTask() tea.Msg {
    task := NewTask(m.focused, m.title.Value(), m.description.Value())
    return task
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch  msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "enter":
            if m.title.Focused() {
                m.title.Blur()
                m.description.Focus()
                return m, textarea.Blink
            } else {
                models[form] = m
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
    return lipgloss.JoinVertical(lipgloss.Left, m.title.View(), m.description.View())
}
