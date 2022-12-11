package ui

import (
	"github.com/aemery-cb/pequod/internal/common"
	tea "github.com/charmbracelet/bubbletea"
)

type LogModel struct {
	logs []string
}

func NewLogModel() LogModel {
	return LogModel{logs: make([]string, 0)}
}

func (m LogModel) Init() tea.Cmd {
	return nil
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

	if len(m.logs) == 0 {
		return "press tab then select a pod to stream logs from"
	}

	for _, log := range m.logs {
		s += log
	}

	s += "\nPress 'tab' to go back"
	return s
}
