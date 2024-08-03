package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/gobeam/stringy"
)

const (
	divisor = len(status_strings) + 1
	pad     = 2
)

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
	title   string
	focused bool
	list    list.Model
}

func (s *SwimLane) Focus() {
	s.focused = true
}

func (s *SwimLane) Blur() {
	s.focused = false
}

// This will create a new list, meant to be rendered next to N number of other lists,
// where N is equal to the number total lists. This number is passed in as a divisor.
func (s *SwimLane) Init(width int, height int, status status) SwimLane {
	title := stringy.New(status.String()).Title()
	s.title = title

	// Fetch items from the DB
	taskDB := GetDB()
	defer taskDB.db.Close()
	tasks, err := taskDB.GetByStatus(status)
	if err != nil {
		log.Fatal(err)
	}

	var items []list.Item
	for _, task := range tasks {
		items = append(items, task)
	}

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
