package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui/window"
)

const ShowHelp = true
const HideHelp = false
const NoDebug = false
const Debug = true

func main() {
	app := backend.New()
	app.Init()
	p := tea.NewProgram(window.New(app, HideHelp, NoDebug))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
