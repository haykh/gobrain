package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v3"

	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui/window"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	default_home := filepath.Join(usr.HomeDir, ".gobrain")
	cmd := &cli.Command{
		Name:    "gobrain",
		Usage:   "a terminal-based notes and tasks organizer",
		Version: "v1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "home",
				Aliases: []string{"h"},
				Value:   default_home,
				Usage:   "home directory for gobrain",
			},
			&cli.StringFlag{
				Name:    "git",
				Aliases: []string{"g"},
				Value:   "",
				Usage:   "initialize home directory from git repository",
			},
			&cli.BoolFlag{
				Name:    "offline",
				Aliases: []string{"o"},
				Usage:   "offline mode",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "enable debug mode",
			},
			&cli.BoolFlag{
				Name:    "keys",
				Aliases: []string{"k"},
				Usage:   "show keybindings",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var app *backend.Backend
			var err error
			if cmd.String("git") != "" {
				if app, err = backend.InitFromGit(cmd.String("home"), cmd.String("git")); err != nil {
					return fmt.Errorf("could not initialize from git: %v", err)
				}
			} else {
				app = backend.New(cmd.String("home"), cmd.Bool("offline"))
				if err = app.Init(); err != nil {
					return fmt.Errorf("could not initialize backend: %v", err)
				}
			}
			p := tea.NewProgram(window.New(app, cmd.Bool("keys"), cmd.Bool("debug")))
			if _, err := p.Run(); err != nil {
				return fmt.Errorf("could not run program: %v", err)
			}
			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
