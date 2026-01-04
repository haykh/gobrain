package tasklist

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/haykh/gobrain/components"
	"github.com/haykh/gobrain/ui"
)

func (m model) View() string {
	if len(m.tasklists) == 0 {
		return ""
	}
	lines := []string{}
	counter := 0
	counter_visible := 0
	for tl_idx, tl := range m.tasklists {
		if m.view_range.IsVisible(counter) {
			lines = append(
				lines,
				tl.View(
					ui.Width_Window-1,
					m.cursor.TaskIndex == 0 && m.cursor.TasklistIndex == tl_idx,
				),
			)
			counter_visible++
		}
		counter++
		for t_idx, t := range tl.Tasks {
			if m.view_range.IsVisible(counter) {
				lines = append(
					lines,
					t.View(
						ui.Width_Window-1,
						m.cursor.TaskIndex == t_idx+1 && m.cursor.TasklistIndex == tl_idx,
					),
				)
				counter_visible++
			}
			counter++
		}
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			lines...,
		),
		components.Scrollbar(
			ui.Height_Window,
			counter,
			m.view_range.IMin,
			counter_visible,
		),
	)
}
