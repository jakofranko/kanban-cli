package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/gobeam/stringy"
)

const (
	divisor       = len(status_strings) + 1
	horizontalPad = 2
	verticalPad   = 1
	bordersize    = 1
)

// Styles
var (
	columnStyle = lipgloss.NewStyle().
			Padding(verticalPad, horizontalPad)
	focusedStyle = lipgloss.NewStyle().
			Padding(verticalPad, horizontalPad).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#C0FFE3"))
	listTitleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("67")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)
	listFocusItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("84")).
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.Color("84")).
				Padding(0, 0, 0, 1)
	listFocusItemDescStyle = listFocusItemStyle.Copy().
				Foreground(lipgloss.Color("71"))
)

type SwimLane struct {
	title      string
	focused    bool
	list       list.Model
	laneStatus status
}

func (s *SwimLane) Focus() {
	s.focused = true
}

func (s *SwimLane) Blur() {
	s.focused = false
}

// w should be the available width of the current view port (minus other UI)
func (s *SwimLane) SetWidth(w int) {
	hOffset := (horizontalPad * 2) + (bordersize * 2)
	s.list.SetHeight(w - hOffset)
}

// h should be the available height of the current view port (minus other UI)
func (s *SwimLane) SetHeight(h int) {
	vOffset := (verticalPad * 2) + (bordersize * 2)
	s.list.SetHeight(h - vOffset)
}

// This will create a new list, meant to be rendered next to N number of other lists,
// where N is equal to the number total lists. This number is passed in as a divisor.
func (s *SwimLane) Init(width int, height int, status status) SwimLane {
	s.laneStatus = status

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

	vOffset := (verticalPad * 2) + (bordersize * 2)
	hOffset := (horizontalPad * 2) + (bordersize * 2)

	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = listFocusItemStyle
	d.Styles.SelectedDesc = listFocusItemDescStyle

	s.list = list.New([]list.Item{}, d, width-hOffset, height-vOffset)
	s.list.Title = title
	s.list.SetItems(items)
	s.list.SetShowHelp(false)

	// Set styles
	s.list.Styles.Title = listTitleStyle

	return *s
}

func (s *SwimLane) View() string {
	if s.focused {
		return focusedStyle.Render(s.list.View())
	}

	return columnStyle.Render(s.list.View())
}
