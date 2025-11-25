package randnotes

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
)

type Cursor struct {
	NoteIndex int
	TagIndex  int
	TagFilter string
}

type model struct {
	keymap         help.KeyMap
	notes          []backend.Note
	filtered_notes []*backend.Note
	cursor         Cursor

	note_view_idx_min int
	note_view_idx_max int

	app *backend.Backend
}

func New(app *backend.Backend) model {
	m := model{
		keymap:         keys,
		notes:          []backend.Note{},
		filtered_notes: []*backend.Note{},
		cursor:         Cursor{0, -1, ""},

		note_view_idx_min: 0,
		note_view_idx_max: 0,

		app: app,
	}
	m.Sync()
	return m
}

func (m model) Path() []string {
	return []string{"dashboard", "random notes"}
}

func (m model) Parent() ui.PanelType {
	return ui.PanelDashboard
}

func (m model) Keys() help.KeyMap {
	return m.keymap
}
