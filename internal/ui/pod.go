package ui

import (
	"sort"
	"strconv"

	"github.com/aemery-cb/pequod/internal/common"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
)

type PodModel struct {
	table table.Model
	pods  []v1.Pod
}

func NewPodModel() PodModel {
	columns := []table.Column{
		{Title: "id", Width: 4},
		{Title: "Namespace", Width: 25},
		{Title: "Name", Width: 30},
	}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	return PodModel{table: t}
}

func (m PodModel) Init() tea.Cmd {
	return nil
}

func (m PodModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
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
			index, _ := strconv.Atoi(m.table.SelectedRow()[0])
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
	for index, pod := range m.pods {
		rows = append(rows, table.Row{
			strconv.Itoa(index),
			pod.Namespace,
			pod.Name,
		})
	}
	m.table.SetRows(rows)
}
func (m PodModel) View() string {

	return m.table.View()
	// s += "\nselect a pod using ↑ ↓ and press enter to stream logs"
	// return s
}
