package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var taskStyle = lipgloss.NewStyle().
	Padding(1).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(secondaryColor).
	Align(lipgloss.Center)

var nameStyle = lipgloss.NewStyle().
	MarginBottom(2).
	Bold(true)

var infoStyle = lipgloss.NewStyle().
	Italic(true)

type ViewTask struct {
	width  int
	height int
	task   Task
	help   help.Model
	keys   viewTaskKeyMap
}

func NewViewTask(width, height int, t Task) *ViewTask {
	model := &ViewTask{
		width:  width,
		height: height,
		task:   t,
		help:   help.New(),
		keys:   viewTaskKeys,
	}

	return model
}

func (v ViewTask) Init() tea.Cmd {
	return nil
}

func (v ViewTask) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch mt := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(mt, v.keys.Back):
			return models[board], nil
		case key.Matches(mt, v.keys.Quit):
			return v, tea.Quit
		}
	}

	return v, nil
}

func (v ViewTask) View() string {
	n := nameStyle.Render(v.task.Name)
	i := infoStyle.Render(v.task.Info)
	taskData := taskStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, n, i),
	)

	render := lipgloss.JoinVertical(
		lipgloss.Center,
		taskData,
		v.help.View(v.keys),
	)

	return lipgloss.Place(v.width, v.height, lipgloss.Center, lipgloss.Center, render)
}
