package calendar

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Sync() {
	m.calendar_days = []CalendarDay{}
	for d := m.shown_date_min; d.Before(m.shown_date_max); d = d.AddDate(0, 0, 1) {
		m.calendar_days = append(m.calendar_days, CalendarDay{
			Date:       d,
			NoteExists: false,
		})
		if isToday(d) {
			m.cursor = len(m.calendar_days) - 1
		}
	}
}

func (m *model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keys.Up):
			if m.cursor > 7 {
				m.cursor -= 7
			} else {
				m.cursor = 0
			}
		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.calendar_days)-8 {
				m.cursor += 7
			} else {
				m.cursor = len(m.calendar_days) - 1
			}
		case key.Matches(msg, keys.Left):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, keys.Right):
			if m.cursor < len(m.calendar_days)-1 {
				m.cursor++
			}

		}
	}
	return nil
}
