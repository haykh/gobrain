package backend

import (
	"sort"
	"time"
)

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

/**
 * Helper function for getting the most urgent tasks across multiple task lists
 */
type byDate []time.Time

func (a byDate) Len() int {
	return len(a)
}

func (a byDate) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byDate) Less(i, j int) bool {
	return a[i].Before(a[j])
}

func GetUrgentTasks(tasklists []TaskList, max_number int) []TaskItem {
	due_tasks := map[time.Time]TaskItem{}
	dates := []time.Time{}

	for _, tl := range tasklists {
		for _, t := range tl.Items {
			if !t.DueDate.IsZero() && !t.Checked {
				due_tasks[t.DueDate] = t
				dates = append(dates, t.DueDate)
			}
		}
	}

	sort.Sort(byDate(dates))
	urgent_tasks := []TaskItem{}
	for i, d := range dates {
		if i >= max_number {
			break
		}
		urgent_tasks = append(urgent_tasks, due_tasks[d])
	}

	return urgent_tasks
}
