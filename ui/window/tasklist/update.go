package tasklist

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haykh/gobrain/ui"
)

func (m *model) Update(msg tea.Msg) tea.Cmd {
	if m.active_action != NoAction {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, ui.Key_Cancel):
				return m.CancelAction()

			case key.Matches(msg, ui.Key_Accept):
				return m.AcceptAction()
			}
		}
	} else {
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch {

			case key.Matches(msg, keys.Up):
				return m.Navigate(ui.DirUp)

			case key.Matches(msg, keys.Down):
				return m.Navigate(ui.DirDown)

			case key.Matches(msg, keys.PgUp):
				return m.Navigate(ui.DirPgUp)

			case key.Matches(msg, keys.PgDown):
				return m.Navigate(ui.DirPgDown)

			case key.Matches(msg, keys.Toggle):
				return m.ToggleTask()

			case key.Matches(msg, keys.Edit):
				return m.EditItem()

			case key.Matches(msg, keys.Delete):
				return m.DeleteItem()

			case key.Matches(msg, keys.Add):
				return m.AddTask()

			case key.Matches(msg, keys.AddList):
				return m.AddTasklist()
			}
		}
	}
	return nil
}
