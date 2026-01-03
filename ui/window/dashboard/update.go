package dashboard

import (
	"path/filepath"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/haykh/gobrain/backend"
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

		case key.Matches(msg, keys.Edit):
			switch m.cursor {

			case TodaysNote:
				today_str := time.Now().Format("2006-Jan-02") + ".md"
				daily_notes, err := m.app.GetMarkdownFilenames(m.app.DailyNotes)
				if err != nil {
					return func() tea.Msg {
						return ui.ErrorMsg{Error: err}
					}
				}
				if !slices.Contains(daily_notes, today_str) {
					fname, err := m.app.CreateNew_DailyNote(time.Now())
					if err != nil {
						return func() tea.Msg {
							return ui.ErrorMsg{Error: err}
						}
					}
					today_str = fname
				}
				return func() tea.Msg {
					return ui.OpenEditorMsg{
						Filename: filepath.Join(m.app.DailyNotes, today_str),
					}
				}

			}

		case key.Matches(msg, keys.Select):
			switch m.cursor {

			case DailyNotes:
				return func() tea.Msg {
					return ui.NavigateFwdMsg{NewPanel: ui.PanelCalendar}
				}

			case TodaysNote:
				today_str := time.Now().Format("2006-Jan-02") + ".md"
				daily_notes, err := m.app.GetMarkdownFilenames(m.app.DailyNotes)
				if err != nil {
					return func() tea.Msg {
						return ui.ErrorMsg{Error: err}
					}
				}
				if slices.Contains(daily_notes, today_str) {
					return func() tea.Msg {
						return ui.NewViewer{
							Filepath: m.app.DailyNotes,
							Filename: today_str,
						}
					}
				}

			case TaskLists:
				return func() tea.Msg {
					return ui.NavigateFwdMsg{NewPanel: ui.PanelTaskLists}
				}

			case RandomNotes:
				return func() tea.Msg {
					return ui.NavigateFwdMsg{NewPanel: ui.PanelRandomNotes}
				}

			case NewRandomNote:
				if newfile, err := m.app.CreateNew_RandomNote(); err == nil {
					return func() tea.Msg {
						return ui.OpenEditorMsg{Filename: filepath.Join(m.app.RandomNotes, newfile)}
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

func (m *model) Sync() error {
	filenamesTasklists, err := m.app.GetMarkdownFilenames(m.app.Tasks)
	if err != nil {
		return err
	}

	tasklists := []backend.TaskList{}
	for _, filename := range filenamesTasklists {
		if title, tasks, checked, importances, dueDates, err := backend.ParseMarkdownTasklist(m.app.Tasks, filename); err != nil {
			return err
		} else {
			tasklist := backend.TaskList{Title: title, Filename: filename, Path: m.app.Tasks}
			for i := range tasks {
				tasklist.Items = append(tasklist.Items, backend.TaskItem{
					Text:       tasks[i],
					Checked:    checked[i],
					Importance: importances[i],
					DueDate:    dueDates[i],
				})
			}
			tasklists = append(tasklists, tasklist)
		}
	}

	m.urgentTasks = backend.GetUrgentTasks(tasklists, 8)
	return nil
}
