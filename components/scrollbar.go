package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/haykh/gobrain/ui"
)

func Scrollbar(height, num_items, first_item, num_displayed_items int) string {
	if num_displayed_items == 0 || num_displayed_items == num_items {
		return " "
	}
	first_idx := int(float64(first_item) / float64(num_items) * float64(height))
	last_idx := int(float64(first_item+num_displayed_items) / float64(num_items) * float64(height))
	bars := []string{}
	for i := range height {
		if i < first_idx || i >= last_idx {
			bars = append(bars, "░")
		} else {
			bars = append(bars, "█")
		}
	}
	return lipgloss.NewStyle().
		Foreground(ui.Color_Scrollbar).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				bars...,
			),
		)
}
