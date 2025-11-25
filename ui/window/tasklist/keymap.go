package tasklist

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/haykh/gobrain/ui"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	PgUp   key.Binding
	PgDown key.Binding
	Add    key.Binding
	Toggle key.Binding
	Edit   key.Binding
	Delete key.Binding
	Back   key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PgUp, k.PgDown},
		{k.Add, k.Toggle, k.Edit, k.Delete},
		{k.Back, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up:     ui.Key_Up,
	Down:   ui.Key_Down,
	PgUp:   ui.Key_PgUp,
	PgDown: ui.Key_PgDown,
	Add:    ui.Key_Add,
	Toggle: ui.Key_Toggle,
	Edit:   ui.Key_Edit,
	Delete: ui.Key_Delete,
	Back:   ui.Key_Back,
	Help:   ui.Key_Help,
	Quit:   ui.Key_Quit,
}
