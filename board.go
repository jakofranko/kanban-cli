package main

import (
    tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const divisor = 4

// Styles
var (
    helpStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241"))
)

// This is the model for the board view, which will implement the
// Bubbletea model methods for rendering etc. It maintains multiple
// lists components from bubbles, and is styled via lipgloss.
type Board struct {
    focused status
    lanes []SwimLane
	err  error
    loaded      bool
    quitting bool
}

type UpdateListMsg struct {
    update status
    newBoard Board
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
    m.lanes[m.focused].Blur()

    if m.focused == done {
        m.focused = todo
    } else {
        m.focused++
    }

    m.lanes[m.focused].Focus()
}

func (m *Board) Prev() {
    m.lanes[m.focused].Blur()

    if m.focused == todo {
        m.focused = done
    } else {
        m.focused--
    }

    m.lanes[m.focused].Focus()
}

func (m *Board) MoveToNext() tea.Msg {
    // First, get the focused lane and selected task
    focusedLane := m.lanes[m.focused]
    selectedItem := focusedLane.list.SelectedItem()

    // Only act if there is a selected item
    if selectedItem != nil {
        // Cast the selected list Item to a Task
		selectedTask := selectedItem.(Task)

        // Get the index of the selection
        itemIndex := focusedLane.list.Index()
        
        // Get the new cursor index for the current lane (back it up)
        var selectIndex int
        selectIndex = itemIndex-1
        if selectIndex < 0 {
            selectIndex = 0
        }

        // Get the current (will be old) status
        oldStatus := selectedTask.status

        // Remove the item from the list
		m.lanes[oldStatus].list.RemoveItem(itemIndex)

        // Put the cursor on the new selectIndex value
        m.lanes[oldStatus].list.Select(selectIndex)

        // Increment the selected task status, which should put it on the next lane
		selectedTask.Next()

        // Get the new lane and insert the task into this lane
        // Leaving this here for learning:
        //
        // newLane := m.lanes[selectedTask.status]
        //
        // The reason this won't work is that it's no longer operating on the model.
        // Pulling this lane into a variable like this means I'm operating on
        // a new object, not the model itself.
		m.lanes[selectedTask.status].
            list.
            InsertItem(
                len(m.lanes[selectedTask.status].list.Items())+1, 
                list.Item(selectedTask),
            )

        // Update target list. The current list will get updated on the next
        // Update tick in the main Board.
        m.lanes[selectedTask.status].list.Update(nil)
        return UpdateListMsg{update: oldStatus}
    }

    return nil
}

func (m *Board) initLists(width, height int) {
    todoLane := new(SwimLane)
    inProgressLane := new(SwimLane)
    doneLane := new(SwimLane)
    m.lanes = []SwimLane{
        todoLane.Init(width, height, 4, "To Do", todoList),
        inProgressLane.Init(width, height, 4, "In Progress", inProgressList),
        doneLane.Init(width, height, 4, "Done", doneList),
    }
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
            m.focused = todo
            m.lanes[m.focused].Focus()
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
            case "e":
                models[board] = m
                currentList := m.lanes[m.focused]
                currentTask := currentList.list.SelectedItem().(Task)
                currentIndex := currentList.list.Index()
                models[form] = UpdateForm(m.focused, currentTask.Title(), currentTask.Description(), currentIndex)
                return models[form], nil
            case "d":
                // We could do a confirmation screen, but for now just delete.
                // Another option would be to archive items that are in the done
                // column. Maybe a feature for when I'm using persistant storage.
                i := m.lanes[m.focused].list.Index()
                m.lanes[m.focused].list.RemoveItem(i)
                return m, nil
        }
    case CreateTaskMsg:
        task := msg.task
        return m, m.lanes[task.status].list.InsertItem(len(m.lanes[task.status].list.Items()), task)
    case EditTaskMsg:
        task := msg.task
        i := msg.index
        return m, m.lanes[task.status].list.SetItem(i, task)
    case UpdateListMsg:
        listToUpdate := msg.update
        m.lanes[listToUpdate].list.Update(nil)
        return m, nil
	}

	var cmd tea.Cmd
	m.lanes[m.focused].list, cmd = m.lanes[m.focused].list.Update(msg)

	return m, cmd
}

func (m Board) View() string {
    if m.quitting {
        return "Quitting KanBan CLI..."
    }

    if m.loaded {
        todoView := m.lanes[todo].View()
        ipView := m.lanes[inProgress].View()
        doneView := m.lanes[done].View()
        
        return lipgloss.JoinHorizontal(
            lipgloss.Left, 
            todoView,
            ipView,
            doneView,
        )
    } else {
        return "loading..."
    }
}
