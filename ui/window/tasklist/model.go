package tasklist

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
	"github.com/haykh/gobrain/ui/window/tasklist/list"
)

type Selector interface {
	View(int, bool) string
	Type() string
	Index() int
}

type Cursor struct {
	Index int
}

type model struct {
	keymap          help.KeyMap
	tasklists       []list.List
	filtered_tasks  []Selector
	cursor          Cursor
	is_adding_list  bool
	is_adding_task  bool
	adding_list_idx int
	adding_task_idx int

	task_view_idx_min int
	task_view_idx_max int

	app *backend.Backend
}

func New(app *backend.Backend) model {
	m := model{
		keymap:          keys,
		tasklists:       []list.List{},
		filtered_tasks:  []Selector{},
		cursor:          Cursor{0},
		is_adding_list:  false,
		is_adding_task:  false,
		adding_list_idx: -1,
		adding_task_idx: -1,

		task_view_idx_min: 0,
		task_view_idx_max: 0,

		app: app,
	}
	m.Sync()
	return m
}

func (m model) Path() []string {
	return []string{"dashboard", "task lists"}
}

func (m model) Parent() ui.PanelType {
	return ui.PanelDashboard
}

func (m model) Keys() help.KeyMap {
	return m.keymap
}
