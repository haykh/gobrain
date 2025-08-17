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
	today          time.Time
	shown_date_min time.Time
	shown_date_max time.Time
	cursor         int

	calendar_days []CalendarDay

	app *backend.Backend
}

func New(app *backend.Backend) model {
	today := time.Now()
	shown_date_min := today.AddDate(0, 0, -int(today.Weekday())-13)
	shown_data_max := today.AddDate(0, 0, 7-int(today.Weekday())+15)
	return model{
		keymap:         keys,
		today:          today,
		shown_date_min: shown_date_min,
		shown_date_max: shown_data_max,
		cursor:         0,
		app:            app,
	}
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
