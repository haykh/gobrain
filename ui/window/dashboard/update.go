package dashboard

import (
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haykh/gobrain/ui"
)

func navigate(current PanelCategory, dir ui.Direction) PanelCategory {
	switch current {
	case DailyNotes:
		switch dir {
		case ui.DirDown:
			return TodaysNote
		case ui.DirRight:
			return TaskLists
		default:
			return current
		}
	case TaskLists:
		switch dir {
		case ui.DirLeft:
			return DailyNotes
		case ui.DirRight:
			return RandomNotes
		default:
			return current
		}
	case TodaysNote:
		switch dir {
		case ui.DirUp:
			return DailyNotes
		case ui.DirRight:
			return TaskLists
		default:
			return current
		}
	case RandomNotes:
		switch dir {
		case ui.DirLeft:
			return TaskLists
		case ui.DirDown:
			return NewRandomNote
		default:
			return current
		}
	case NewRandomNote:
		switch dir {
		case ui.DirUp:
			return RandomNotes
		case ui.DirLeft:
			return TaskLists
		default:
			return current
		}
	default:
		panic("unknown note category")
	}
}

func (m *model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keys.Up):
			m.cursor = navigate(m.cursor, ui.DirUp)
			return nil

		case key.Matches(msg, keys.Down):
			m.cursor = navigate(m.cursor, ui.DirDown)
			return nil

		case key.Matches(msg, keys.Left):
			m.cursor = navigate(m.cursor, ui.DirLeft)
			return nil

		case key.Matches(msg, keys.Right):
			m.cursor = navigate(m.cursor, ui.DirRight)
			return nil

		case key.Matches(msg, keys.Select):
			switch m.cursor {

			case RandomNotes:
				return func() tea.Msg {
					return ui.NavigateFwdMsg{NewPanel: ui.PanelRandomNotes}
				}

			case NewRandomNote:
				if newfile, err := m.app.CreateNew_RandomNote(); err == nil {
					return func() tea.Msg {
						return ui.OpenEditorMsg{Filename: filepath.Join(m.app.RandomNotesPath, newfile)}
					}
				} else {
					return func() tea.Msg {
						return ui.ErrorMsg{Error: err}
					}
				}
			}
		}
	}

	return nil
}

func (m *model) Sync() {}
