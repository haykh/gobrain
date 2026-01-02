package list

import (
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
	"github.com/haykh/gobrain/ui/window/tasklist/task"
)

type List struct {
	backend.TaskList
	index int
	Tasks []task.Task

	is_editing bool
	input      *textinput.Model
	old_title  string
}

func (l List) Type() string {
	return "tasklist"
}

func New(title, filename, path string, index int) List {
	return List{
		TaskList:   backend.TaskList{Title: title, Filename: filename, Path: path},
		index:      index,
		Tasks:      []task.Task{},
		is_editing: false,
		input:      nil,
		old_title:  "",
	}
}

func (l *List) StartEditing(input *textinput.Model, placeholder string) {
	l.input = input
	l.input.Prompt = ""
	l.input.Placeholder = placeholder
	l.input.SetValue(l.Title)
	l.old_title = l.Title
	l.is_editing = true
}

func (l *List) StopEditing(accept bool) {
	if l.input == nil {
		return
	}

	if accept {
		l.Title = l.input.Value()
	} else {
		l.Title = l.old_title
	}

	l.input.SetValue("")
	l.input = nil
	l.is_editing = false
}

func (l *List) AppendTask(t task.Task) {
	l.Tasks = append(l.Tasks, t)
}

func (l *List) AddRawTask(text string, i int) {
	pattern := regexp.MustCompile(`^(.*?)( \{(\!*)\})?( \{([0-9]{4}-[0-9]{2}-[0-9]{2})?\})?$`)
	matches := pattern.FindStringSubmatch(text)
	var (
		taskText   string
		importance int
		dueTime    time.Time
	)
	if len(matches) > 2 {
		if matches[3] != "" {
			importance = len(matches[3])
		}
	}
	if len(matches) > 4 {
		if matches[5] != "" {
			dueTime, _ = time.Parse("2006-01-02", matches[5])
		}
	}
	if len(matches) > 1 {
		taskText = matches[1]
	} else {
		taskText = text
	}
	newtask := task.New(taskText, false, importance, dueTime, i, l.index)
	l.Tasks = append(l.Tasks[:i], append([]task.Task{newtask}, l.Tasks[i:]...)...)
}

func (l *List) AddTask(t task.Task, i int) {
	l.Tasks = append(l.Tasks[:i], append([]task.Task{t}, l.Tasks[i:]...)...)
}

func (l *List) Delete(app *backend.Backend) error {
	return app.TrashTasklist(l.Filename, l.Path)
}

func (l *List) RemoveTask(index int) {
	if index < 0 || index >= len(l.Tasks) {
		return
	}
	l.Tasks = append(l.Tasks[:index], l.Tasks[index+1:]...)
}

func (l List) NumTasks() int {
	return len(l.Tasks)
}

func (l List) View(width int, hover bool) string {
	titleText := l.Title
	if l.is_editing && l.input != nil {
		titleText = l.input.View()
	}

	title := lipgloss.NewStyle().Underline(true).Render(titleText)
	if hover {
		title += " <"
	}
	spacer := " "
	space := lipgloss.NewStyle().
		Foreground(ui.Color_Bg).Render(
		strings.Repeat(spacer, width-lipgloss.Width(title)),
	)
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		space,
	)
}

func (l List) Sync() error {
	tasks := []string{}
	checked := []bool{}
	importances := []int{}
	dueDates := []time.Time{}

	for _, t := range l.Tasks {
		tasks = append(tasks, t.Text)
		checked = append(checked, t.Checked)
		importances = append(importances, t.Importance)
		if !t.DueDate.IsZero() {
			dueDates = append(dueDates, t.DueDate)
		} else {
			dueDates = append(dueDates, time.Time{})
		}
	}
	return backend.WriteMarkdownTasklist(l.Path, l.Filename, l.Title, tasks, checked, importances, dueDates)
}

func (l List) Index() int {
	return l.index
}
