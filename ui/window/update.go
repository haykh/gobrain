package window

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
)

func (w *Window) FetchWeather() {
	// if resp, err := http.Get(fmt.Sprintf("https://wttr.in/?m&format=%s", ui.Fmt_Weather)); err == nil {
	// 	w.DebugLog("Fetching weather data")
	// 	defer resp.Body.Close()
	//
	// 	body, err := io.ReadAll(resp.Body)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	w.weather_last_updated = time.Now()
	// 	w.weather = string(body)
	// }
	w.weather_last_updated = time.Now()
	w.weather = "hello"
}

func (w Window) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if w.weather == "" || time.Since(w.weather_last_updated) > ui.UpdateInterval_Weather {
		w.FetchWeather()
	}

	if !w.mdviewport_show {
		active_panel_msg := w.panels[w.active_panel].Update(msg)
		cmd = tea.Batch(active_panel_msg)
	}

	if w.is_typing {
		var cmd_typing tea.Cmd
		w.app.TypingInput, cmd_typing = w.app.TypingInput.Update(msg)
		cmd = tea.Batch(cmd, cmd_typing)

		switch msg.(type) {
		case ui.TypingEndMsg:
			w.is_typing = false
			w.app.TypingInput.Blur()
			return w, func() tea.Msg {
				return ui.DebugMsg{Message: "Ended typing"}
			}
		}

	} else {

		switch msg := msg.(type) {
		case ui.TypingStartMsg:
			w.is_typing = true
			w.app.TypingInput.Focus()
			return w, func() tea.Msg {
				return ui.DebugMsg{Message: "Started typing"}
			}

		case ui.NavigateFwdMsg:
			if msg.NewPanel != w.active_panel {
				w.active_panel = msg.NewPanel
				w.panels[w.active_panel].Sync()
			}
			w.DebugLog(fmt.Sprintf("Navigating to panel: %s", w.panels[w.active_panel].Path()))
			return w, nil

		case ui.OpenEditorMsg:
			return w, backend.OpenEditor(msg.Filename)

		case backend.EditorFinishedMsg:
			if msg.Error != nil {
				w.DebugLog(fmt.Sprintf("Error opening editor: %v", msg.Error))
				return w, tea.Quit
			}
			w.panels[w.active_panel].Sync()
			return w, nil

		case ui.NewViewer:
			if !w.mdviewport_show {
				glamourRenderWidth := ui.Width_Window - w.mdviewport.Style.GetHorizontalFrameSize()
				renderer, err := glamour.NewTermRenderer(
					glamour.WithStandardStyle("dark"),
					glamour.WithWordWrap(glamourRenderWidth),
				)
				if err != nil {
					return nil, func() tea.Msg {
						return ui.ErrorMsg{Error: err}
					}
				}
				w.mdviewport_show = true
				w.mdviewport_filename = msg.Filename
				w.mdviewport = viewport.New(ui.Width_Window-1, ui.Height_Window)
				content, err := backend.ReadMarkdownNote(msg.Filepath, msg.Filename)
				if err != nil {
					return nil, func() tea.Msg {
						return ui.ErrorMsg{Error: fmt.Errorf("error reading file %s: %w", msg.Filename, err)}
					}
				}
				str, err := renderer.Render(content)
				if err != nil {
					return nil, func() tea.Msg {
						return ui.ErrorMsg{Error: err}
					}
				}

				w.mdviewport.SetContent(str)
			} else {
				w.DebugLog("Markdown viewer already open")
				return w, tea.Quit
			}
			return w, nil

		case ui.TrashNoteMsg:
			if err := w.app.TrashNote(msg.Filename, msg.Filepath); err != nil {
				w.DebugLog(fmt.Sprintf("Error trashing random note: %v", err))
				return w, tea.Quit
			}
			w.panels[w.active_panel].Sync()
			return w, nil

		case tea.KeyMsg:
			switch {
			case key.Matches(msg, ui.Key_Back):
				if w.mdviewport_show {
					w.mdviewport_show = false
					w.mdviewport_filename = ""
				} else {
					new_panel := w.panels[w.active_panel].Parent()
					if new_panel != w.active_panel {
						w.active_panel = new_panel
						w.panels[w.active_panel].Sync()
					}
				}
				return w, nil

			case key.Matches(msg, ui.Key_Help):
				w.help.ShowAll = !w.help.ShowAll

			case key.Matches(msg, ui.Key_Quit):
				return w, tea.Quit

			default:
				if w.mdviewport_show {
					var cmd_md tea.Cmd
					w.mdviewport, cmd_md = w.mdviewport.Update(msg)
					return w, tea.Batch(cmd, cmd_md)
				}
			}
		}
	}

	switch msg := msg.(type) {
	case ui.ErrorMsg:
		if msg.Error != nil {
			w.DebugLog(fmt.Sprintf("Error: %v", msg.Error))
			return w, tea.Quit
		} else {
			return w, nil
		}

	case ui.DebugMsg:
		w.DebugLog(msg.Message)
		return w, nil
	}

	return w, cmd
}
