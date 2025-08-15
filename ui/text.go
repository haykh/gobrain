package ui

import (
	_ "embed"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

//go:embed phm-lcdmatrix.flf
var fontFile []byte

func Title(text string) string {
	io_reader := strings.NewReader(string(fontFile))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		figure.NewFigureWithFont(text, io_reader, false).Slicify()...,
	)
}
