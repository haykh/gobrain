package task

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
)

type Task struct {
	backend.TaskItem
	is_editing bool
	input      *textinput.Model
	old_text   string
}

func (t Task) Type() string {
	return "task"
}

func New(text string, checked bool, importance int, dueDate time.Time) Task {
	return Task{
		TaskItem: backend.TaskItem{
			Text:       text,
			Checked:    checked,
			Importance: importance,
			DueDate:    dueDate,
		},
		is_editing: false,
		input:      nil,
		old_text:   "",
	}
}

func (t *Task) Toggle() {
	t.Checked = !t.Checked
}

func (t *Task) StartEditing(input *textinput.Model) {
	t.input = input
	t.input.Prompt = ""
	t.input.Placeholder = ""
	t.input.SetValue(t.Text)
	t.old_text = t.Text
	t.is_editing = true
}

func (t *Task) StopEditing(accept bool) {
	if accept {
		t.Text = t.input.Value()
	} else {
		t.Text = t.old_text
	}
	t.input.SetValue("")
	t.is_editing = false
	t.input = nil
}

func (t Task) View(width int, hover bool) string {
	line := lipgloss.NewStyle().Foreground(ui.Color_Fg_Braces).Render(" [")
	textstyle := lipgloss.NewStyle()
	if t.Checked {
		line += lipgloss.NewStyle().Foreground(ui.Color_Fg_Checkmark).Render("✓")
		textstyle = textstyle.Strikethrough(true).Foreground(ui.Color_Fg_Tasklist_Done)
	} else {
		line += " "
	}
	text := t.Text
	if t.is_editing && t.input != nil {
		text = t.input.View()
	}
	line += lipgloss.NewStyle().Foreground(ui.Color_Fg_Braces).Render("] ") + textstyle.Render(text)
	if hover && !t.is_editing {
		line += " <"
	}
	spacer := " "
	space := lipgloss.NewStyle().
		Foreground(ui.Color_Bg).Render(
		strings.Repeat(spacer, width-lipgloss.Width(line)-4),
	)
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		line,
		space,
	)
}
