package ui

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	/**
	 * Main Dashboard
	**/
	//// Colors
	// general colors
	Color_Border_Window = lipgloss.Color("236")
	Color_InactiveBg    = lipgloss.Color("236")
	Color_ActiveBg      = lipgloss.Color("57")
	Color_Dividers      = lipgloss.Color("245")

	// element colors
	Color_ActiveFg_Today       = lipgloss.Color("63")
	Color_ActiveFg_New         = lipgloss.Color("200")
	Color_ActiveFg_Breadcrumbs = lipgloss.Color("63")
	Color_Fg_Breadcrumbs       = lipgloss.Color("245")
	Color_Fg_Datetime          = lipgloss.Color("4")
	Color_Fg_Weather           = lipgloss.Color("12")
	Color_Bg_DebugLabel        = lipgloss.Color("124")

	// background
	Color_Bg = lipgloss.Color("236")

	//// Sizing
	// window params
	Width_Window   = 90
	Height_Window  = 30
	MarginV_Window = 2
	MarginH_Window = 1

	// large button params
	PaddingH_LargeBtn = 3
	PaddingV_LargeBtn = 1
	MarginT_LargeBtn  = 2

	// small button params
	PaddingH_SmallBtn = 1
	PaddingV_SmallBtn = 0
	MarginT_SmallBtn  = 2

	//// Formats and strings
	String_Divider_Time = "⋮"
	String_Bg           = "⢕"
	Fmt_Weather         = "%C+%t"
	Fmt_Time            = "3:04pm"
	Fmt_Date            = "Mon, Jan 2"

	//// Behavior
	UpdateInterval_Weather = 15 * time.Minute

	/**
	 * Random Notes Window
	**/
	//// Colors
	Color_RandomNotes_Icon            = lipgloss.Color("63")
	Color_RandomNotes_Tags            = lipgloss.Color("10")
	Color_RandomNotes_FilteredTags    = lipgloss.Color("34")
	Color_RandomNotes_Created         = lipgloss.Color("9")
	Color_RandomNotes_HighlightBorder = lipgloss.Color("7")

	/**
	 * Common Elements
	**/
	Color_Scrollbar = lipgloss.Color("237")
)

var Style_Window = lipgloss.NewStyle().
	Margin(MarginV_Window, MarginH_Window).
	Align(lipgloss.Center).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Color_Border_Window)
