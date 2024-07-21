
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
