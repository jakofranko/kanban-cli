package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const pad = 2

// Styles
var (
    columnStyle = lipgloss.NewStyle().
        Padding(1, pad)
    focusedStyle = lipgloss.NewStyle().
        Padding(1, pad).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("62"))
)

type SwimLane struct {
    title string
    focused bool
    list list.Model
}

func (s *SwimLane) Focus() {
    s.focused = true
}

func (s *SwimLane) Blur() {
    s.focused = false
}

// This will create a new list, meant to be rendered next to N number of other lists,
// where N is equal to the number total lists. This number is passed in as a divisor.
func (s *SwimLane) Init(width int, height int, divisor int, title string, items []list.Item) SwimLane {
    s.title = title

    s.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height-(pad*2))
    s.list.Title = title
    s.list.SetItems(items)
    s.list.SetShowHelp(false)

    return *s
}

func (s *SwimLane) View() string {
    if s.focused {
        return focusedStyle.Render(s.list.View())
    }

    return columnStyle.Render(s.list.View())
}
