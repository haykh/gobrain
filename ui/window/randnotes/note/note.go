package note

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/haykh/gobrain/ui"
)

type Note struct {
	Filename string
	Path     string
	Title    string
	Icon     string
	Tags     []string
	Created  time.Time
}

func New(fname, path, title, icon string, tags []string, created time.Time) Note {
	return Note{
		Filename: fname,
		Path:     path,
		Title:    title,
		Icon:     icon,
		Tags:     tags,
		Created:  created,
	}
}

func (n Note) View(width int, is_highlighted bool, hl_tag_idx int, filtered_tag string) string {
	icon_style := lipgloss.NewStyle().Foreground(ui.Color_RandomNotes_Icon).MarginRight(1)

	title_style := lipgloss.NewStyle()

	if is_highlighted && hl_tag_idx == -1 {
		title_style = title_style.Underline(true)
	}

	spacer := " "

	icon := icon_style.Render(n.Icon)
	title := title_style.Render(n.Title)

	tags := "{"
	for i, tag := range n.Tags {
		tags += "#"
		tags_style := lipgloss.NewStyle()
		tags_style = tags_style.Foreground(ui.Color_RandomNotes_Tags)
		if tag == filtered_tag {
			tags_style = tags_style.Bold(true)
		}

		if is_highlighted && i == hl_tag_idx {
			tags += tags_style.Underline(true).Render(tag)
		} else {
			tags += tags_style.Render(tag)
		}
		if i < len(n.Tags)-1 {
			tags += ", "
		}
	}
	tags += "}"

	if len(n.Tags) == 0 {
		tags = ""
	}
	tags = lipgloss.NewStyle().Render(tags)

	space := lipgloss.NewStyle().
		Foreground(ui.Color_Bg).Render(
		strings.Repeat(spacer, width-(lipgloss.Width(icon)+lipgloss.Width(title)+lipgloss.Width(tags))-4),
	)

	first_line := lipgloss.JoinHorizontal(
		lipgloss.Top,
		icon,
		title,
		space,
		tags,
	)

	var second_line string
	if n.Created.IsZero() {
		second_line = lipgloss.NewStyle().
			Foreground(ui.Color_Bg).Render(
			strings.Repeat(spacer, width-4),
		)
	} else {
		second_line = fmt.Sprintf("%s %s %s",
			lipgloss.NewStyle().
				Foreground(ui.Color_Dividers).Render(
				"  [ created ",
			),
			lipgloss.NewStyle().
				Foreground(ui.Color_RandomNotes_Created).Render(
				humanize.Time(n.Created),
			),
			lipgloss.NewStyle().
				Foreground(ui.Color_Dividers).Render(
				"]",
			),
		)
		space = lipgloss.NewStyle().
			Foreground(ui.Color_Bg).Render(
			strings.Repeat(spacer, width-lipgloss.Width(second_line)-4),
		)

		second_line = lipgloss.JoinHorizontal(
			lipgloss.Top,
			second_line,
			space,
		)
	}

	left_pad := "  \n  "
	right_pad := left_pad
	if is_highlighted {
		left_pad = "⎧ \n⎩ "
		right_pad = " ⎫\n ⎭"
	}
	left_pad = lipgloss.NewStyle().
		Foreground(ui.Color_RandomNotes_HighlightBorder).
		Render(left_pad)
	right_pad = lipgloss.NewStyle().
		Foreground(ui.Color_RandomNotes_HighlightBorder).
		Render(right_pad)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		left_pad,
		lipgloss.JoinVertical(
			lipgloss.Left,
			first_line,
			second_line,
		),
		right_pad,
	)
}
