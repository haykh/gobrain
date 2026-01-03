package tasklist

import (
	"sort"
	"time"

	"github.com/haykh/gobrain/backend"
)

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

func (m model) GetUrgentTasks(max_number int) []backend.TaskItem {
	due_tasks := map[time.Time]backend.TaskItem{}
	dates := []time.Time{}

	for _, tl := range m.tasklists {
		for _, t := range tl.Tasks {
			if !t.DueDate.IsZero() && !t.Checked {
				due_tasks[t.DueDate] = t.TaskItem
				dates = append(dates, t.DueDate)
			}
		}
	}

	sort.Sort(byDate(dates))
	urgent_tasks := []backend.TaskItem{}
	for i, d := range dates {
		if i >= max_number {
			break
		}
		urgent_tasks = append(urgent_tasks, due_tasks[d])
	}

	return urgent_tasks
}
