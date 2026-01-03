package window

import (
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
	"github.com/haykh/gobrain/ui/window/calendar"
	"github.com/haykh/gobrain/ui/window/dashboard"
	"github.com/haykh/gobrain/ui/window/randnotes"
	"github.com/haykh/gobrain/ui/window/tasklist"
)

type Window struct {
	debug          bool
	debug_logs     []string
	max_debug_logs int

	weather              string
	weather_last_updated time.Time
	weather_fetching     bool
	httpClient           *http.Client

	typing_input textinput.Model
	is_typing    bool

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
	calendar_panel := calendar.New(app)
	tasklists_panel := tasklist.New(app)
	randnotes_panel := randnotes.New(app)
	panels := map[ui.PanelType]ui.PanelView{
		ui.PanelDashboard:   &dashboard_panel,
		ui.PanelTaskLists:   &tasklists_panel,
		ui.PanelCalendar:    &calendar_panel,
		ui.PanelRandomNotes: &randnotes_panel,
	}
	window := Window{
		debug:          debug,
		debug_logs:     []string{},
		max_debug_logs: 20,

		weather:              "fetching weather...",
		weather_last_updated: time.Time{},
		weather_fetching:     true,
		httpClient:           &http.Client{Timeout: 30 * time.Second},

		is_typing: false,

		show_help: show_help,
		help:      help.New(),

		active_panel: ui.PanelDashboard,
		panels:       panels,

		mdviewport_filename: "",
		mdviewport:          viewport.New(0, 0),
		mdviewport_show:     false,

		app: app,
	}
	return window
}

func (w Window) Init() tea.Cmd {
	if w.weather_fetching {
		return fetchWeatherCmd(w.httpClient)
	}
	return textinput.Blink
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
