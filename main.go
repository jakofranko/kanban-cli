package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const divisor = 4

// Styles
var (
    columnStyle = lipgloss.NewStyle().Padding(1, 2)
    focusedStyle = lipgloss.NewStyle().
        Padding(1, 2).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("62"))
    helpStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241"))
    )

// Type defs
type status int

const (
    todo status = iota
    inProgress
    done
)
type Task struct {
	title       string
	description string
	status      status
}

func (t *Task) Next() {
    if t.status == done {
        t.status = todo
    } else {
        t.status++
    }
}

// Implement list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

// Main model
type Model struct {
    focused status
	lists []list.Model
	err  error
    loaded      bool
    quitting bool
}

func New() *Model {
	return &Model{}
}

func (m *Model) Next() {
    if m.focused == done {
        m.focused = todo
    } else {
        m.focused++
    }
}

func (m *Model) Prev() {
    if m.focused == todo {
        m.focused = done
    } else {
        m.focused--
    }
}

func (m *Model) MoveToNext() tea.Msg {
    focusedItem := m.lists[m.focused]
    selectedItem := focusedItem.SelectedItem()
    if selectedItem != nil {
		selectedTask := selectedItem.(Task)
		m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
        _, updateCurrentListCmd := m.lists[selectedTask.status].Update(nil)

		selectedTask.Next()
		m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
        _, updateNextListCmd := m.lists[selectedTask.status].Update(nil)

        tea.Sequence(updateNextListCmd, updateCurrentListCmd)
    }

    return nil
}

// TODO: call on tea.WindowSizeMsg
func (m *Model) initLists(width, height int) {
    defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height-divisor)
    defaultList.SetShowHelp(false)

    m.lists = []list.Model{defaultList, defaultList, defaultList}
    
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "the good stuff"},
		Task{status: todo, title: "get gud", description: "the good stuff"},
		Task{status: todo, title: "this is cool", description: "the good stuff"},
	})
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: inProgress, title: "Get good at Go", description: "I'm working on it"},
	})
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: done, title: "Be good at go", description: "Don't question it"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
        if !m.loaded {
            columnStyle.Width(msg.Width / divisor)
            focusedStyle.Width(msg.Width / divisor)
            columnStyle.Height(msg.Height - divisor)
            focusedStyle.Height(msg.Height - divisor)
            m.initLists(msg.Width, msg.Height)
            m.loaded = true
        }
    case tea.KeyMsg:
        switch msg.String() {
            case "ctrl+c", "q":
                m.quitting = true
                return m, tea.Quit
            case "left", "h":
                m.Prev()
            case "right", "l":
                m.Next()
            case "enter":
                return m, m.MoveToNext
        }
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)

	return m, cmd
}

func (m Model) View() string {
    if m.quitting {
        return "Quitting KanBan CLI..."
    }

    if m.loaded {
        todoView := m.lists[todo].View()
        ipView := m.lists[inProgress].View()
        doneView := m.lists[done].View()
        
        switch m.focused {
            case todo:
                return lipgloss.JoinHorizontal(
                    lipgloss.Left, 
                    focusedStyle.Render(todoView),
                    columnStyle.Render(ipView),
                    columnStyle.Render(doneView),
                )
            case inProgress:
                return lipgloss.JoinHorizontal(
                    lipgloss.Left, 
                    columnStyle.Render(todoView),
                    focusedStyle.Render(ipView),
                    columnStyle.Render(doneView),
                )
            case done:
                return lipgloss.JoinHorizontal(
                    lipgloss.Left, 
                    columnStyle.Render(todoView),
                    columnStyle.Render(ipView),
                    focusedStyle.Render(doneView),
                )
            default:
                return lipgloss.JoinHorizontal(
                    lipgloss.Left, 
                    focusedStyle.Render(todoView),
                    columnStyle.Render(ipView),
                    columnStyle.Render(doneView),
                )
        }
    } else {
        return "loading..."
    }
}

func main() {
	m := New()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
