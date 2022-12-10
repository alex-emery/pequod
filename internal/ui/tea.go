package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

/**
This is just an empty model with a minimum implementation.
**/

type Model struct {
}

func NewModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "new view"
}
