package dashboard

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/haykh/gobrain/ui"
)

func (m model) View() string {
	// general params
	column_width := ui.Width_Window / 3

	large_button_style := lipgloss.NewStyle().
		Padding(ui.PaddingV_LargeBtn, ui.PaddingH_LargeBtn).
		Align(lipgloss.Center).
		Background(ui.Color_InactiveBg)

	small_button_style := lipgloss.NewStyle().
		Padding(ui.PaddingV_SmallBtn, ui.PaddingH_SmallBtn)

	dly_style := large_button_style
	tsk_style := large_button_style
	rnd_style := large_button_style

	todaysnote_style := small_button_style
	newnote_style := small_button_style

	column_style := lipgloss.NewStyle().
		Width(column_width).
		Align(lipgloss.Center)

	// apply highlight styles
	switch m.cursor {
	case DailyNotes:
		dly_style = dly_style.Background(ui.Color_ActiveBg)
	case TaskLists:
		tsk_style = tsk_style.Background(ui.Color_ActiveBg)
	case RandomNotes:
		rnd_style = rnd_style.Background(ui.Color_ActiveBg)
	case TodaysNote:
		todaysnote_style = todaysnote_style.Background(ui.Color_ActiveFg_Today)
	case NewRandomNote:
		newnote_style = newnote_style.Background(ui.Color_ActiveFg_New)
	}

	place_large_button := func(btn string) string {
		return lipgloss.Place(
			column_width, ui.MarginT_LargeBtn+4,
			lipgloss.Center, lipgloss.Bottom,
			btn,
			lipgloss.WithWhitespaceChars(ui.String_Bg),
			lipgloss.WithWhitespaceForeground(ui.Color_Bg),
		)
	}

	place_small_button := func(btn string) string {
		return lipgloss.Place(
			column_width, ui.MarginT_SmallBtn+1,
			lipgloss.Center, lipgloss.Bottom,
			btn,
			lipgloss.WithWhitespaceChars(ui.String_Bg),
			lipgloss.WithWhitespaceForeground(ui.Color_Bg),
		)
	}

	renderUrgentTasks := func() string {
		if len(m.urgentTasks) == 0 {
			return place_small_button(small_button_style.Render("no urgent tasks"))
		}

		taskLines := []string{}
		for i, task := range m.urgentTasks {
			line := fmt.Sprintf("%d. %s", i+1, task.Text)
			if !task.DueDate.IsZero() {
				line = fmt.Sprintf("%s — %s", line, task.DueDate.Format("Jan 02"))
			} else {
				line = fmt.Sprintf("%s — due soon", line)
			}
			taskLines = append(taskLines, line)
		}

		urgentList := small_button_style.
			Align(lipgloss.Left).
			Render(lipgloss.JoinVertical(lipgloss.Left, taskLines...))

		height := ui.MarginT_SmallBtn + lipgloss.Height(urgentList)
		if height < ui.MarginT_SmallBtn+1 {
			height = ui.MarginT_SmallBtn + 1
		}
		return lipgloss.Place(
			column_width, height,
			lipgloss.Center, lipgloss.Bottom,
			urgentList,
			lipgloss.WithWhitespaceChars(ui.String_Bg),
			lipgloss.WithWhitespaceForeground(ui.Color_Bg),
		)
	}

	dly_button := place_large_button(dly_style.Render(ui.Title("daily")))
	tsk_button := place_large_button(tsk_style.Render(ui.Title("tasks")))
	rnd_buttom := place_large_button(rnd_style.Render(ui.Title("random")))

	todaysnote_button := place_small_button(todaysnote_style.Render("today"))
	urgent_tasks := renderUrgentTasks()
	newnote_button := place_small_button(newnote_style.Render("new"))

	dly_column := lipgloss.JoinVertical(
		lipgloss.Center,
		dly_button,
		todaysnote_button,
	)

	tsk_column := lipgloss.JoinVertical(
		lipgloss.Center,
		tsk_button,
		urgent_tasks,
	)

	rnd_column := lipgloss.JoinVertical(
		lipgloss.Center,
		rnd_buttom,
		newnote_button,
	)

	dashboard_content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		column_style.Render(dly_column),
		column_style.Render(tsk_column),
		column_style.Render(rnd_column),
	)

	return dashboard_content
}
