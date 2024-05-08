package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

// COLORS
var unfocusedBorderColor = lipgloss.Color("#353635")
var focusedBorderColor = lipgloss.Color("#e8e8e8")
var notificationForegroundColor = lipgloss.Color("4")
var senderColor = lipgloss.Color("10")

// BORDERS
var unfocusedBorderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(unfocusedBorderColor)
var focusedBorderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(focusedBorderColor)

// TEXT
var notificationTextStyle = lipgloss.NewStyle().Foreground(notificationForegroundColor)
var senderTextStyle = lipgloss.NewStyle().Foreground(senderColor)

// HELP
var helpStyle = help.Styles{ShortSeparator: lipgloss.NewStyle().Foreground(senderColor)}
