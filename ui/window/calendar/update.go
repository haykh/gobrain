package calendar

import (
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/haykh/gobrain/ui"
)

func (m *model) Reset() {
	today := time.Now()
	m.shown_date_min = today.AddDate(0, 0, -int(today.Weekday())-20)
	m.shown_date_max = today.AddDate(0, 0, 7-int(today.Weekday())+8)
	m.Sync()
	m.ResetCursor()
}

func (m *model) ResetCursor() {
	di := 0
	for d := m.shown_date_min; d.Before(m.shown_date_max); d = d.AddDate(0, 0, 1) {
		if isToday(d) {
			m.cursor = di
		}
		di++
	}
}

func (m *model) Sync() error {
	m.calendar_days = []CalendarDay{}
	filenames, err := m.app.GetMarkdownFilenames(m.app.DailyNotes)
	if err != nil {
		return err
	}
	for d := m.shown_date_min; d.Before(m.shown_date_max); d = d.AddDate(0, 0, 1) {
		note_exists := false
		for _, filename := range filenames {
			if d.Format("2006-Jan-02")+".md" == filename {
				note_exists = true
			}
		}
		m.calendar_days = append(m.calendar_days, CalendarDay{
			Date:       d,
			NoteExists: note_exists,
		})
	}
	return nil
}

func (m *model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keys.Reset):
			m.Reset()
		case key.Matches(msg, keys.Up):
			if m.cursor > 7 {
				m.cursor -= 7
			} else {
				m.shown_date_min = m.shown_date_min.AddDate(0, 0, -7)
				m.shown_date_max = m.shown_date_max.AddDate(0, 0, -7)
				m.Sync()
			}
			return nil

		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.calendar_days)-7 {
				m.cursor += 7
			} else {
				m.shown_date_min = m.shown_date_min.AddDate(0, 0, 7)
				m.shown_date_max = m.shown_date_max.AddDate(0, 0, 7)
				m.Sync()
			}
			return nil

		case key.Matches(msg, keys.Left):
			if m.cursor > 0 {
				m.cursor--
			}
			return nil

		case key.Matches(msg, keys.Right):
			if m.cursor < len(m.calendar_days)-1 {
				m.cursor++
			}
			return nil

		case key.Matches(msg, keys.Select):
			if m.calendar_days[m.cursor].NoteExists {
				return func() tea.Msg {
					return ui.NewViewer{
						Filepath: m.app.DailyNotes,
						Filename: m.calendar_days[m.cursor].Date.Format("2006-Jan-02") + ".md",
					}
				}
			} else {
				return nil
			}

		case key.Matches(msg, keys.Edit):
			fname := m.calendar_days[m.cursor].Date.Format("2006-Jan-02") + ".md"
			if !m.calendar_days[m.cursor].NoteExists {
				filename, err := m.app.CreateNew_DailyNote(m.calendar_days[m.cursor].Date)
				if err != nil {
					return func() tea.Msg {
						return ui.ErrorMsg{
							Error: err,
						}
					}
				}
				fname = filename
			}
			return func() tea.Msg {
				return ui.OpenEditorMsg{
					Filename: filepath.Join(m.app.DailyNotes, fname),
				}
			}

		}
	}
	return nil
}
