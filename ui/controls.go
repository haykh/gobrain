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
	DirPgUp
	DirPgDown
)

type PanelType int

const (
	PanelDashboard PanelType = iota
	PanelCalendar
	PanelTaskLists
	PanelRandomNotes
	PanelMdViewer
)

type PanelView interface {
	Path() []string
	Parent() PanelType

	View() string
	Sync() error
	Update(msg tea.Msg) tea.Cmd

	Keys() help.KeyMap
}

// Messages/signals
type ErrorMsg struct {
	Error error
}

type DebugMsg struct {
	Message string
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

type TrashNoteMsg struct {
	Filepath string
	Filename string
}

type TypingStartMsg struct{}
type TypingEndMsg struct{}

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
	key.WithHelp("↵/␣", "select"),
)
var Key_Add = key.NewBinding(
	key.WithKeys("a"),
	key.WithHelp("a", "add/create"),
)
var Key_AddList = key.NewBinding(
	key.WithKeys("shift+a", "A"),
	key.WithHelp("A", "add tasklist"),
)
var Key_Edit = key.NewBinding(
	key.WithKeys("e"),
	key.WithHelp("e", "edit"),
)
var Key_Delete = key.NewBinding(
	key.WithKeys("delete", "x"),
	key.WithHelp("del/x", "delete"),
)
var Key_Reset = key.NewBinding(
	key.WithKeys("r"),
	key.WithHelp("r", "reset"),
)
var Key_Filter = key.NewBinding(
	key.WithKeys("f"),
	key.WithHelp("f", "filter"),
)
var Key_Toggle = key.NewBinding(
	key.WithKeys("t"),
	key.WithHelp("t", "toggle"),
)
var Key_PgUp = key.NewBinding(
	key.WithKeys("pgup", "u"),
	key.WithHelp("pgup/u", "page up"),
)
var Key_PgDown = key.NewBinding(
	key.WithKeys("pgdown", "d"),
	key.WithHelp("pgdn/d", "page down"),
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

// keys for typing input
var Key_Cancel = key.NewBinding(
	key.WithKeys("esc", "ctrl+c"),
	key.WithHelp("esc/ctrl+c", "cancel"),
)
var Key_Accept = key.NewBinding(
	key.WithKeys("enter"),
	key.WithHelp("↵", "accept"),
)
