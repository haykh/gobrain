package randnotes

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/haykh/gobrain/components"
	"github.com/haykh/gobrain/ui"
	"github.com/haykh/gobrain/ui/window/randnotes/note"
)

func (m model) View() string {
	if len(m.filtered_notes) == 0 {
		return ""
	}
	n_views := []string{}
	for i := m.note_view_idx_min; i < m.note_view_idx_max; i++ {
		n_views = append(
			n_views,
			note.View(
				*m.filtered_notes[i],
				ui.Width_Window-1,
				m.cursor.NoteIndex == i,
				m.cursor.TagIndex,
				m.cursor.TagFilter,
			),
		)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			n_views...,
		),
		components.Scrollbar(
			ui.Height_Window,
			len(m.filtered_notes),
			m.note_view_idx_min,
			m.note_view_idx_max-m.note_view_idx_min,
		),
	)
}
