package ui

import (
	"encoding/json"
	"strings"

	"github.com/aemery-cb/pequod/internal/common"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	v1 "k8s.io/api/core/v1"
)

type LogModel struct {
	logs     []string
	pod      *v1.Pod
	logsView viewport.Model
	podView  viewport.Model
	focussed int
	ready    bool
}

var titleStyle = func() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
}()

func NewLogModel() LogModel {
	return LogModel{logs: make([]string, 0)}
}

func (m LogModel) Init() tea.Cmd {
	return nil
}

func (m LogModel) headerView() string {
	title := titleStyle.Render("Log View")
	if m.pod != nil {
		title = titleStyle.Render(m.pod.Name)
	}
	line := strings.Repeat("─", max(0, m.logsView.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m LogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	headerHeight := lipgloss.Height(m.headerView())
	verticalHeight := headerHeight + 3

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right":
			m.focussed = 1
		case "left":
			m.focussed = 0
		}
		if m.focussed == 0 {
			m.logsView, cmd = m.logsView.Update(msg)
			return m, cmd
		} else {
			m.podView, cmd = m.podView.Update(msg)
			return m, cmd
		}

	case tea.MouseMsg:
		if m.focussed == 0 {
			m.logsView, cmd = m.logsView.Update(msg)
			return m, cmd
		} else {
			m.podView, cmd = m.podView.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		if !m.ready {

			m.logsView = viewport.New(msg.Width/2, msg.Height-verticalHeight)
			m.podView = viewport.New(msg.Width/2, msg.Height-verticalHeight)

			m.logsView.YPosition = headerHeight
			m.podView.YPosition = headerHeight

			m.logsView.SetContent(strings.Join(m.logs, ""))
			json, _ := json.MarshalIndent(m.pod, "", " ")
			m.podView.SetContent(string(json))

			m.ready = true
		} else {
			m.logsView.Width = msg.Width / 2

			m.logsView.Height = msg.Height - verticalHeight

			m.podView.Width = msg.Width - m.logsView.Width
			m.logsView.Height = msg.Height - verticalHeight
		}

	case common.NewLogMsg:
		m.pod = msg.Pod
		m.logs = append(m.logs, msg.Message)
		m.logsView.SetContent(strings.Join(m.logs, ""))
		json, _ := json.MarshalIndent(m.pod, "", " ")
		m.podView.SetContent(string(json))
		return m, common.WaitForActivity()
	case common.ClearPodLogsMsg:
		m.logs = make([]string, 0)
	}

	m.logsView, cmd = m.logsView.Update(msg)
	cmds = append(cmds, cmd)
	m.podView, cmd = m.podView.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m LogModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	return lipgloss.JoinHorizontal(0, m.headerView()+"\n"+m.logsView.View(), m.podView.View())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
