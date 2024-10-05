# Kanban CLI

A kanban board written in Go + [Bubble Tea](https://github.com/charmbracelet/bubbletea)! This project was heavily inspired by a couple of the [charm_](https://charm.sh/) team's tutorials for using their toolchains, namely the [kancli](https://github.com/charmbracelet/kancli) and [taskcli](https://github.com/charmbracelet/taskcli) example projects. This was kind of a mish-mash of the two concepts in a way that was the most useful for me. This project should properly be called a TUI at the moment, as it does not currently support CLI arguments or STD I/O, but that may be a venue for future improvement.

## Installation

`go install github.com/jakofranko/kanban-cli`

to install the binary directly via Go.

## Usage

Installing the binary will give you access to the `kanban-cli` command. Run this to start up your fresh kanban project tool! Help notes should be at the bottom of every view if you ever get lost.

### Projects

The initial view of the TUI is a projects view. It will initially be empty. Note the help view at the bottom. Press '?' for more options. 'n' will create a new project. Enter a name for your project and press 'enter'. You should see a new empty project in your project list!

Use the arrow or vim keys to navigate.

Press 'enter' on a highlighted project to view that project's kanban board.

When you have completed a project, you can press 'a' to archive the project and move it out of your main view. Press 'v' to toggle between active and archived projects.

### Kanban Board

Every project has its own board of tasks, with three "swim lanes": todo, in progress, and done. Create a new task with 'n'. Enter a name for the task and press 'ctrl+y' to confirm and then enter a description. Press 'ctrl+y' to create the new task.

Use the arrow or vim keys to navigate between tasks and swim lanes.

'Enter' will move a task to the next status. Don't worry if you accidentally move a task to done or in progress, you can cycle tasks in the done column back to the todo column.

Don't forget to press the '?' key to view all the options you have! You can delete tasks, edit tasks, and view tasks so that you can read all the details you put in the description.
