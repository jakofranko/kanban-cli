package main

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

type CreateTaskMsg struct {
    task Task
}

type EditTaskMsg struct {
    task Task
    index int
}

type DeleteTaskMsg struct {
    index int
}

func (t *Task) Next() {
    if t.status == done {
        t.status = todo
    } else {
        t.status++
    }
}

// Implement the bubbles/list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

func NewTask(status status, title string, description string) Task {
    return Task{status: status, title: title, description: description}
}

