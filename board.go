package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
	focused  status
	lanes    []SwimLane
	err      error
	loaded   bool
	quitting bool
	project  string
}

type UpdateListMsg struct {
	update   status
	newBoard Board
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
		selectIndex = itemIndex - 1
		if selectIndex < 0 {
			selectIndex = 0
		}

		// Get the current (will be old) status
		oldStatus := selectedTask.Status

		// Remove the item from the list
		m.lanes[oldStatus].list.RemoveItem(itemIndex)

		// Put the cursor on the new selectIndex value
		m.lanes[oldStatus].list.Select(selectIndex)

		// Increment the selected task status, via the taskDB method,
		// which should put it on the next lane in the UI AND the db.
		taskDB := GetDB()
		defer taskDB.db.Close()

		updatedTask, err := taskDB.NextStatus(selectedTask)
		if err != nil {
			log.Fatal(err)
		}

		selectedTask = updatedTask

		// Get the new lane and insert the task into this lane
		// Leaving this here for learning:
		//
		// newLane := m.lanes[selectedTask.status]
		//
		// The reason this won't work is that it's no longer operating on the model.
		// Pulling this lane into a variable like this means I'm operating on
		// a new object, not the model itself.
		m.lanes[selectedTask.Status].
			list.
			InsertItem(
				len(m.lanes[selectedTask.Status].list.Items())+1,
				list.Item(selectedTask),
			)

		// Update target list. The current list will get updated on the next
		// Update tick in the main Board.
		m.lanes[selectedTask.Status].list.Update(nil)
		return UpdateListMsg{update: oldStatus}
	}

	return nil
}

func (m *Board) initLists(width, height int) {
	todoLane := new(SwimLane)
	inProgressLane := new(SwimLane)
	doneLane := new(SwimLane)

	// Get lists from the TaskDB by status.
	// For now though, just use empty lists.
	var todoList []list.Item
	var inProgressList []list.Item
	var doneList []list.Item

	taskDB := GetDB()
	defer taskDB.db.Close()

	todoRows, err := taskDB.GetByStatus(todo)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range todoRows {
		todoList = append(todoList, row)
	}

	ipRows, err := taskDB.GetByStatus(inProgress)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range ipRows {
		inProgressList = append(inProgressList, row)
	}

	doneRows, err := taskDB.GetByStatus(done)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range doneRows {
		doneList = append(doneList, row)
	}

	m.lanes = []SwimLane{
		todoLane.Init(width, height, 4, "To Do", todoList),
		inProgressLane.Init(width, height, 4, "In Progress", inProgressList),
		doneLane.Init(width, height, 4, "Done", doneList),
	}
}

func (m Board) Init() tea.Cmd {
	// TODO set project when initializing a new board
	// from a as of yet non-existent project view
	m.project = "test project"
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
			models[form] = NewForm(m.focused, m.project)

			return models[form], nil
		case "e":
			models[board] = m // save current model
			currentTask := m.lanes[m.focused].list.SelectedItem().(Task)
			currentIndex := m.lanes[m.focused].list.Index()
			models[form] = UpdateForm(currentTask, currentIndex)

			return models[form], nil
		case "d":
			// We could do a confirmation screen, but for now just delete.
			// Another option would be to archive items that are in the done
			// column. Maybe a feature for when I'm using persistant storage.
			i := m.lanes[m.focused].list.Index()
			m.lanes[m.focused].list.RemoveItem(i)

			// And remove from DB
			taskDB := GetDB()
			defer taskDB.db.Close()

			currentList := m.lanes[m.focused]
			task := currentList.list.SelectedItem().(Task)

			taskDB.Delete(task.Id)
			return m, nil
		}
	case CreateTaskMsg:
		task := msg.task

		// Insert into list
		return m, m.lanes[task.Status].list.InsertItem(len(m.lanes[task.Status].list.Items()), task)
	case EditTaskMsg:
		task := msg.task
		i := msg.index

		// Update in list
		return m, m.lanes[task.Status].list.SetItem(i, task)
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
