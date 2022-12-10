package ui

import (
	"github.com/aemery-cb/pequod/internal/common"
	tea "github.com/charmbracelet/bubbletea"
)

type LogModel struct {
	focus bool
	logs  []string
}

func NewLogModel() LogModel {
	return LogModel{focus: false, logs: make([]string, 0)}
}

func (m LogModel) Init() tea.Cmd {
	return nil
}

func (m LogModel) Focus() Page {
	m.focus = true
	return m
}

func (m LogModel) Blur() Page {
	m.focus = false
	return m
}

func (m LogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.NewLogMsg:
		m.logs = append(m.logs, msg.Message)
		return m, common.WaitForActivity()
	case common.ClearPodLogsMsg:
		m.logs = make([]string, 0)
	}
	return m, nil
}

func (m LogModel) View() string {
	s := ""

	for _, log := range m.logs {
		s += log
	}

	s += "\nPress 'tab' to go back"
	return s
}
