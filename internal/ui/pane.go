package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Pane struct {
	model tea.Model
	focus bool
}

func NewPane(model tea.Model) Pane {
	return Pane{model: model, focus: false}

}

func (p Pane) Init() tea.Cmd {
	return p.model.Init()
}

func (p *Pane) Focus() {
	p.focus = true
}
func (p *Pane) Blur() {
	p.focus = false
}

func (p Pane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok && !p.focus {
		return p, nil
	}

	model, cmd := p.model.Update(msg)
	p.model = model

	return p, cmd
}

func (p Pane) View() string {
	return p.model.View()
}
