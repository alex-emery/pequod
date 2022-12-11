package ui

import (
	"sort"

	"github.com/aemery-cb/pequod/internal/common"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	v1 "k8s.io/api/core/v1"
)

type PodModel struct {
	table table.Model
	pods  []v1.Pod
	ready bool
}

func NewPodModel() PodModel {
	// rows := []table.Row{}
	return PodModel{}
}

func (m PodModel) Init() tea.Cmd {
	return nil
}

func (m PodModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			columnWidth := msg.Width/2 - 3
			var columns = []table.Column{
				{Title: "Namespace", Width: columnWidth},
				{Title: "Name", Width: columnWidth},
			}
			t := table.New(
				table.WithColumns(columns),
				table.WithFocused(true),
			)

			s := table.DefaultStyles()
			s.Header = s.Header.
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")).
				BorderBottom(true).
				Bold(false)
			s.Selected = s.Selected.
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57")).
				Bold(false)
			t.SetStyles(s)
			m.table = t

			m.table.SetHeight(msg.Height - 5)
			m.table.SetWidth(msg.Width)
			m.ready = true
		} else {
			m.table.SetHeight(msg.Height - 5)
			m.table.SetWidth(msg.Width)
		}

	case common.NewPodMsg:
		m.pods = append(m.pods, *msg.Pod)
		sort.Slice(m.pods, func(i, j int) bool {
			return m.pods[i].Name < m.pods[j].Name
		})
		m.UpdateTable()
		return m, common.WaitForActivity()
	case common.UpdatePodMsg:
		return m, common.WaitForActivity()
	case common.DeletePodMsg:
		pods := m.pods
		newList := []v1.Pod{}
		for _, pod := range pods {
			if pod.Name != msg.Pod.Name {
				newList = append(newList, pod)
			}
		}
		m.pods = newList
		m.UpdateTable()
		return m, common.WaitForActivity()
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			index := m.table.Cursor()
			if index <= len(m.pods) {
				selected := m.pods[index]
				return m, common.WatchPodLogs(&selected)
			}
		}
	}

	table, cmd := m.table.Update(msg)
	m.table = table
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *PodModel) UpdateTable() {
	var rows []table.Row
	for _, pod := range m.pods {
		rows = append(rows, table.Row{
			pod.Namespace,
			pod.Name,
		})
	}
	m.table.SetRows(rows)
}
func (m PodModel) View() string {

	return m.table.View()

}
