package tasklist

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
	"github.com/haykh/gobrain/ui/window/tasklist/list"
)

type Action int

type ItemIndex struct {
	ListIndex int
	TaskIndex int
}

type ViewRange struct {
	IMin      int
	MaxHeight int
}

const (
	NoAction Action = iota
	AddingTask
	AddingList
	EditingTask
	EditingList
)

type Cursor struct {
	TasklistIndex int
	TaskIndex     int
}

type model struct {
	keymap        help.KeyMap
	tasklists     []list.List
	active_action Action
	active_item   ItemIndex

	cursor     Cursor
	view_range ViewRange

	app *backend.Backend
}

func New(app *backend.Backend) model {
	m := model{
		keymap:    keys,
		tasklists: []list.List{},

		active_action: NoAction,
		active_item:   ItemIndex{-1, -1},

		cursor:     Cursor{0, 0},
		view_range: ViewRange{0, ui.Height_Window},
		app:        app,
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

/**
 * Cursor helper
**/

func (c Cursor) ViewIndex(tasklists []list.List) int {
	index := 0
	for i := 0; i < c.TasklistIndex; i++ {
		index += 1 + tasklists[i].NumTasks()
	}
	index += c.TaskIndex
	return index
}

/**
 * ViewRange helper
**/

func (v ViewRange) NumItems(tasklists []list.List) int {
	count := 0
	for _, tl := range tasklists {
		count += tl.NumTasks() + 1
	}
	return count
}

func (v ViewRange) IsVisible(index int) bool {
	return index >= v.IMin && index < v.IMin+v.MaxHeight
}

func (v *ViewRange) RemoveItems(n int, tasklists []list.List) {
	num_items := v.NumItems(tasklists)
	if num_items < v.MaxHeight {
		v.IMin = 0
	} else {
		v.IMin = max(0, v.IMin-n)
	}
}

func (v *ViewRange) EnsureCursorVisible(cursorIndex int) {
	if cursorIndex < v.IMin {
		v.IMin = cursorIndex
	} else if cursorIndex >= v.IMin+v.MaxHeight {
		v.IMin = cursorIndex - v.MaxHeight + 1
	}
}
