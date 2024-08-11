package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Models
var models []tea.Model

const (
	board status = iota
	form
	projects
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
	log.Println("Starting Cli...")

	// TODO confirm that the new form project here doesn't matter?
	models = []tea.Model{NewBoard("this doesn't exist", 0, 0), NewForm(todo, "not real"), NewProjectsTable()}
	m := models[projects]
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
