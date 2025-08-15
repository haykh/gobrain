package dashboard

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
)

type PanelCategory int

const (
	DailyNotes PanelCategory = iota
	TaskLists
	RandomNotes
	TodaysNote
	NewRandomNote
)

type model struct {
	cursor PanelCategory
	keymap help.KeyMap

	app *backend.Backend
}

func New(app *backend.Backend) model {
	return model{
		cursor: DailyNotes,
		keymap: keys,
		app:    app,
	}
}

func (m model) Path() []string {
	return []string{"dashboard"}
}

func (m model) Parent() ui.PanelType {
	return ui.PanelDashboard
}

func (m model) Keys() help.KeyMap {
	return m.keymap
}
