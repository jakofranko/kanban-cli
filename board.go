package main

import (
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	helpStyle = lipgloss.NewStyle().
			Foreground(grey)
	progressStyle = lipgloss.NewStyle().
			Margin(1)
)

// This is the model for the board view, which will implement the
// Bubbletea model methods for rendering etc. It maintains multiple
// lists components from bubbles, and is styled via lipgloss.
type Board struct {
	focused        status
	lanes          []SwimLane
	err            error
	loaded         bool
	quitting       bool
	project        int
	help           help.Model
	keys           boardKeyMap
	height         int
	width          int
	progress       progress.Model
	totalTasks     int
	completedTasks int
}

type UpdateListMsg struct {
	update         status
	newBoard       Board
	totalTasks     int
	completedTasks int
}

type ResetListHeightMsg struct{}

func resetListHeight() tea.Msg {
	return &ResetListHeightMsg{}
}

func NewBoard(project int, width int, height int) *Board {
	b := &Board{
		project:  project,
		keys:     boardKeys,
		help:     help.New(),
		progress: progress.New(progress.WithScaledGradient(secondary, highlight)),
		width:    width,
		height:   height,
		focused:  todo,
		loaded:   true,
	}

	b.initLists(width, height)

	// Focus the todo lane
	b.lanes[todo].Focus()

	return b
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

		// Before we change the status, handle changing the completed tasks
		if oldStatus == done {
			m.completedTasks--
		}

		updatedTask, err := taskDB.NextStatus(selectedTask)
		if err != nil {
			log.Fatal(err)
		}

		// Adjust completed tasks
		if updatedTask.Status == done {
			m.completedTasks++
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

		return UpdateListMsg{update: oldStatus, completedTasks: m.completedTasks}
	}

	return nil
}

func (m *Board) initLists(width, height int) {
	todoLane := new(SwimLane)
	inProgressLane := new(SwimLane)
	doneLane := new(SwimLane)

	m.lanes = []SwimLane{
		todoLane.Init(width, m.getListHeight(height), m.project, todo),
		inProgressLane.Init(width, m.getListHeight(height), m.project, inProgress),
		doneLane.Init(width, m.getListHeight(height), m.project, done),
	}

	// Count total and completed tasks for the progress bar.
	m.totalTasks = 0
	for _, lane := range m.lanes {
		m.totalTasks += len(lane.list.Items())

		if lane.laneStatus == done {
			m.completedTasks = len(lane.list.Items())
		}
	}
}

// This will return a height minus the height of other UI elements
func (m *Board) getListHeight(height int) int {
	// This can be expanded later if additional UI elements are
	// added to the Board view
	helpHeight := lipgloss.Height(m.help.View(boardKeys))
	progressHeight := lipgloss.Height(m.progress.ViewAs(0.0))
	progressMargin := 1
	return height - helpHeight - progressHeight - progressMargin*2
}

func (m Board) Init() tea.Cmd {
	// I don't understand what Init is for
	log.Print("initing board")
	return nil
}

func (m Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.initLists(msg.Width, msg.Height)
		m.lanes[m.focused].Focus()
		m.progress.Width = (msg.Width / 2) - horizontalPad*2
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
			models[form] = NewForm(m.width, m.height, m.focused, m.project)

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
			// Remove from DB
			taskDB := GetDB()
			defer taskDB.db.Close()

			currentList := m.lanes[m.focused]
			task := currentList.list.SelectedItem().(Task)

			err := taskDB.Delete(task.Id)
			if err != nil {
				log.Fatal(err)
			}

			m.totalTasks--
			if task.Status == done {
				m.completedTasks--
			}

			// And remove from UI
			i := m.lanes[m.focused].list.Index()
			m.lanes[m.focused].list.RemoveItem(i)
			return m, nil
		case key.Matches(msg, m.keys.Projects):
			// Back to the projects view
			return models[projects], m.RefreshProjects
		}
	case CreateTaskMsg:
		task := msg.task

		m.totalTasks++
		if task.Status == done {
			m.completedTasks++
		}

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
		m.completedTasks = msg.completedTasks

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

		return lipgloss.JoinVertical(
			lipgloss.Center,
			listsView,
			helpStyle.Render(m.help.View(m.keys)),
			progressStyle.Render(m.progress.ViewAs(float64(m.completedTasks)/float64(m.totalTasks))),
		)
	} else {
		return "loading..."
	}
}

func (b *Board) RefreshProjects() tea.Msg {
	return RefreshProjectsMsg{}
}
