package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

type PodMsg struct {
	pods []PodStatus
}
type PodStatus struct {
	name   string
	status string
	uptime string
}

func (r PodStatus) String() string {
	return r.name + " : " + r.status + " - " + r.uptime
}

type model struct {
	spinner  spinner.Model
	pods     []PodStatus
	quitting bool
}

func newModel() model {
	s := spinner.New()
	s.Style = spinnerStyle
	return model{
		spinner: s,
		pods:    make([]PodStatus, 5),
	}
}

func (m model) Init() tea.Cmd {
	return spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case PodMsg:
		m.pods = msg.pods
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	var s string

	if m.quitting {
		s += "That’s all for today!"
	} else {
		s += m.spinner.View() + "Watching..."
	}

	s += "\n\n"

	for _, res := range m.pods {
		s += res.String() + "\n"
	}

	if !m.quitting {
		s += "Press any key to exit"
	}

	if m.quitting {
		s += "\n"
	}

	return s
}
