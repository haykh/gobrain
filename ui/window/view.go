package window

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/haykh/gobrain/components"
	"github.com/haykh/gobrain/ui"
)

func (w Window) View() string {
	if w.mdviewport_show {
		tot_lines := w.mdviewport.TotalLineCount()
		vis_lines := w.mdviewport.Height
		first_line := w.mdviewport.YOffset

		return w.FinalView(lipgloss.JoinHorizontal(
			lipgloss.Top,
			w.mdviewport.View(),
			components.Scrollbar(
				ui.Height_Window,
				tot_lines,
				first_line,
				vis_lines,
			),
		))
	} else {
		return w.FinalView(w.Active().View())
	}
}

func (w Window) FinalView(content string) string {
	view := ui.Style_Window.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			w.BreadcrumbsView(),
			lipgloss.Place(
				ui.Width_Window, ui.Height_Window,
				lipgloss.Center, lipgloss.Top,
				content,
				lipgloss.WithWhitespaceChars(ui.String_Bg),
				lipgloss.WithWhitespaceForeground(ui.Color_Bg),
			),
			w.StatusView(),
		),
	)
	help_style := lipgloss.NewStyle().MarginTop(5)

	help_view := ""
	if w.show_help {
		help_view = help_style.Render(w.HelpView())
	}

	if w.debug {
		nrows := ui.Height_Window - lipgloss.Height(help_view) - 20
		debug_style := lipgloss.NewStyle().MarginTop(nrows)
		debug_view := debug_style.Render(w.DebugView())
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			view,
			lipgloss.JoinVertical(
				lipgloss.Left,
				help_view,
				debug_view,
			),
		)
	} else {
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			view,
			help_view,
		)
	}
}

func (w Window) BreadcrumbsView() string {
	breadcrumbs_style := lipgloss.NewStyle().
		Width(ui.Width_Window).
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false).
		BorderForeground(ui.Color_Border_Window)

	divider_style := lipgloss.NewStyle().
		Foreground(ui.Color_Dividers)
	element_style := lipgloss.NewStyle().
		Foreground(ui.Color_Fg_Breadcrumbs)
	active_style := lipgloss.NewStyle().
		Foreground(ui.Color_ActiveFg_Breadcrumbs)

	breadcrumbs_text := ""
	for i, part := range w.Active().Path() {
		breadcrumbs_text += divider_style.Render(" / ")
		if i == len(w.Active().Path())-1 && !w.mdviewport_show {
			breadcrumbs_text += active_style.Render(part)
		} else {
			breadcrumbs_text += element_style.Render(part)
		}
	}
	if w.mdviewport_show {
		breadcrumbs_text += divider_style.Render(" / ")
		breadcrumbs_text += active_style.Render(w.mdviewport_filename)
	}
	return breadcrumbs_style.Render(breadcrumbs_text)
}

func (w Window) StatusView() string {
	status_style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(ui.Width_Window).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderForeground(ui.Color_Border_Window)

	datetime_style := lipgloss.NewStyle().
		Foreground(ui.Color_Fg_Datetime)
	divider_style := lipgloss.NewStyle().
		Foreground(ui.Color_Dividers)

	time_now := time.Now().Format(ui.Fmt_Date)
	date_now := time.Now().Format(ui.Fmt_Time)

	datetime_text := fmt.Sprintf(" %s %s %s",
		datetime_style.Render(time_now),
		divider_style.Render(ui.String_Divider_Time),
		datetime_style.Render(date_now))

	remote_text := ""
	if sync, err := w.app.InSync(); err != nil {
		panic(err)
	} else if sync {
		remote_text = lipgloss.NewStyle().Foreground(ui.Color_Fg_RemoteSynced).Render(ui.String_Synced)
	} else {
		remote_text = lipgloss.NewStyle().Foreground(ui.Color_Fg_RemoteUnsynced).Render(ui.String_Desynced)
	}
	remote_text += lipgloss.NewStyle().Foreground(ui.Color_Fg_RemoteDivider).Render(" " + ui.String_Divider_Sync + " ")
	if w.app.OfflineMode() {
		remote_text += lipgloss.NewStyle().Foreground(ui.Color_Fg_RemoteOffline).Render(ui.String_Offline)
	} else {
		remote_text += lipgloss.NewStyle().Foreground(ui.Color_Fg_RemoteOnline).Render(ui.String_Online)
	}

	weather_style := lipgloss.NewStyle().
		Foreground(ui.Color_Fg_Weather).
		MaxWidth(ui.Width_Window - lipgloss.Width(datetime_text) - lipgloss.Width(remote_text) - 1)

	weather_text := weather_style.Render(w.weather)

	nspace := ui.Width_Window - lipgloss.Width(weather_text) - lipgloss.Width(datetime_text) - lipgloss.Width(remote_text) - 1
	nspace_l := nspace / 3
	nspace_r := nspace - nspace_l
	return status_style.Render(
		datetime_text + strings.Repeat(" ", nspace_l) + remote_text + strings.Repeat(" ", nspace_r) + weather_text + " ",
	)
}

func (w Window) HelpView() string {
	return w.help.View(w.Active().Keys())
}

func (w Window) DebugView() string {
	debug_label_style := lipgloss.NewStyle().
		Background(ui.Color_Bg_DebugLabel).
		Padding(0, 1)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		debug_label_style.Render("DEBUG"),
		lipgloss.JoinVertical(
			lipgloss.Left,
			w.debug_logs...,
		),
	)
}
