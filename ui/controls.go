package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Direction int

const (
	DirUp Direction = iota
	DirDown
	DirLeft
	DirRight
)

type PanelType int

const (
	PanelDashboard PanelType = iota
	PanelCalendar
	PanelTasklist
	PanelRandomNotes
	PanelMdViewer
)

type PanelView interface {
	Path() []string
	Parent() PanelType

	View() string
	Sync()
	Update(msg tea.Msg) tea.Cmd

	Keys() help.KeyMap
}

// Messages/signals
type ErrorMsg struct {
	Error error
}

type NavigateFwdMsg struct {
	NewPanel PanelType
}

type NavigateBackMsg struct{}

type NewViewer struct {
	Filepath string
	Filename string
}

type OpenEditorMsg struct {
	Filename string
}

type TrashRandomNoteMsg struct {
	Filename string
}

// Key bindings
var Key_Up = key.NewBinding(
	key.WithKeys("up", "k"),
	key.WithHelp("↑/k", "up"),
)
var Key_Down = key.NewBinding(
	key.WithKeys("down", "j"),
	key.WithHelp("↓/j", "down"),
)
var Key_Left = key.NewBinding(
	key.WithKeys("left", "h"),
	key.WithHelp("←/h", "left"),
)
var Key_Right = key.NewBinding(
	key.WithKeys("right", "l"),
	key.WithHelp("→/l", "right"),
)
var Key_Select = key.NewBinding(
	key.WithKeys("enter", " "),
	key.WithHelp("enter/space", "select"),
)
var Key_Edit = key.NewBinding(
	key.WithKeys("e"),
	key.WithHelp("e", "edit"),
)
var Key_Delete = key.NewBinding(
	key.WithKeys("delete", "d"),
	key.WithHelp("del/d", "delete"),
)
var Key_Filter = key.NewBinding(
	key.WithKeys("f"),
	key.WithHelp("f", "filter"),
)
var Key_Back = key.NewBinding(
	key.WithKeys("backspace"),
	key.WithHelp("backspace", "back"),
)
var Key_Help = key.NewBinding(
	key.WithKeys("?"),
	key.WithHelp("?", "toggle help"),
)
var Key_Quit = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("q", "quit"),
)
