package note

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"

	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
)

// View renders a backend.Note for the random-notes panel with styling and tag
// highlighting. It keeps presentation logic close to the panel while relying
// on the shared backend.Note data model.
func View(n backend.Note, width int, isHighlighted bool, hlTagIdx int, filteredTag string) string {
	iconStyle := lipgloss.NewStyle().Foreground(ui.Color_RandomNotes_Icon).MarginRight(1)
	titleStyle := lipgloss.NewStyle()

	if isHighlighted && hlTagIdx == -1 {
		titleStyle = titleStyle.Underline(true)
	}

	spacer := " "

	icon := iconStyle.Render(n.Icon)
	title := titleStyle.Render(n.Title)

	tags := "{"
	for i, tag := range n.Tags {
		tags += "#"
		tagsStyle := lipgloss.NewStyle()
		tagsStyle = tagsStyle.Foreground(ui.Color_RandomNotes_Tags)
		if tag == filteredTag {
			tagsStyle = tagsStyle.Bold(true)
		}

		if isHighlighted && i == hlTagIdx {
			tags += tagsStyle.Underline(true).Render(tag)
		} else {
			tags += tagsStyle.Render(tag)
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

	firstLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		icon,
		title,
		space,
		tags,
	)

	var secondLine string
	if n.Created.IsZero() {
		secondLine = lipgloss.NewStyle().
			Foreground(ui.Color_Bg).Render(
			strings.Repeat(spacer, width-4),
		)
	} else {
		secondLine = fmt.Sprintf("%s %s %s",
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
			strings.Repeat(spacer, width-lipgloss.Width(secondLine)-4),
		)

		secondLine = lipgloss.JoinHorizontal(
			lipgloss.Top,
			secondLine,
			space,
		)
	}

	leftPad := "  \n  "
	rightPad := leftPad
	if isHighlighted {
		leftPad = "⎧ \n⎩ "
		rightPad = " ⎫\n ⎭"
	}
	leftPad = lipgloss.NewStyle().
		Foreground(ui.Color_RandomNotes_HighlightBorder).
		Render(leftPad)
	rightPad = lipgloss.NewStyle().
		Foreground(ui.Color_RandomNotes_HighlightBorder).
		Render(rightPad)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPad,
		lipgloss.JoinVertical(
			lipgloss.Left,
			firstLine,
			secondLine,
		),
		rightPad,
	)
}

// New assembles a backend.Note instance, keeping construction alongside the
// rendering helpers for the random-notes panel.
func New(fname, path, title, icon string, tags []string, created time.Time) backend.Note {
	return backend.Note{
		Filename: fname,
		Path:     path,
		Title:    title,
		Icon:     icon,
		Tags:     tags,
		Created:  created,
	}
}
