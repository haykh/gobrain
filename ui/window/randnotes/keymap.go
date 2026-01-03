package randnotes

import (
	"github.com/charmbracelet/bubbles/key"

	"github.com/haykh/gobrain/ui"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Add    key.Binding
	Select key.Binding
	Edit   key.Binding
	Filter key.Binding
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
		{k.Up, k.Down, k.Left, k.Right},
		{k.Add, k.Select, k.Edit, k.Filter, k.Delete},
		{k.Back, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up:     ui.Key_Up,
	Down:   ui.Key_Down,
	Left:   ui.Key_Left,
	Right:  ui.Key_Right,
	Add:    ui.Key_Add,
	Select: ui.Key_Select,
	Edit:   ui.Key_Edit,
	Filter: ui.Key_Filter,
	Delete: ui.Key_Delete,
	Back:   ui.Key_Back,
	Help:   ui.Key_Help,
	Quit:   ui.Key_Quit,
}
