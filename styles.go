package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	// AppStyle is the main container for the entire application.
	// We use a simple margin to keep content from touching the terminal edges.
	AppStyle = lipgloss.NewStyle().Margin(1, 2)

	// TitleStyle is for headers.
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#6A5ACD")).
			Padding(0, 1)

	// ErrorStyle is for displaying errors.
	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5370"))

	// PaginationStyle is for the list's pagination.
	PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)

	// HelpStyle is for the list's help text.
	HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)
