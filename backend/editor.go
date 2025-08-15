package backend

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type EditorFinishedMsg struct{ Error error }

func OpenEditor(filepath string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, filepath) //nolint:gosec
	return tea.ExecProcess(
		c,
		func(err error) tea.Msg {
			return EditorFinishedMsg{err}
		},
	)
}
