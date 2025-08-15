package window

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
	"github.com/haykh/gobrain/ui/window/dashboard"
	"github.com/haykh/gobrain/ui/window/randnotes"
)

type Window struct {
	debug          bool
	debug_logs     []string
	max_debug_logs int

	weather              string
	weather_last_updated time.Time

	show_help bool
	help      help.Model

	active_panel ui.PanelType
	panels       map[ui.PanelType]ui.PanelView

	mdviewport_filename string
	mdviewport          viewport.Model
	mdviewport_show     bool

	// backend
	app *backend.Backend
}

func New(app *backend.Backend, show_help bool, debug bool) Window {
	dashboard_panel := dashboard.New(app)
	randnotes_panel := randnotes.New(app)
	// randnotes_mdviewer_panel := mdviewer.New(app, &randnotes_panel, ui.PanelRandomNotes)
	panels := map[ui.PanelType]ui.PanelView{
		ui.PanelDashboard:   &dashboard_panel,
		ui.PanelRandomNotes: &randnotes_panel,
	}
	window := Window{
		debug:          debug,
		debug_logs:     []string{},
		max_debug_logs: 30,

		weather:              "",
		weather_last_updated: time.Time{},

		show_help: show_help,
		help:      help.New(),

		active_panel: ui.PanelDashboard,
		panels:       panels,

		mdviewport_filename: "",
		mdviewport:          viewport.New(0, 0),
		mdviewport_show:     false,

		app: app,
	}
	window.FetchWeather()
	return window
}

func (w Window) Init() tea.Cmd {
	return nil
}

func (w Window) Active() ui.PanelView {
	return w.panels[w.active_panel]
}

func (w *Window) DebugLog(log string) {
	if w.debug {
		if len(w.debug_logs) >= w.max_debug_logs {
			w.debug_logs = w.debug_logs[1:]
		}
		w.debug_logs = append(w.debug_logs, log)
	}
}
