package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

var defaultStyle = lipgloss.NewStyle()

// COLORS
var successColor = lipgloss.Color("#007d69")
var errorColor = lipgloss.Color("#ff3333")
var unfocusedBorderColor = lipgloss.Color("#353635")
var focusedBorderColor = lipgloss.Color("#e8e8e8")
var notificationForegroundColor = lipgloss.Color("4")
var senderColor = lipgloss.Color("10")

// BORDERS
var unfocusedBorderStyle = defaultStyle.Copy().BorderStyle(lipgloss.NormalBorder()).BorderForeground(unfocusedBorderColor)
var focusedBorderStyle = defaultStyle.Copy().BorderStyle(lipgloss.NormalBorder()).BorderForeground(focusedBorderColor)

// TEXT
var notificationTextStyle = defaultStyle.Copy().Foreground(notificationForegroundColor)
var senderTextStyle = defaultStyle.Copy().Foreground(senderColor)
var errorTextStyle = defaultStyle.Copy().Foreground(errorColor)
var successTextStyle = defaultStyle.Copy().Foreground(successColor)

// HELP
var helpStyle = help.Styles{ShortSeparator: defaultStyle.Copy().Foreground(senderColor)}
