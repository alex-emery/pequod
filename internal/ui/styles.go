package ui

import "github.com/charmbracelet/lipgloss"

// keep all the styles in a single place so we can work out different heights/widths for views

var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))
