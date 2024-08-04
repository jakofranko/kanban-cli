package main

import (
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
	help     help.Model
	keys     boardKeyMap
	height   int
	width    int
}

type UpdateListMsg struct {
	update   status
	newBoard Board
}

type ResetListHeightMsg struct{}

func resetListHeight() tea.Msg {
	return &ResetListHeightMsg{}
}

func NewBoard() *Board {
	return &Board{
		keys: boardKeys,
		help: help.New(),
	}
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

	m.lanes = []SwimLane{
		todoLane.Init(width, m.getListHeight(height), todo),
		inProgressLane.Init(width, m.getListHeight(height), inProgress),
		doneLane.Init(width, m.getListHeight(height), done),
	}
}

// This will return a height minus the height of other UI elements
func (m *Board) getListHeight(height int) int {
	// This can be expanded later if additional UI elements are
	// added to the Board view
	return height - lipgloss.Height(m.help.View(boardKeys))
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
			m.width = msg.Width
			m.height = msg.Height
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
			m.focused = todo
			m.lanes[m.focused].Focus()
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Left):
			l := m.lanes[m.focused].list
			p := l.Paginator
			if p.TotalPages > 1 && !p.OnFirstPage() {
				m.lanes[m.focused].list.Paginator.PrevPage()
			} else {
				m.Prev()
			}
		case key.Matches(msg, m.keys.Right):
			l := m.lanes[m.focused].list
			p := l.Paginator
			if p.TotalPages > 1 && !p.OnLastPage() {
				m.lanes[m.focused].list.Paginator.NextPage()
			} else {
				m.Next()
			}
		case key.Matches(msg, m.keys.Move):
			return m, m.MoveToNext
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			model, _ := m.help.Update(nil)
			m.help = model

			// Update lanes
			var cmds []tea.Cmd
			for i, lane := range m.lanes {
				lane.SetHeight(m.getListHeight(m.height))
				_, cmd := lane.list.Update(nil)
				m.lanes[i] = lane
				cmds = append(cmds, cmd)
			}

			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keys.New):
			models[board] = m // save current model
			models[form] = NewForm(m.focused, m.project)

			return models[form], nil
		case key.Matches(msg, m.keys.Edit):
			models[board] = m // save current model
			currentTask := m.lanes[m.focused].list.SelectedItem().(Task)
			currentIndex := m.lanes[m.focused].list.Index()
			models[form] = UpdateForm(currentTask, currentIndex)

			return models[form], nil
		case key.Matches(msg, m.keys.Delete):
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

		listsView := lipgloss.JoinHorizontal(
			lipgloss.Left,
			todoView,
			ipView,
			doneView,
		)

		return lipgloss.JoinVertical(lipgloss.Center, listsView, m.help.View(m.keys))
	} else {
		return "loading..."
	}
}
