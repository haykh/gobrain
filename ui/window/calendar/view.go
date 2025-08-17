package calendar

import (
	"math"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/haykh/gobrain/ui"
)

func isToday(date time.Time) bool {
	today := time.Now()
	return date.Year() == today.Year() && date.Month() == today.Month() && date.Day() == today.Day()
}

func (c CalendarDay) View(is_highlighted bool) string {
	day_view := c.Date.Format("2")
	note_view := ""
	if c.NoteExists {
		note_view = "*"
	}
	width := int(math.Floor(float64(ui.Width_Window)/float64(7)) - 4)
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	if isToday(c.Date) {
		style = style.BorderForeground(ui.Color_Border_CalendarDay_Today)
	}
	if is_highlighted {
		if isToday(c.Date) {
			style = style.BorderForeground(ui.Color_Border_CalendarDay_Today_Active)
		}
		style = style.BorderBackground(ui.Color_Bg_CalendarDay_Active)
		style = style.Background(ui.Color_Bg_CalendarDay_Active)
	}

	return style.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.NewStyle().Width(width).MarginBottom(1).Render(day_view),
			),
			note_view,
		),
	)
}

func monthColumn(months []string) string {
	month_fmt := []string{}
	for _, month := range months {
		month_fmt = append(
			month_fmt,
			lipgloss.NewStyle().
				Margin(2, 2, 2, 2).
				Foreground(ui.Color_Fg_Calendar_Helper).
				Render(month),
		)
	}
	return lipgloss.JoinVertical(lipgloss.Left, month_fmt...)
}

func (c model) View() string {
	weekdays_str := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	weekdays_fmt := []string{}
	for _, wd := range weekdays_str {
		weekdays_fmt = append(
			weekdays_fmt,
			lipgloss.NewStyle().
				Foreground(ui.Color_Fg_Calendar_Helper).
				Margin(0, 0, 0, 7).Render(wd),
		)
	}
	weeks := []string{}
	week := []string{}
	week_time := []time.Time{}
	left_months := []string{}
	right_months := []string{}
	for di, day := range c.calendar_days {
		if di > 0 && di%7 == 0 {
			weeks = append(weeks, lipgloss.JoinHorizontal(lipgloss.Left, week...))
			left_months = append(left_months, week_time[0].Format("Jan"))
			right_months = append(right_months, week_time[len(week_time)-1].Format("Jan"))
			week = []string{}
			week_time = []time.Time{}
		}
		week = append(week, day.View(di == c.cursor))
		week_time = append(week_time, day.Date)
	}
	if len(week) > 0 {
		weeks = append(weeks, lipgloss.JoinHorizontal(lipgloss.Left, week...))
		left_months = append(left_months, week_time[0].Format("Jan"))
		right_months = append(right_months, week_time[len(week_time)-1].Format("Jan"))
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().MarginLeft(3).Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				weekdays_fmt...,
			)),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			monthColumn(left_months),
			lipgloss.JoinVertical(
				lipgloss.Left,
				weeks...,
			),
			monthColumn(right_months),
		),
	)
}
