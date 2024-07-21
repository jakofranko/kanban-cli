package main

import (
    tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/list"
)

const divisor = 4

// Styles
var (
    columnStyle = lipgloss.NewStyle().
        Padding(1, 2)
    focusedStyle = lipgloss.NewStyle().
        Padding(1, 2).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("62"))
    helpStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241"))
)

// This is the model for the board view, which will implement the
// Bubbletea model methods for rendering etc. It maintains multiple
// lists components from bubbles, and is styled via lipgloss.
type Board struct {
    focused status
	lists []list.Model
	err  error
    loaded      bool
    quitting bool
}

type UpdateListMsg struct {
    update status
}

var todoList = []list.Item{
		Task{status: todo, title: "buy milk", description: "the good stuff"},
		Task{status: todo, title: "get gud", description: "the good stuff"},
		Task{status: todo, title: "this is cool", description: "the good stuff"},
}
var inProgressList = []list.Item{
		Task{status: inProgress, title: "Get good at Go", description: "I'm working on it"},
}
var doneList = []list.Item{
		Task{status: done, title: "Be good at go", description: "Don't question it"},
}

func NewBoard() *Board {
	return &Board{}
}

func (m *Board) Next() {
    if m.focused == done {
        m.focused = todo
    } else {
        m.focused++
    }
}

func (m *Board) Prev() {
    if m.focused == todo {
        m.focused = done
    } else {
        m.focused--
    }
}

// This is not working
func (m *Board) MoveToNext() tea.Msg {
    focusedItem := m.lists[m.focused]
    selectedItem := focusedItem.SelectedItem()
    if selectedItem != nil {
		selectedTask := selectedItem.(Task)
        oldStatus := selectedTask.status
        itemIndex := m.lists[m.focused].Index()
        var selectIndex int
        selectIndex = itemIndex-1
        if selectIndex < 0 {
            selectIndex = 0
        }

		m.lists[selectedTask.status].RemoveItem(itemIndex)
        m.lists[selectedTask.status].Select(selectIndex)

		selectedTask.Next()
		m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))

        // Update target list. The current list will get updated on the next
        // Update tick in the main Board.
        m.lists[selectedTask.status].Update(nil)
        return UpdateListMsg{update: oldStatus}
    }

    return nil
}

func (m *Board) initLists(width, height int) {
    defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height-divisor)
    defaultList.SetShowHelp(false)

    m.lists = []list.Model{defaultList, defaultList, defaultList}
    
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems(todoList)

	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems(inProgressList)

	m.lists[done].Title = "Done"
	m.lists[done].SetItems(doneList)
}

func (m Board) Init() tea.Cmd {
	return nil
}

func (m Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
            case "n":
               models[board] = m // save current model
               models[form] = NewForm(m.focused)
               return models[form], nil
        }
    case Task:
        task := msg
        return m, m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)
    case UpdateListMsg:
        listToUpdate := msg.update
        m.lists[listToUpdate].Update(nil)
        return m, nil
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)

	return m, cmd
}

func (m Board) View() string {
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
