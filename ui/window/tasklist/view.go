package tasklist

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/haykh/gobrain/components"
	"github.com/haykh/gobrain/ui"
)

func (m model) View() string {
	if len(m.filtered_tasks) == 0 {
		return ""
	}
	l_views := []string{}
	for i := m.task_view_idx_min; i < m.task_view_idx_max; i++ {
		l_views = append(
			l_views,
			m.filtered_tasks[i].View(
				ui.Width_Window-1,
				(m.cursor.Index == i) && !m.is_adding_task && !m.is_adding_list,
			),
		)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			l_views...,
		),
		components.Scrollbar(
			ui.Height_Window,
			len(m.filtered_tasks),
			m.task_view_idx_min,
			m.task_view_idx_max-m.task_view_idx_min,
		),
	)
}
