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
	client   *Client
	pods     []PodStatus
	quitting bool
}

func newModel(client *Client) model {
	pods, _ := client.GetPods()
	return model{
		client: client,
		pods:   pods,
	}
}

func (m model) Init() tea.Cmd {
	return spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.quitting = true
			return m, tea.Quit
		case "r":
			pods, _ := m.client.GetPods()
			m.pods = pods
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	var s string

	if m.quitting {
		s += "That’s all for today!"
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