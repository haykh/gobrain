package calendar

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
)

type CalendarDay struct {
	Date       time.Time
	NoteExists bool
}

type model struct {
	keymap         help.KeyMap
	shown_date_min time.Time
	shown_date_max time.Time
	cursor         int

	calendar_days []CalendarDay

	app *backend.Backend
}

func New(app *backend.Backend) model {
	m := model{
		keymap:         keys,
		shown_date_min: time.Time{},
		shown_date_max: time.Time{},
		cursor:         0,
		app:            app,
	}
	m.Reset()
	return m
}

func (m model) Path() []string {
	return []string{"dashboard", "calendar"}
}

func (m model) Parent() ui.PanelType {
	return ui.PanelDashboard
}

func (m model) Keys() help.KeyMap {
	return m.keymap
}
