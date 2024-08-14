package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type boardKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Help     key.Binding
	Edit     key.Binding
	New      key.Binding
	Move     key.Binding
	Delete   key.Binding
	Projects key.Binding
	Quit     key.Binding
}

type formKeyMap struct {
	Next key.Binding
	Quit key.Binding
}

type projectListKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	New      key.Binding
	Archive  key.Binding
	MoveUp   key.Binding
	MoveDown key.Binding
	Quit     key.Binding
	Help     key.Binding
	Select   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k boardKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.New, k.Projects, k.Quit, k.Help}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k boardKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.New, k.Edit, k.Delete},       // second column
		{k.Projects, k.Quit, k.Help},    // third column
	}
}

func (k formKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k formKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Quit},
	}
}

func (k projectListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Help}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k projectListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.MoveUp, k.MoveDown},
		{k.Select, k.New, k.Archive},
		{k.Help, k.Quit},
	}
}

var boardKeys = boardKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit task"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new task"),
	),
	Move: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "move task"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete task"),
	),
	Projects: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "projects"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

var formKeys = formKeyMap{
	Next: key.NewBinding(
		key.WithKeys("ctrl+y"),
		key.WithHelp("ctrl+y", "next field/confirm"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

var projectListKeys = projectListKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	MoveUp: key.NewBinding(
		key.WithKeys("ctrl+up", "ctrl+k"),
		key.WithHelp("ctrl+ ↑/k", "move project up"),
	),
	MoveDown: key.NewBinding(
		key.WithKeys("ctrl+down", "ctrl+j"),
		key.WithHelp("ctrl+ ↓/j", "move project down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new project"),
	),
	Archive: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "archive project"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", "space"),
		key.WithHelp("enter/space", "select project"),
	),
}
