package dashboard

import (
	"github.com/charmbracelet/bubbles/key"

	"github.com/haykh/gobrain/ui"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Edit   key.Binding
	Select key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Select, k.Edit, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up:     ui.Key_Up,
	Down:   ui.Key_Down,
	Left:   ui.Key_Left,
	Right:  ui.Key_Right,
	Edit:   ui.Key_Edit,
	Select: ui.Key_Select,
	Help:   ui.Key_Help,
	Quit:   ui.Key_Quit,
}
