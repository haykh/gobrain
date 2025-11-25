package backend

import "time"

// Note represents a markdown-backed note with metadata. Using a shared struct
// makes it straightforward to extend support for libraries and synced sources.
type Note struct {
	Filename string
	Title    string
	Icon     string
	Tags     []string
	Created  time.Time
	Path     string
}

// TaskItem represents a single checkbox item in a task list.
type TaskItem struct {
	Text       string
	Checked    bool
	Importance int
	DueDate    time.Time
}

// TaskList holds multiple TaskItems alongside a title and location.
type TaskList struct {
	Title    string
	Filename string
	Items    []TaskItem
	Path     string
}
