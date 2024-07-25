package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)


// Models
var models []tea.Model
const (
    board status = iota
    form
)

func main() {
    f, err := tea.LogToFile("debug.log", "debug")
    if err != nil {
        fmt.Println("fatal:", err)
        os.Exit(1)
    }
    defer f.Close()

    // NewBoard is defined in board.go
    // NewForm is defined in form.go
    models = []tea.Model{NewBoard(), NewForm(todo)}
	m := models[board]
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
