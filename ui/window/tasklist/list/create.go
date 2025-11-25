package list

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type NewList struct {
	textInput textinput.Model
	err       error
}

//	func (m NewList) Init() tea.Cmd {
//		return textinput.Blink
//	}
func (m *NewList) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return tea.Quit
		}

	case error:
		m.err = msg
		return nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return cmd
}

func (m NewList) View() string {
	return m.textInput.View()
}
