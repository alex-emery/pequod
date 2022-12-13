package ui

import (
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
	viewport viewport.Model
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
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m LogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	headerHeight := lipgloss.Height(m.headerView())
	verticalHeight := headerHeight + 3

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {

			m.viewport = viewport.New(msg.Width-2, msg.Height-verticalHeight)
			m.viewport.YPosition = headerHeight

			m.viewport.SetContent(strings.Join(m.logs, ""))
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalHeight
		}
	case common.WatchPodLogsMsg:
		m.pod = msg.Pod
		return m, common.WaitForActivity()
	case common.NewLogMsg:
		if m.pod != nil && m.pod.Name != msg.Pod.Name {
			return m, common.WaitForActivity()
		}
		m.pod = msg.Pod
		m.logs = append(m.logs, msg.Message)
		m.viewport.SetContent(strings.Join(m.logs, ""))
		return m, common.WaitForActivity()
	case common.ClearPodLogsMsg:
		m.logs = make([]string, 0)
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m LogModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	return m.headerView() + "\n" + m.viewport.View()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
