package dashboard

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mergestat/timediff"

	"github.com/haykh/gobrain/ui"
)

func (m model) View() string {
	// general params
	column_width := ui.Width_Window / 3
	urgent_width := 3 * ui.Width_Window / 4
	urgent_offset_l := (ui.Width_Window - urgent_width) / 2
	empty_line := lipgloss.Place(
		ui.Width_Window, 1,
		lipgloss.Left, lipgloss.Top,
		"",
		lipgloss.WithWhitespaceChars(ui.String_Bg),
		lipgloss.WithWhitespaceForeground(ui.Color_Bg),
	)

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
			return "no urgent tasks"
		}

		taskLines := []string{}
		dueTimes := []string{}
		dueTimeCategories := []int{}
		max_dueTime_width := 0
		for _, task := range m.urgentTasks {
			line := fmt.Sprintf("· %s", task.Text)
			dueTime := ""
			if !task.DueDate.IsZero() {
				dueTime = task.DueDate.Format("Jan 02")
				if time.Now().After(task.DueDate) {
					dueTime = "overdue"
				} else if time.Now().Add(30 * 24 * time.Hour).After(task.DueDate) {
					dueTime = timediff.TimeDiff(task.DueDate)
				}
				if (time.Now().Add(48 * time.Hour)).After(task.DueDate) {
					dueTimeCategories = append(dueTimeCategories, 0)
				} else if (time.Now().Add(7 * 24 * time.Hour)).After(task.DueDate) {
					dueTimeCategories = append(dueTimeCategories, 1)
				} else {
					dueTimeCategories = append(dueTimeCategories, 2)
				}
				dueTime = fmt.Sprintf("[%s]", dueTime)
			}
			if lipgloss.Width(dueTime) > max_dueTime_width {
				max_dueTime_width = lipgloss.Width(dueTime)
			}
			dueTimes = append(dueTimes, dueTime)
			taskLines = append(taskLines, line)
		}

		splitLines := []string{}
		for i, line := range taskLines {
			dueTime := dueTimes[i]
			dueTimeCategory := dueTimeCategories[i]
			for lipgloss.Width(line) > urgent_width-max_dueTime_width {
				splitAt := urgent_width - max_dueTime_width
				for j := urgent_width - max_dueTime_width; j > 0; j-- {
					if line[j] == ' ' {
						splitAt = j
						break
					}
				}
				splitLines = append(splitLines, line[:splitAt])
				line = line[splitAt:]
			}
			// last line with due time
			padding := urgent_width - lipgloss.Width(line) - lipgloss.Width(dueTime)
			dots := lipgloss.NewStyle().Foreground(ui.Color_Bg_UrgentTasks).Render(strings.Repeat(".", padding))
			switch dueTimeCategory {
			case 0:
				dueTime = lipgloss.NewStyle().Foreground(ui.Color_Fg_UrgentTasks_DateUrgent).Render(dueTime)
			case 1:
				dueTime = lipgloss.NewStyle().Foreground(ui.Color_Fg_UrgentTasks_DateSoon).Render(dueTime)
			default:
				dueTime = lipgloss.NewStyle().Foreground(ui.Color_Fg_UrgentTasks_Date).Render(dueTime)
			}
			splitLines = append(splitLines, fmt.Sprintf("%s%s%s", line, dots, dueTime))
		}

		taskLines = splitLines

		return lipgloss.JoinVertical(
			lipgloss.Left,
			taskLines...,
		)
	}

	dly_button := place_large_button(dly_style.Render(ui.Title("daily")))
	tsk_button := place_large_button(tsk_style.Render(ui.Title("tasks")))
	rnd_buttom := place_large_button(rnd_style.Render(ui.Title("random")))

	todaysnote_button := place_small_button(todaysnote_style.Render("today"))
	dummy := place_small_button("")
	newnote_button := place_small_button(newnote_style.Render("new"))
	urgent_tasks := renderUrgentTasks()
	tasks_height := lipgloss.Height(urgent_tasks)
	tasks_width := lipgloss.Width(urgent_tasks)
	urgent_tasks_filler := func(width int) string {
		return lipgloss.Place(
			width, tasks_height,
			lipgloss.Center, lipgloss.Bottom,
			"",
			lipgloss.WithWhitespaceChars(ui.String_Bg),
			lipgloss.WithWhitespaceForeground(ui.Color_Bg),
		)
	}
	urgent_tasks = lipgloss.JoinHorizontal(
		lipgloss.Center,
		urgent_tasks_filler(urgent_offset_l),
		urgent_tasks,
		urgent_tasks_filler(ui.Width_Window-urgent_offset_l-tasks_width),
	)

	dly_column := lipgloss.JoinVertical(
		lipgloss.Center,
		dly_button,
		todaysnote_button,
	)

	tsk_column := lipgloss.JoinVertical(
		lipgloss.Center,
		tsk_button,
		dummy,
	)

	rnd_column := lipgloss.JoinVertical(
		lipgloss.Center,
		rnd_buttom,
		newnote_button,
	)

	dashboard_content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			column_style.Render(dly_column),
			column_style.Render(tsk_column),
			column_style.Render(rnd_column),
		),
		empty_line,
		empty_line,
		urgent_tasks,
	)

	return dashboard_content
}
